package domain

import "errors"

// Errores genéricos del dominio. Otros errores los envuelven con %w
// para que se puedan detectar con errors.Is.
//
// Convención: si una función de dominio falla por datos inválidos,
// devuelve un error envuelto sobre ErrInvalidInput.
// Si falla porque algo no existe, sobre ErrNotFound. Etc.
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
)
