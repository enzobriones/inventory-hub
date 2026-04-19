// Package domain contiene las entidades, reglas de negocio e invariantes
// del hub de inventario. No depende de ningún detalle de infraestructura.
package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// TenantID identifica unívocamente a una organización cliente del hub.
// Todas las entidades del sistema llevan un TenantID para aislarlas entre clientes.
type TenantID string

// ProductID identifica unívocamente a un producto del catálogo.
type ProductID string

// VariantID identifica unívocamente a una variante concreta de un producto
// (ej: "polera M roja"). Es la unidad real de venta y de stock.
type VariantID string

// LocationID identifica unívocamente a una ubicación física (tienda, bodega, sucursal).
type LocationID string

// ChannelID identifica unívocamente a un canal de venta configurado por el tenant
// (ej: una cuenta específica de Mercado Libre).
type ChannelID string

// InventoryID identifica unívocamente a una fila de la proyección de inventario,
// que asocia (variante, ubicación) con cantidades actuales.
type InventoryID string

// MovementID identifica unívocamente a un movimiento de stock (entrada, salida,
// ajuste, transferencia). Los movimientos son la fuente de verdad del inventario.
type MovementID string

// ReservationID identifica unívocamente a una reserva temporal de stock.
type ReservationID string

// SaleID identifica unívocamente a una venta consolidada en el hub,
// sin importar el canal de origen.
type SaleID string

// SaleItemID identifica unívocamente a una línea de una venta.
type SaleItemID string

// NewTenantID genera un nuevo TenantID aleatorio (UUID v4).
func NewTenantID() TenantID { return TenantID(uuid.NewString()) }

// NewProductID genera un nuevo ProductID aleatorio (UUID v4).
func NewProductID() ProductID { return ProductID(uuid.NewString()) }

// NewVariantID genera un nuevo VariantID aleatorio (UUID v4).
func NewVariantID() VariantID { return VariantID(uuid.NewString()) }

// NewLocationID genera un nuevo LocationID aleatorio (UUID v4).
func NewLocationID() LocationID { return LocationID(uuid.NewString()) }

// NewChannelID genera un nuevo ChannelID aleatorio (UUID v4).
func NewChannelID() ChannelID { return ChannelID(uuid.NewString()) }

// NewInventoryID genera un nuevo InventoryID aleatorio (UUID v4).
func NewInventoryID() InventoryID { return InventoryID(uuid.NewString()) }

// NewMovementID genera un nuevo MovementID aleatorio (UUID v4).
func NewMovementID() MovementID { return MovementID(uuid.NewString()) }

// NewReservationID genera un nuevo ReservationID aleatorio (UUID v4).
func NewReservationID() ReservationID { return ReservationID(uuid.NewString()) }

// NewSaleID genera un nuevo SaleID aleatorio (UUID v4).
func NewSaleID() SaleID { return SaleID(uuid.NewString()) }

// NewSaleItemID genera un nuevo SaleItemID aleatorio (UUID v4).
func NewSaleItemID() SaleItemID { return SaleItemID(uuid.NewString()) }

// parseUUID verifica que un string sea un UUID válido y no vacío.
// Helper interno reutilizado por los parsers públicos.
func parseUUID(s string) error {
	if s == "" {
		return fmt.Errorf("%w: id is empty", ErrInvalidInput)
	}
	if _, err := uuid.Parse(s); err != nil {
		return fmt.Errorf("%w: id %q is not a valid uuid", ErrInvalidInput, s)
	}
	return nil
}

// ParseTenantID valida un string UUID y lo convierte a TenantID tipado.
func ParseTenantID(s string) (TenantID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return TenantID(s), nil
}

// ParseProductID valida un string UUID y lo convierte a ProductID tipado.
func ParseProductID(s string) (ProductID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return ProductID(s), nil
}

// ParseVariantID valida un string UUID y lo convierte a VariantID tipado.
func ParseVariantID(s string) (VariantID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return VariantID(s), nil
}

// ParseLocationID valida un string UUID y lo convierte a LocationID tipado.
func ParseLocationID(s string) (LocationID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return LocationID(s), nil
}

// ParseChannelID valida un string UUID y lo convierte a ChannelID tipado.
func ParseChannelID(s string) (ChannelID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return ChannelID(s), nil
}

// ParseSaleID valida un string UUID y lo convierte a SaleID tipado.
func ParseSaleID(s string) (SaleID, error) {
	if err := parseUUID(s); err != nil {
		return "", err
	}
	return SaleID(s), nil
}
