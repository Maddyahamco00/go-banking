package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go-banking/internal/config"
	"go-banking/internal/pkg/response"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requestsPerMin int
	users          map[string]*userBucket
	mu             sync.RWMutex
	cleanupInterval time.Duration
}

type userBucket struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg *config.RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		requestsPerMin: cfg.RequestsPerMin,
		users:          make(map[string]*userBucket),
		cleanupInterval: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes stale entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.users {
			if now.Sub(bucket.lastReset) > 10*time.Minute {
				delete(rl.users, key)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.users[key]

	if !exists {
		rl.users[key] = &userBucket{
			tokens:    rl.requestsPerMin - 1,
			lastReset:  now,
		}
		return true
	}

	// Reset bucket if minute has passed
	if now.Sub(bucket.lastReset) >= time.Minute {
		bucket.tokens = rl.requestsPerMin - 1
		bucket.lastReset = now
		return true
	}

	// Check if tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// Handle returns a Gin middleware for rate limiting
func (rl *RateLimiter) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use user ID if authenticated, otherwise use IP
		var key string
		if userID, exists := c.Get("user_id"); exists {
			key = "user:" + userID.(string)
		} else {
			key = "ip:" + c.ClientIP()
		}

		if !rl.Allow(key) {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}