package unit

import (
	"testing"

	"gobanking-v2/internal/health"
)

func TestHealthUsecase_Check_ReturnsOK(t *testing.T) {
	uc := health.New("gobanking")

	res, err := uc.Check()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if res.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", res.Status)
	}
	if res.Service != "gobanking" {
		t.Fatalf("expected service=gobanking, got %q", res.Service)
	}
}

