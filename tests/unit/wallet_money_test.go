package unit

import (
	"testing"

	"gobanking-v2/domain/wallet"
)

func TestMoney_FromIntCents_NegativeRejected(t *testing.T) {
	_, err := wallet.NewMoneyFromCents(-1)
	if err == nil {
		t.Fatalf("expected error for negative cents")
	}
}

func TestMoney_FromIntCents_ZeroAllowed(t *testing.T) {
	m, err := wallet.NewMoneyFromCents(0)
	if err != nil {
		t.Fatalf("expected no error for zero")
	}
	if m.Cents() != 0 {
		t.Fatalf("expected 0 cents")
	}
}

func TestMoney_Add(t *testing.T) {
	a, _ := wallet.NewMoneyFromCents(100)
	b, _ := wallet.NewMoneyFromCents(250)
	c := a.Add(b)
	if c.Cents() != 350 {
		t.Fatalf("expected 350 cents, got %d", c.Cents())
	}
}

func TestMoney_Sub_DoesNotAllowNegative(t *testing.T) {
	a, _ := wallet.NewMoneyFromCents(50)
	b, _ := wallet.NewMoneyFromCents(60)
	_, err := a.Sub(b)
	if err == nil {
		t.Fatalf("expected error when subtracting to negative")
	}
}

