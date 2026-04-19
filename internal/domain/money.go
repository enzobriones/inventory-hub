package domain

import "fmt"

var validCurrencies = map[string]struct{}{
	"CLP": {},
	"USD": {},
	"EUR": {},
}

// Money representa un monto en una moneda específica.
// amount se guarda en unidades mínimas: para CLP son pesos (sin decimales),
// para USD son centavos (1.99 USD = 199).
//
// Money es un value object inmutable: no hay setters, todas las operaciones
// devuelven un Money nuevo
type Money struct {
	amount   int64
	currency string
}

// NewMoney construye un Money validando la moneda
// El amount puede ser negativo (para ajustes, reversas), cero, o positivo.
// Las reglas de "precio debe ser > 0" van en la entidad que lo usa, no acá.
func NewMoney(amount int64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("%w: currency is empty", ErrInvalidInput)
	}
	if _, ok := validCurrencies[currency]; !ok {
		return Money{}, fmt.Errorf("%w: unsupported currency %q", ErrInvalidInput, currency)
	}
	return Money{amount: amount, currency: currency}, nil
}

// MustMoney es útil para tests y constantes. Panics si la construcción falla.
// NO usar en código de producción - siempre preferir NewMoney
func MustMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// Accessors: son solo-lectura para preservar inmutabilidad.
func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) IsZero() bool     { return m.amount == 0 }
func (m Money) IsNegative() bool { return m.amount < 0 }
func (m Money) IsPositive() bool { return m.amount > 0 }

// Add suma dos montos. Falla si las monedas no coinciden.
func (m Money) Add(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: currency mismatch: %s vs %s",
			ErrInvalidInput, m.currency, o.currency)
	}
	return Money{amount: m.amount + o.amount, currency: m.currency}, nil
}

// Sub resta dos montos. Falla si las monedas no coinciden.
func (m Money) Sub(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: currency mismatch: %s vs %s",
			ErrInvalidInput, m.currency, o.currency)
	}
	return Money{amount: m.amount - o.amount, currency: m.currency}, nil
}

// Mul multiplica el monto por un entero (ej: precio unitario × cantidad).
// No falla porque no involucra conversión de moneda.
func (m Money) Mul(n int64) Money {
	return Money{amount: m.amount * n, currency: m.currency}
}

// Equal compara dos montos por valor.
func (m Money) Equal(o Money) bool {
	return m.amount == o.amount && m.currency == o.currency
}

// String formatea el monto para logs/debug. NO es el formato para mostrar
// al usuario — eso es responsabilidad de la capa de presentación.
func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.currency)
}
