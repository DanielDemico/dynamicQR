package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"dynamicqr/internal/config"
	"dynamicqr/internal/domain"
	"dynamicqr/internal/repository"
	"dynamicqr/internal/service"

	"github.com/go-chi/chi/v5"
)

type QRHandler struct {
	cfg       *config.Config
	repo      *repository.QRCodeRepository
	qrService *service.QRService
}

func NewQRHandler(cfg *config.Config, repo *repository.QRCodeRepository, qrService *service.QRService) *QRHandler {
	return &QRHandler{
		cfg:       cfg,
		repo:      repo,
		qrService: qrService,
	}
}

func (h *QRHandler) resolveBaseURL(r *http.Request) string {
	if h.cfg.BaseURL != "" && !strings.Contains(h.cfg.BaseURL, "localhost") {
		return strings.TrimRight(h.cfg.BaseURL, "/")
	}

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	if r.Host != "" {
		return fmt.Sprintf("%s://%s", scheme, r.Host)
	}

	return strings.TrimRight(h.cfg.BaseURL, "/")
}

// Create cria um novo QR Code codificado em Base32
func (h *QRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateQRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload JSON inválido"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "O título é obrigatório"})
		return
	}

	targetURL := strings.TrimSpace(req.TargetURL)
	if err := h.qrService.ValidateTargetURL(targetURL); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var slug string
	customSlug := strings.TrimSpace(req.CustomSlug)
	if customSlug != "" {
		if err := h.qrService.ValidateSlug(customSlug); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		slug = strings.ToLower(customSlug)
	} else {
		// Gera slug aleatório com tentativa de retry
		for i := 0; i < 3; i++ {
			s, err := h.qrService.GenerateRandomSlug(6)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao gerar identificador"})
				return
			}
			// Verifica se já existe
			if _, err := h.repo.FindBySlug(r.Context(), s); errors.Is(err, repository.ErrNotFound) {
				slug = s
				break
			}
		}
		if slug == "" {
			respondJSON(w, http.StatusConflict, map[string]string{"error": "Falha ao gerar slug único. Tente novamente."})
			return
		}
	}

	baseURL := h.resolveBaseURL(r)
	shortURL := baseURL + "/" + slug

	// Gera imagem PNG e serializa em Base32
	base32Img, err := h.qrService.GenerateBase32QRCode(shortURL)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao gerar imagem do QR code"})
		return
	}

	qr := &domain.QRCode{
		Title:         title,
		Slug:          slug,
		TargetURL:     targetURL,
		QRImageBase32: base32Img,
		QRMimeType:    "image/png",
	}

	if err := h.repo.Create(r.Context(), qr); err != nil {
		if errors.Is(err, repository.ErrSlugConflict) {
			respondJSON(w, http.StatusConflict, map[string]string{"error": "O slug informado já está em uso."})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao persistir no banco de dados"})
		return
	}

	respondJSON(w, http.StatusCreated, qr.ToResponse(baseURL))
}

// List lista QR Codes paginados
func (h *QRHandler) List(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	items, total, err := h.repo.List(r.Context(), search, page, limit)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao consultar registros"})
		return
	}

	baseURL := h.resolveBaseURL(r)
	respItems := make([]domain.QRCodeResponse, len(items))
	for i, item := range items {
		respItems[i] = item.ToResponse(baseURL)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	response := domain.ListQRCodeResponse{
		Data: respItems,
		Pagination: domain.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// GetByID busca um QR Code por ID
func (h *QRHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qr, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "QR Code não encontrado."})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao consultar registro"})
		return
	}

	baseURL := h.resolveBaseURL(r)
	respondJSON(w, http.StatusOK, qr.ToResponse(baseURL))
}

// Update edita dinamicamente o QR Code (preserva slug e imagem)
func (h *QRHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.UpdateQRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload JSON inválido"})
		return
	}

	if req.TargetURL != nil {
		trimmed := strings.TrimSpace(*req.TargetURL)
		if err := h.qrService.ValidateTargetURL(trimmed); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		req.TargetURL = &trimmed
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "O título não pode ser vazio"})
			return
		}
		req.Title = &trimmed
	}

	updated, err := h.repo.Update(r.Context(), id, req.Title, req.TargetURL, req.IsActive)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "QR Code não encontrado."})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar registro"})
		return
	}

	baseURL := h.resolveBaseURL(r)
	respondJSON(w, http.StatusOK, updated.ToResponse(baseURL))
}

// Delete remove um QR Code
func (h *QRHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "QR Code não encontrado."})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao remover registro"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetImage decodifica a imagem salva em Base32 no MongoDB e entrega como PNG
func (h *QRHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qr, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "QR Code não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao recuperar imagem", http.StatusInternalServerError)
		return
	}

	pngBytes, err := h.qrService.DecodeBase32(qr.QRImageBase32)
	if err != nil {
		http.Error(w, "Erro ao decodificar imagem Base32", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(pngBytes)))

	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="qrcode-%s.png"`, qr.Slug))
	}

	_, _ = w.Write(pngBytes)
}

// Healthz endpoint de verificação de integridade e conectividade com MongoDB
func (h *QRHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	// Faz uma consulta leve para garantir que o MongoDB está respondendo
	_, _, err := h.repo.List(r.Context(), "", 1, 1)
	if err != nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "error",
			"mongo":  "disconnected",
			"detail": err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"mongo":  "connected",
	})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
