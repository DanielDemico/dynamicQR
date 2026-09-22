package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dynamicqr/internal/repository"

	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	repo *repository.QRCodeRepository
}

func NewRedirectHandler(repo *repository.QRCodeRepository) *RedirectHandler {
	return &RedirectHandler{
		repo: repo,
	}
}

// Redirect executa o redirecionamento público 302 com cabeçalhos anti-cache e contagem atômica
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		h.renderNotFound(w, r, "Identificador não fornecido")
		return
	}

	// Ignora rotas reservadas que possam cair aqui
	if slug == "favicon.ico" || slug == "robots.txt" || strings.HasPrefix(slug, "api") || strings.HasPrefix(slug, "static") {
		http.NotFound(w, r)
		return
	}

	qr, err := h.repo.FindBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.renderNotFound(w, r, fmt.Sprintf("O link dinâmico '/%s' não existe ou foi desativado.", slug))
			return
		}
		http.Error(w, "Erro interno ao processar redirecionamento", http.StatusInternalServerError)
		return
	}

	// Incremento atômico assíncrono para garantir latência ultra-baixa (< 15ms)
	go func(targetSlug string) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = h.repo.IncrementClickCount(ctx, targetSlug)
	}(qr.Slug)

	// Cabeçalhos anti-cache estritos: garante que alterações imediatas de destino sejam refletidas
	w.Header().Set("Location", qr.TargetURL)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusFound) // 302 Found
}

func (h *RedirectHandler) renderNotFound(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusNotFound)

	// Se for requisição de navegador, exibe página amigável
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - Link Não Encontrado | DynamicQR</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
        body {
            background: #0f172a;
            color: #f8fafc;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            padding: 20px;
        }
        .card {
            background: #1e293b;
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 16px;
            padding: 40px;
            max-width: 480px;
            text-align: center;
            box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
        }
        .icon {
            font-size: 56px;
            margin-bottom: 16px;
        }
        h1 { font-size: 24px; font-weight: 700; margin-bottom: 12px; color: #f1f5f9; }
        p { font-size: 15px; color: #94a3b8; line-height: 1.6; margin-bottom: 24px; }
        a {
            display: inline-block;
            background: #6366f1;
            color: white;
            padding: 12px 24px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 600;
            transition: background 0.2s;
        }
        a:hover { background: #4f46e5; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">🔍</div>
        <h1>Link Não Encontrado</h1>
        <p>%s</p>
        <a href="/">Ir para o Painel DynamicQR</a>
    </div>
</body>
</html>`, message)
		_, _ = w.Write([]byte(html))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Link não encontrado","detail":"%s"}`, message)))
}
