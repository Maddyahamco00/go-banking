package http

import (
	"net/http"
)


func NewRouter(healthHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler)
	return mux
}

