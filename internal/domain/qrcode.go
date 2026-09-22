package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// QRCode representa o documento no MongoDB
type QRCode struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string        `bson:"title" json:"title"`
	Slug          string        `bson:"slug" json:"slug"`
	TargetURL     string        `bson:"target_url" json:"target_url"`
	ClickCount    int64         `bson:"click_count" json:"click_count"`
	IsActive      bool          `bson:"is_active" json:"is_active"`
	QRImageBase32 string        `bson:"qr_image_base32" json:"qr_image_base32"`
	QRMimeType    string        `bson:"qr_mime_type" json:"qr_mime_type"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}

// CreateQRCodeRequest payload para criar um novo QR Code
type CreateQRCodeRequest struct {
	Title      string `json:"title"`
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug,omitempty"`
}

// UpdateQRCodeRequest payload para editar um QR Code
type UpdateQRCodeRequest struct {
	Title     *string `json:"title,omitempty"`
	TargetURL *string `json:"target_url,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// QRCodeResponse retorno formatado da API com links compostos
type QRCodeResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	ShortURL       string    `json:"short_url"`
	TargetURL      string    `json:"target_url"`
	ClickCount     int64     `json:"click_count"`
	IsActive       bool      `json:"is_active"`
	QRCodeImageURL string    `json:"qr_code_image_url"`
	QRImageBase32  string    `json:"qr_image_base32,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Pagination metadados de paginação
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}

// ListQRCodeResponse lista paginada de QR Codes
type ListQRCodeResponse struct {
	Data       []QRCodeResponse `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// ToResponse converte a entidade para DTO com a short URL e URL da imagem
func (q *QRCode) ToResponse(baseURL string) QRCodeResponse {
	return QRCodeResponse{
		ID:             q.ID.Hex(),
		Title:          q.Title,
		Slug:           q.Slug,
		ShortURL:       baseURL + "/" + q.Slug,
		TargetURL:      q.TargetURL,
		ClickCount:     q.ClickCount,
		IsActive:       q.IsActive,
		QRCodeImageURL: "/api/v1/qr-codes/" + q.ID.Hex() + "/image",
		QRImageBase32:  q.QRImageBase32,
		CreatedAt:      q.CreatedAt,
		UpdatedAt:      q.UpdatedAt,
	}
}
