package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-banking/internal/config"
	"go-banking/internal/pkg/response"
)

// HMACMiddleware handles HMAC request signing verification
type HMACMiddleware struct {
	cfg *config.HMACConfig
}

// NewHMACMiddleware creates a new HMACMiddleware
func NewHMACMiddleware(cfg *config.HMACConfig) *HMACMiddleware {
	return &HMACMiddleware{cfg: cfg}
}

// Handle verifies HMAC signatures on sensitive endpoints
func (m *HMACMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only verify on sensitive mutating endpoints
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		signature := c.GetHeader("X-Signature")
		timestampStr := c.GetHeader("X-Timestamp")

		if signature == "" || timestampStr == "" {
			response.BadRequest(c, "Missing X-Signature or X-Timestamp header")
			c.Abort()
			return
		}

		// Parse and validate timestamp
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid X-Timestamp format")
			c.Abort()
			return
		}

		requestTime := time.Unix(timestamp, 0)
		now := time.Now()

		// Reject requests older than 5 minutes
		if math.Abs(float64(now.Sub(requestTime).Minutes())) > 5 {
			response.BadRequest(c, "Request timestamp too old or too far in future")
			c.Abort()
			return
		}

		// Compute expected signature
		// Signature = HMAC-SHA256(timestamp + body, secret)
		body, _ := c.GetRawData()
		message := timestampStr + string(body)

		expectedSig := m.computeHMAC(message)
		if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
			response.Unauthorized(c, "Invalid HMAC signature")
			c.Abort()
			return
		}

		// Restore body for downstream handlers
		c.Request.Body = newBody(string(body))

		c.Next()
	}
}

// computeHMAC computes HMAC-SHA256
func (m *HMACMiddleware) computeHMAC(message string) string {
	h := hmac.New(sha256.New, []byte(m.cfg.Secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// bodyReader wraps a string as a reader
type bodyReader struct {
	content []byte
	pos     int
}

func newBody(content string) *bodyReader {
	return &bodyReader{content: []byte(content), pos: 0}
}

func (b *bodyReader) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.content) {
		return 0, nil
	}
	n = copy(p, b.content[b.pos:])
	b.pos += n
	return n, nil
}

func (b *bodyReader) Close() error {
	return nil
}