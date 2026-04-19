package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// Longitudes permitidas para un SKU tras normalizar.
const (
	MinSKULength = 3
	MaxSKULength = 40
)

// skuRegex valida la forma canónica del SKU: uno o más segmentos
// alfanuméricos separados opcionalmente por '.', '_' o '-'. Sin
// separadores consecutivos ni en los extremos.
var skuRegex = regexp.MustCompile(`^[A-Z0-9]+(?:[._-][A-Z0-9]+)*$`)

// SKU es el código único de una ProductVariant dentro de un tenant.
// Se normaliza a uppercase en el constructor; toda comparación o
// persistencia usa esa forma canónica. Ejemplo: "POL-2026-ROJ-M".
type SKU struct {
	value string
}

// NewSKU valida y construye un SKU. Aplica trim de espacios y uppercase
// antes de validar. Devuelve ErrInvalidInput envuelto si el string
// resultante no cumple el formato canónico o la longitud permitida.
func NewSKU(s string) (SKU, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if len(s) < MinSKULength {
		return SKU{}, fmt.Errorf("sku too short (min %d): %w", MinSKULength, ErrInvalidInput)
	}
	if len(s) > MaxSKULength {
		return SKU{}, fmt.Errorf("sku too long (max %d): %w", MaxSKULength, ErrInvalidInput)
	}
	if !skuRegex.MatchString(s) {
		return SKU{}, fmt.Errorf("invalid sku format %q: %w", s, ErrInvalidInput)
	}
	return SKU{value: s}, nil
}

// String devuelve la representación textual del SKU (siempre uppercase).
func (s SKU) String() string {
	return s.value
}

// IsZero indica si el SKU es el valor cero (no inicializado).
func (s SKU) IsZero() bool {
	return s.value == ""
}

// Equal compara dos SKUs por valor canónico.
func (s SKU) Equal(other SKU) bool {
	return s.value == other.value
}
