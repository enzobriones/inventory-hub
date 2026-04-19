package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/enzobriones/inventory-hub/internal/domain"
)

func TestNewHandle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid simple", "polera", nil},
		{"valid with numbers", "polera-2026", nil},
		{"valid multi-segment", "polera-manga-corta-2026", nil},
		{"valid only digits", "12345", nil},
		{"empty", "", domain.ErrInvalidInput},
		{"uppercase", "Polera", domain.ErrInvalidInput},
		{"with space", "polera manga", domain.ErrInvalidInput},
		{"leading dash", "-polera", domain.ErrInvalidInput},
		{"trailing dash", "polera-", domain.ErrInvalidInput},
		{"double dash", "polera--manga", domain.ErrInvalidInput},
		{"special char", "polera!", domain.ErrInvalidInput},
		{"accent", "polerá", domain.ErrInvalidInput},
		{"too long", strings.Repeat("a", domain.MaxHandleLength+1), domain.ErrInvalidInput},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, err := domain.NewHandle(tc.input)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("esperaba sin error, obtuve %v", err)
				}
				if h.String() != tc.input {
					t.Errorf("esperaba %q, obtuve %q", tc.input, h.String())
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("esperaba error %v, obtuve %v", tc.wantErr, err)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Polera", "polera"},
		{"multi word", "Polera Manga Corta", "polera-manga-corta"},
		{"with numbers", "Polera 2026", "polera-2026"},
		{"accents", "Camiseta Básica", "camiseta-basica"},
		{"with n-tilde", "Camión Ñandú", "camion-nandu"},
		{"special chars", "¡Oferta! 50%", "oferta-50"},
		{"multiple spaces", "polera   basica", "polera-basica"},
		{"leading trailing spaces", "  polera  ", "polera"},
		{"mixed symbols", "Ñandú Azul & Rojo", "nandu-azul-rojo"},
		{"all invalid", "¡!¡!", ""},
		{"empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := domain.Slugify(tc.in); got != tc.want {
				t.Errorf("esperaba %q, obtuve %q", tc.want, got)
			}
		})
	}
}

// TestSlugifyProducesValidHandle verifica la propiedad de composición:
// todo slug no vacío debe ser aceptado por NewHandle.
func TestSlugifyProducesValidHandle(t *testing.T) {
	inputs := []string{
		"Polera Manga Corta",
		"Cámara Réflex Nikon D850",
		"Teclado Mecánico RGB 60%",
		"Año Nuevo 2026",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			s := domain.Slugify(in)
			if s == "" {
				t.Fatalf("slug vacío para %q", in)
			}
			if _, err := domain.NewHandle(s); err != nil {
				t.Errorf("slug %q rechazado por NewHandle: %v", s, err)
			}
		})
	}
}
