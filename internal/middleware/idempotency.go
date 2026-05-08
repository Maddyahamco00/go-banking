package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-redis/redis/v8"
	"go-banking/internal/repository"
)

// IdempotencyMiddleware handles duplicate request prevention
type IdempotencyMiddleware struct {
	redisClient *redis.Client
	repo        *repository.IdempotencyRepository
	ttlHours    int
}

// NewIdempotencyMiddleware creates a new IdempotencyMiddleware
func NewIdempotencyMiddleware(redisClient *redis.Client, repo *repository.IdempotencyRepository, ttlHours int) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		redisClient: redisClient,
		repo:        repo,
		ttlHours:    ttlHours,
	}
}

// responseCapture wraps gin.ResponseWriter to capture the response body
type responseCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Handle checks for duplicate requests based on X-Idempotency-Key header
func (m *IdempotencyMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for GET requests
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader("X-Idempotency-Key")
		if idempotencyKey == "" {
			// No idempotency key provided - continue without idempotency protection
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// Check Redis first (fast path)
		cacheKey := "idempotency:" + idempotencyKey
		cachedResponse, err := m.redisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedResponse != "" {
			// Found in cache - return cached response
			var resp map[string]interface{}
			if json.Unmarshal([]byte(cachedResponse), &resp) == nil {
				c.AbortWithStatusJSON(http.StatusOK, resp)
				return
			}
		}

		// Check PostgreSQL (persistent storage)
		record, err := m.repo.Get(ctx, idempotencyKey)
		if err == nil && record != nil {
			// Found in DB - return cached response
			var resp map[string]interface{}
			if json.Unmarshal(record.Response, &resp) == nil {
				// Cache in Redis for faster subsequent lookups
				m.redisClient.Set(ctx, cacheKey, record.Response, time.Duration(m.ttlHours)*time.Hour)

				c.AbortWithStatusJSON(http.StatusOK, resp)
				return
			}
		}

		// Generate a unique request ID for tracking
		requestID := uuid.New().String()
		c.Set("idempotency_key", idempotencyKey)
		c.Set("request_id", requestID)

		// Capture response
		capture := &responseCapture{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = capture

		c.Next()

		// After handler execution, store the response if successful
		if c.IsAborted() {
			return
		}

		// Only cache successful responses
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			responseBody := capture.body.String()
			if responseBody != "" {
				var respData map[string]interface{}
				if json.Unmarshal([]byte(responseBody), &respData) == nil {
					respJSON, _ := json.Marshal(respData)

					// Store in PostgreSQL
					m.repo.Set(ctx, idempotencyKey, respJSON, m.ttlHours)

					// Cache in Redis
					m.redisClient.Set(ctx, cacheKey, respJSON, time.Duration(m.ttlHours)*time.Hour)
				}
			}
		}
	}
}

// SkipIdempotency returns true if the request should skip idempotency check
func SkipIdempotency(c *gin.Context) bool {
	method := c.Request.Method
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}