package service

import (
	"bytes"
	"testing"
)

func TestValidateTargetURL(t *testing.T) {
	svc := NewQRService()

	validURLs := []string{
		"https://example.com",
		"http://localhost:8080/test",
		"https://sub.domain.com.br/path?query=123#hash",
	}

	for _, u := range validURLs {
		if err := svc.ValidateTargetURL(u); err != nil {
			t.Errorf("expected valid URL for '%s', got error: %v", u, err)
		}
	}

	invalidURLs := []string{
		"",
		"ftp://example.com",
		"javascript:alert(1)",
		"example.com",
		"just-text",
		"file:///etc/passwd",
	}

	for _, u := range invalidURLs {
		if err := svc.ValidateTargetURL(u); err == nil {
			t.Errorf("expected invalid URL for '%s', got nil error", u)
		}
	}
}

func TestValidateSlug(t *testing.T) {
	svc := NewQRService()

	validSlugs := []string{
		"abc",
		"cardapio",
		"mesa-12",
		"mesa_14",
		"promo2026",
	}

	for _, s := range validSlugs {
		if err := svc.ValidateSlug(s); err != nil {
			t.Errorf("expected valid slug for '%s', got: %v", s, err)
		}
	}

	invalidSlugs := []string{
		"ab",                            // less than 3 chars
		"slug with space",
		"slug@special",
		"slug!",
		"super_long_slug_that_exceeds_thirty_two_characters_limit_test",
	}

	for _, s := range invalidSlugs {
		if err := svc.ValidateSlug(s); err == nil {
			t.Errorf("expected invalid slug for '%s', got nil", s)
		}
	}
}

func TestGenerateRandomSlug(t *testing.T) {
	svc := NewQRService()

	slug1, err := svc.GenerateRandomSlug(6)
	if err != nil {
		t.Fatalf("unexpected error generating slug: %v", err)
	}

	if len(slug1) != 6 {
		t.Errorf("expected slug length 6, got %d (%s)", len(slug1), slug1)
	}

	slug2, err := svc.GenerateRandomSlug(6)
	if err != nil {
		t.Fatalf("unexpected error generating slug: %v", err)
	}

	if slug1 == slug2 {
		t.Errorf("expected distinct random slugs, but got identical: %s", slug1)
	}
}

func TestBase32EncodingRoundtrip(t *testing.T) {
	svc := NewQRService()

	originalData := []byte("fake-png-binary-stream-for-testing-purposes-1234567890")
	encoded := svc.EncodeBase32(originalData)

	if encoded == "" {
		t.Fatal("expected non-empty Base32 string")
	}

	decoded, err := svc.DecodeBase32(encoded)
	if err != nil {
		t.Fatalf("unexpected error decoding Base32: %v", err)
	}

	if !bytes.Equal(originalData, decoded) {
		t.Errorf("roundtrip data mismatch: expected %v, got %v", originalData, decoded)
	}
}

func TestGenerateBase32QRCode(t *testing.T) {
	svc := NewQRService()

	url := "http://localhost:8080/test-slug"
	base32Img, err := svc.GenerateBase32QRCode(url)
	if err != nil {
		t.Fatalf("failed to generate base32 qr code: %v", err)
	}

	if base32Img == "" {
		t.Fatal("generated base32 qr code is empty")
	}

	// Decodifica para verificar integridade do cabeçalho PNG (\x89PNG\r\n\x1a\n)
	pngBytes, err := svc.DecodeBase32(base32Img)
	if err != nil {
		t.Fatalf("failed to decode generated base32 qr image: %v", err)
	}

	pngHeader := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if !bytes.HasPrefix(pngBytes, pngHeader) {
		t.Errorf("decoded bytes do not start with PNG magic header")
	}
}
