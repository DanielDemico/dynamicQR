package service

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strings"

	"github.com/skip2/go-qrcode"
)

var (
	ErrInvalidURL      = errors.New("URL de destino inválida. Apenas links HTTP/HTTPS são permitidos")
	ErrInvalidSlug     = errors.New("Slug inválido. Use entre 3 e 32 caracteres alfanuméricos, hífen ou sublinhado")
	slugRegex          = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)
	charsetSlug        = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type QRService struct{}

func NewQRService() *QRService {
	return &QRService{}
}

// ValidateTargetURL verifica se a URL possui protocolo HTTP ou HTTPS e host válido
func (s *QRService) ValidateTargetURL(rawURL string) error {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ErrInvalidURL
	}

	u, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return ErrInvalidURL
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return ErrInvalidURL
	}

	if u.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// ValidateSlug valida o formato de um slug customizado
func (s *QRService) ValidateSlug(slug string) error {
	if !slugRegex.MatchString(slug) {
		return ErrInvalidSlug
	}
	return nil
}

// GenerateRandomSlug gera um slug aleatório criptograficamente seguro
func (s *QRService) GenerateRandomSlug(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	b := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charsetSlug)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("falha ao gerar caractere aleatório: %w", err)
		}
		b[i] = charsetSlug[num.Int64()]
	}

	return string(b), nil
}

// GenerateQRCodePNG gera a imagem PNG do QR code em memória
func (s *QRService) GenerateQRCodePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	pngBytes, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("falha ao codificar QR code: %w", err)
	}
	return pngBytes, nil
}

// EncodeBase32 codifica bytes em string Base32 (RFC 4648)
func (s *QRService) EncodeBase32(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

// DecodeBase32 decodifica string Base32 para os bytes originais
func (s *QRService) DecodeBase32(encoded string) ([]byte, error) {
	data, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar Base32: %w", err)
	}
	return data, nil
}

// GenerateBase32QRCode gera o QR code para a URL fornecida e já retorna a imagem em Base32
func (s *QRService) GenerateBase32QRCode(fullURL string) (string, error) {
	pngBytes, err := s.GenerateQRCodePNG(fullURL, 256)
	if err != nil {
		return "", err
	}
	return s.EncodeBase32(pngBytes), nil
}
