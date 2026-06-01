package wallet

import "fmt"

// Money represents an amount of currency in the smallest unit (cents/kobo/etc.)
// using an integer.
//
// FinTech decision: DO NOT use floats.
// - floats cannot exactly represent most decimal fractions
// - rounding drift accumulates across credits/debits
// - ledger reconciliation becomes non-deterministic
//
// This project uses integer cents/kobo to guarantee deterministic behavior.

type Money struct {
	cents int64
}

func NewMoneyFromCents(cents int64) (Money, error) {
	// Negative amounts are allowed only at the caller layer (e.g., validation).
	// For construction we reject negative values to enforce the invariant.
	if cents < 0 {
		return Money{}, fmt.Errorf("money_cents_must_be_non_negative")
	}
	return Money{cents: cents}, nil
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Add(other Money) Money {
	return Money{cents: m.cents + other.cents}
}

func (m Money) Sub(other Money) (Money, error) {
	if m.cents < other.cents {
		return Money{}, fmt.Errorf("money_subtract_would_go_negative")
	}
	return Money{cents: m.cents - other.cents}, nil
}
