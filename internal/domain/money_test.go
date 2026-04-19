package domain_test

import (
	"errors"
	"testing"

	"github.com/enzobriones/inventory-hub/internal/domain"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		wantErr  error
	}{
		{"CLP positivo", 1990, "CLP", nil},
		{"USD positivo", 199, "USD", nil},
		{"EUR cero", 0, "EUR", nil},
		{"CLP negativo permitido", -500, "CLP", nil},
		{"moneda vacía falla", 100, "", domain.ErrInvalidInput},
		{"moneda no soportada falla", 100, "ARS", domain.ErrInvalidInput},
		{"moneda en minúscula falla", 100, "clp", domain.ErrInvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewMoney(tt.amount, tt.currency)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("esperaba éxito, obtuve error: %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("esperaba error %v, obtuve %v", tt.wantErr, err)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	t.Run("suma válida", func(t *testing.T) {
		a := domain.MustMoney(1000, "CLP")
		b := domain.MustMoney(500, "CLP")

		got, err := a.Add(b)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if got.Amount() != 1500 {
			t.Errorf("amount = %d, quería 1500", got.Amount())
		}
		if got.Currency() != "CLP" {
			t.Errorf("currency = %s, quería CLP", got.Currency())
		}
	})

	t.Run("monedas distintas falla", func(t *testing.T) {
		a := domain.MustMoney(1000, "CLP")
		b := domain.MustMoney(10, "USD")

		_, err := a.Add(b)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("esperaba ErrInvalidInput, obtuve %v", err)
		}
	})
}

func TestMoney_Sub(t *testing.T) {
	a := domain.MustMoney(1000, "CLP")
	b := domain.MustMoney(300, "CLP")

	got, err := a.Sub(b)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.Amount() != 700 {
		t.Errorf("amount = %d, quería 700", got.Amount())
	}
}

func TestMoney_Mul(t *testing.T) {
	price := domain.MustMoney(1990, "CLP")

	total := price.Mul(3)
	if total.Amount() != 5970 {
		t.Errorf("amount = %d, quería 5970", total.Amount())
	}
}

func TestMoney_Predicates(t *testing.T) {
	zero := domain.MustMoney(0, "CLP")
	positive := domain.MustMoney(100, "CLP")
	negative := domain.MustMoney(-100, "CLP")

	if !zero.IsZero() {
		t.Error("zero.IsZero() debería ser true")
	}
	if !positive.IsPositive() {
		t.Error("positive.IsPositive() debería ser true")
	}
	if !negative.IsNegative() {
		t.Error("negative.IsNegative() debería ser true")
	}
}

func TestMoney_Equal(t *testing.T) {
	a := domain.MustMoney(1000, "CLP")
	b := domain.MustMoney(1000, "CLP")
	c := domain.MustMoney(1000, "USD")
	d := domain.MustMoney(500, "CLP")

	if !a.Equal(b) {
		t.Error("a y b deberían ser iguales")
	}
	if a.Equal(c) {
		t.Error("a y c tienen moneda distinta, no deberían ser iguales")
	}
	if a.Equal(d) {
		t.Error("a y d tienen monto distinto, no deberían ser iguales")
	}
}
