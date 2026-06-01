package wallet

// WalletStatus represents lifecycle state for a wallet.
//
// FinTech decision: status gates money operations.
// - ACTIVE: credits/debits allowed
// - SUSPENDED: credits/debits rejected (future: partial actions)
// - CLOSED: terminal state, no further mutations
//
// We keep mapping/parsing logic close to the domain to ensure consistent validation.

type WalletStatus int

const (
	WalletStatusUnknown WalletStatus = iota
	WalletStatusActive
	WalletStatusSuspended
	WalletStatusClosed
)

func WalletStatusFromString(s string) WalletStatus {
	switch s {
	case "ACTIVE":
		return WalletStatusActive
	case "SUSPENDED":
		return WalletStatusSuspended
	case "CLOSED":
		return WalletStatusClosed
	default:
		return WalletStatusUnknown
	}
}

func (ws WalletStatus) String() string {
	switch ws {
	case WalletStatusActive:
		return "ACTIVE"
	case WalletStatusSuspended:
		return "SUSPENDED"
	case WalletStatusClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}
