package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Tipos de ID del dominio. Son strings internamente, pero el compilador
// los trata como tipos distintos. Intentar pasar un ProductID donde se
// espera un TenantID es un error de compilación - no un bug en runtime.

type (
	TenantID      string
	ProductID     string
	VariantID     string
	LocationID    string
	ChannelID     string
	InventoryID   string
	MovementID    string
	ReservationID string
	SaleID        string
	SaleItemID    string
)

// Generadores: crean un ID nuevo usando UUID v4
// Los usamos cuando el sistema es la fuente del ID (al crear entidades)
func NewTenantID() TenantID           { return TenantID(uuid.NewString()) }
func NewProductID() ProductID         { return ProductID(uuid.NewString()) }
func NewVariantID() VariantID         { return VariantID(uuid.NewString()) }
func NewLocationID() LocationID       { return LocationID(uuid.NewString()) }
func NewChannelID() ChannelID         { return ChannelID(uuid.NewString()) }
func NewInventoryID() InventoryID     { return InventoryID(uuid.NewString()) }
func NewMovementID() MovementID       { return MovementID(uuid.NewString()) }
func NewReservationID() ReservationID { return ReservationID(uuid.NewString()) }
func NewSaleID() SaleID               { return SaleID(uuid.NewString()) }
func NewSaleItemID() SaleItemID       { return SaleItemID(uuid.NewString()) }

// parseUUID es helper interno para validar que un string es un UUID válido.
func parseUUID(s string) error {
	if s == "" {
		return fmt.Errorf("%w: id is empty", ErrInvalidInput)
	}
	if _, err := uuid.Parse(s); err != nil {
		return fmt.Errorf("%w: id %q is not a valid uuid", ErrInvalidInput, s)
	}
	return nil
}

// Parsers: convierten strings externos (de HTTP, BD, etc.) a IDs tipados,
// validando el formato. Los usamos cuando el ID viene de afuera.
func ParseTenantID(s string) (TenantID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return TenantID(s), nil
}

func ParseProductID(s string) (ProductID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return ProductID(s), nil
}

func ParseVariantID(s string) (VariantID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return VariantID(s), nil
}

func ParseLocationID(s string) (LocationID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return LocationID(s), nil
}

func ParseChannelID(s string) (ChannelID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return ChannelID(s), nil
}

func ParseSaleID(s string) (SaleID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return SaleID(s), nil
}
