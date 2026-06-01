package wallet

import (
	"fmt"
	"time"
)

type Wallet struct {
	id        int64
	uuid      string
	userID    int64
	currency  string
	available Money
	held      Money
	status    WalletStatus
	createdAt time.Time
	updatedAt time.Time
}

// NewWallet constructs a domain wallet and enforces invariants.
//
// Invariants:
// - status must be a known business state
// - balances must be non-negative
//
// IDs are intentionally internal to the domain.
func NewWallet(status WalletStatus, userID int64, currency string, availableCents int64, heldCents int64, uuid string) (Wallet, error) {
	if status != WalletStatusActive && status != WalletStatusSuspended && status != WalletStatusClosed {
		return Wallet{}, fmt.Errorf("wallet_status_invalid")
	}
	if currency == "" {
		return Wallet{}, fmt.Errorf("wallet_currency_required")
	}

	available, err := NewMoneyFromCents(availableCents)
	if err != nil {
		return Wallet{}, fmt.Errorf("wallet_available_balance_invalid")
	}
	held, err := NewMoneyFromCents(heldCents)
	if err != nil {
		return Wallet{}, fmt.Errorf("wallet_held_balance_invalid")
	}

	now := time.Now().UTC()
	return Wallet{
		id:        0,
		uuid:      uuid,
		userID:    userID,
		currency:  currency,
		available: available,
		held:      held,
		status:    status,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (w *Wallet) AvailableBalanceCents() int64 {
	return w.available.Cents()
}

// Credit increases available balance.
//
// Validation rules for fintech safety:
// - amount must be strictly positive
// - wallet must be ACTIVE
// - amount is applied deterministically to integer cents
func (w *Wallet) Credit(amountCents int64) (Wallet, error) {
	if w.status != WalletStatusActive {
		return Wallet{}, fmt.Errorf("wallet_status_not_active")
	}
	if amountCents <= 0 {
		return Wallet{}, fmt.Errorf("credit_amount_must_be_positive")
	}
	m, err := NewMoneyFromCents(amountCents)
	if err != nil {
		return Wallet{}, err
	}
	w.available = w.available.Add(m)
	w.updatedAt = time.Now().UTC()
	return *w, nil
}

// Debit decreases available balance.
// Validation rules:
// - amount must be strictly positive
// - wallet must be ACTIVE
// - sufficient available balance required
func (w *Wallet) Debit(amountCents int64) (Wallet, error) {
	if w.status != WalletStatusActive {
		return Wallet{}, fmt.Errorf("wallet_status_not_active")
	}
	if amountCents <= 0 {
		return Wallet{}, fmt.Errorf("debit_amount_must_be_positive")
	}
	m, err := NewMoneyFromCents(amountCents)
	if err != nil {
		return Wallet{}, err
	}
	newAvail, err := w.available.Sub(m)
	if err != nil {
		return Wallet{}, fmt.Errorf("insufficient_available_balance")
	}
	w.available = newAvail
	w.updatedAt = time.Now().UTC()
	return *w, nil
}
