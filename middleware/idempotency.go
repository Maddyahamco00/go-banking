package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type requestState struct {
	done chan struct{}
}

var store = sync.Map{}

// IdempotencyMiddleware is a middleware that prevents duplicate requests
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.Next()
			return
		}

		if existing, exists := store.Load(key); exists {
			state := existing.(*requestState)
			select {
			case <-state.done:
				c.JSON(http.StatusConflict, gin.H{"error": "duplicate request"})
			default:
				c.JSON(http.StatusConflict, gin.H{"error": "request already in progress"})
			}
			c.Abort()
			return
		}

		state := &requestState{done: make(chan struct{})}
		store.Store(key, state)
		defer close(state.done)

		c.Next()

		if c.Writer.Status() >= http.StatusInternalServerError {
			store.Delete(key)
		}
	}
}
