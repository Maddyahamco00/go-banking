package http

import (
	"net/http"

	"gobanking-v2/internal/health"
)

func NewRouter(healthHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler)
	return mux
}

