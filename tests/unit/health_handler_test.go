package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gobanking-v2/internal/health"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func TestHealthHandler_ServeHTTP_GET_ReturnsOKJSON(t *testing.T) {
	uc := health.New("gobanking")
	h := health.NewHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct == "" {
		t.Fatalf("expected Content-Type header to be set")
	}

	var body healthResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid JSON body, got error: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", body.Status)
	}
	if body.Service != "gobanking" {
		t.Fatalf("expected service=gobanking, got %q", body.Service)
	}
}

func TestHealthHandler_ServeHTTP_NonGET_Returns405(t *testing.T) {
	uc := health.New("gobanking")
	h := health.NewHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", res.StatusCode)
	}
}

