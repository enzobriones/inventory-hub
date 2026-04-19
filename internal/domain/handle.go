package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// MaxHandleLength es la longitud máxima permitida en un handle
const MaxHandleLength = 100

// handleRegex valida la forma canónica de un handle: segmentos de
// [a-z0-9] separados por un único guion, sin guiones en los extremos
var handleRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Handle es un identificador legible y URL-safe de un Product.
// Ejemplo: "polera-manga-corta-2026"
type Handle struct {
	value string
}

// NewHandle valida y construye un Handle. Devuelve ErrInvalidInput
// envuelto si el string no cumple con el formato canónico
func NewHandle(s string) (Handle, error) {
	if s == "" {
		return Handle{}, fmt.Errorf("handle cannot be empty: %w", ErrInvalidInput)
	}
	if len(s) > MaxHandleLength {
		return Handle{}, fmt.Errorf("handle exceeds max length %d: %w", MaxHandleLength, ErrInvalidInput)
	}
	if !handleRegex.MatchString(s) {
		return Handle{}, fmt.Errorf("invalid handle format %q: %w", s, ErrInvalidInput)
	}
	return Handle{value: s}, nil
}

// String devuelve la representación textual del handle.
func (h Handle) String() string {
	return h.value
}

// IsZero indica si el hadnle es el valor cero (no inicializado).
func (h Handle) IsZero() bool {
	return h.value == ""
}

// handleAccentReplacer normaliza tildes y ñ del español a ASCII.
var handleAccentReplacer = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
	"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u", "Ü", "u", "Ñ", "n",
)

// handleInvalidChars matchea todo carácter que no sea [a-z0-9-].
var handleInvalidChars = regexp.MustCompile(`[^a-z0-9-]+`)

// handleDashRun colapsa secuencias de guiones en uno solo.
var handleDashRun = regexp.MustCompile(`-+`)

// Slugify convierte un nombre arbitrario en un candidato a handle.
// Pasos: lowercase, trim, normalizar tildes, espacios a guiones,
// eliminar caracteres inválidos, colapsar guiones, trim de guiones,
// truncar a MaxHandleLength. Puede devolver "" si el input no tiene
// caracteres válidos; el caller debe manejarlo.
func Slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = handleAccentReplacer.Replace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = handleInvalidChars.ReplaceAllString(s, "")
	s = handleDashRun.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > MaxHandleLength {
		s = strings.TrimRight(s[:MaxHandleLength], "-")
	}
	return s
}
