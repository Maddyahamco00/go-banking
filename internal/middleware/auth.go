package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-banking/internal/config"
	"go-banking/internal/pkg/response"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Tier   string `json:"tier"`
	jwt.RegisteredClaims
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			response.Unauthorized(c, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			response.Unauthorized(c, "Invalid token claims")
			c.Abort()
			return
		}

		// Check expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			response.Unauthorized(c, "Token has expired")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("tier", claims.Tier)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetRole extracts role from context
func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetRole(c)
		if !exists {
			response.Unauthorized(c, "Role not found in token")
			c.Abort()
			return
		}

		if role != requiredRole && role != "admin" {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken generates a new JWT token (for testing purposes)
func GenerateToken(cfg *config.JWTConfig, userID, role, tier string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}