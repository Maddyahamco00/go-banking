package unit

import (
	"testing"

	"gobanking-v2/domain/wallet"
)

func TestWalletStatus_ValidValues(t *testing.T) {
	// ACTIVE
	if got := wallet.WalletStatusFromString("ACTIVE"); got == wallet.WalletStatusUnknown {
		t.Fatalf("expected ACTIVE to map to a known status")
	}

	// SUSPENDED
	if got := wallet.WalletStatusFromString("SUSPENDED"); got == wallet.WalletStatusUnknown {
		t.Fatalf("expected SUSPENDED to map to a known status")
	}

	// CLOSED
	if got := wallet.WalletStatusFromString("CLOSED"); got == wallet.WalletStatusUnknown {
		t.Fatalf("expected CLOSED to map to a known status")
	}
}

func TestWalletStatus_InvalidValue(t *testing.T) {
	if got := wallet.WalletStatusFromString("NOT_A_STATUS"); got != wallet.WalletStatusUnknown {
		t.Fatalf("expected invalid status to map to WalletStatusUnknown")
	}
}

