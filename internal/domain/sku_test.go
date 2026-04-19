package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/enzobriones/inventory-hub/internal/domain"
)

func TestNewSKU(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		// Válidos — forma básica
		{"simple alpha", "ABC", "ABC", nil},
		{"simple numeric", "12345", "12345", nil},
		{"typical retail", "POL-2026-ROJ-M", "POL-2026-ROJ-M", nil},
		{"with underscore", "ABC_123", "ABC_123", nil},
		{"with dot", "ABC.123", "ABC.123", nil},
		{"mixed separators", "A-B_C.D", "A-B_C.D", nil},
		{"ean-like numeric", "7801234567890", "7801234567890", nil},

		// Normalización
		{"lowercase normalized", "pol-2026", "POL-2026", nil},
		{"mixed case normalized", "Pol-2026-Roj", "POL-2026-ROJ", nil},
		{"trimmed spaces", "  ABC-123  ", "ABC-123", nil},
		{"trim then upper", "  pol-m  ", "POL-M", nil},

		// Inválidos — longitud
		{"empty", "", "", domain.ErrInvalidInput},
		{"only spaces", "   ", "", domain.ErrInvalidInput},
		{"too short", "AB", "", domain.ErrInvalidInput},
		{"too long", strings.Repeat("A", domain.MaxSKULength+1), "", domain.ErrInvalidInput},

		// Inválidos — formato
		{"internal space", "POL 2026", "", domain.ErrInvalidInput},
		{"leading dash", "-ABC", "", domain.ErrInvalidInput},
		{"trailing dash", "ABC-", "", domain.ErrInvalidInput},
		{"double dash", "ABC--123", "", domain.ErrInvalidInput},
		{"double dot", "ABC..123", "", domain.ErrInvalidInput},
		{"special char at", "ABC@123", "", domain.ErrInvalidInput},
		{"special char slash", "ABC/123", "", domain.ErrInvalidInput},
		{"only separators", "---", "", domain.ErrInvalidInput},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := domain.NewSKU(tc.input)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("esperaba sin error, obtuve %v", err)
				}
				if s.String() != tc.want {
					t.Errorf("esperaba %q, obtuve %q", tc.want, s.String())
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("esperaba error %v, obtuve %v", tc.wantErr, err)
			}
		})
	}
}

func TestSKUEqual(t *testing.T) {
	a, err := domain.NewSKU("POL-2026-M")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	b, err := domain.NewSKU("pol-2026-m") // igual tras normalizar
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	c, err := domain.NewSKU("POL-2026-L")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if !a.Equal(b) {
		t.Errorf("esperaba SKUs iguales tras normalizar, obtuve distintos")
	}
	if a.Equal(c) {
		t.Errorf("esperaba SKUs distintos, obtuve iguales")
	}
}

func TestSKUIsZero(t *testing.T) {
	var zero domain.SKU
	if !zero.IsZero() {
		t.Errorf("esperaba IsZero=true para valor cero")
	}
	s, err := domain.NewSKU("ABC-123")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if s.IsZero() {
		t.Errorf("esperaba IsZero=false para SKU inicializado")
	}
}
