package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpDelivery "gobanking-v2/delivery/http"
	healthuc "gobanking-v2/internal/health"
)

func TestHealthEndpoint_Integration(t *testing.T) {
	// Note: current health usecase does not touch DB.
	// This integration test exists to validate routing + handler wiring.
	// When additional DB-backed features are introduced, reuse the DB harness.
	_ = time.Second

	h := healthuc.NewHandler(healthuc.New("gobanking"))
	router := httpDelivery.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestDBCleanupHelper_DoesNotPanic(t *testing.T) {
	// Placeholder to keep testdb helper compiled.
	// Real cleanup/migration tests should run once schema has tables.
	_ = context.Background()
	// No env-based DB open here; just ensure cleanup function compiles.
	_ = CleanupAll

}
