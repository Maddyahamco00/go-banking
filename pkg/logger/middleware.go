package logger

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogging(lg *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewString()
			start := time.Now()
			ww := &ResponseWriter{ResponseWriter: w, status: http.StatusOK}

			lg2 := lg.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			)

			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(ww, r)

			lg2.Info("http_request",
				zap.Int("status", ww.status),
				zap.Duration("latency_ms", time.Since(start)),
			)
		})
	}
}

