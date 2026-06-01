package unit

import (
	"testing"

	"gobanking-v2/domain/wallet"
)

func TestWallet_New_InvalidStatus(t *testing.T) {
	_, err := wallet.NewWallet(wallet.WalletStatusUnknown, 1, "USD", 0, 0, "")
	if err == nil {
		t.Fatalf("expected error for unknown status")
	}
}

func TestWallet_New_NonNegativeBalancesInvariant(t *testing.T) {
	w, err := wallet.NewWallet(wallet.WalletStatusActive, 1, "USD", -1, 0, "")
	if err == nil {
		t.Fatalf("expected error for negative available balance")
	}
	_ = w
}

func TestWallet_CanCreditOnlyWhenActive(t *testing.T) {
	w, err := wallet.NewWallet(wallet.WalletStatusSuspended, 1, "USD", 100, 0, "")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	_, err = w.Credit(10)
	if err == nil {
		t.Fatalf("expected error when crediting a suspended wallet")
	}
}

func TestWallet_CanDebitOnlyWhenActiveAndSufficientFunds(t *testing.T) {
	w, err := wallet.NewWallet(wallet.WalletStatusActive, 1, "USD", 50, 0, "")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	_, err = w.Debit(60)
	if err == nil {
		t.Fatalf("expected error when debiting more than available")
	}

	_, err = w.Debit(40)
	if err != nil {
		t.Fatalf("expected debit to succeed: %v", err)
	}
	if w.AvailableBalanceCents() != 10 {
		t.Fatalf("expected available balance to be 10, got %d", w.AvailableBalanceCents())
	}
}
