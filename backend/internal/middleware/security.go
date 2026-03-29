package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	// AllowedOrigins specifies CORS allowed origins
	AllowedOrigins []string
	// AllowedMethods specifies CORS allowed methods
	AllowedMethods []string
	// AllowedHeaders specifies CORS allowed headers
	AllowedHeaders []string
	// ExposedHeaders specifies CORS exposed headers
	ExposedHeaders []string
	// AllowCredentials indicates if credentials are allowed
	AllowCredentials bool
	// MaxAge specifies CORS max age
	MaxAge int
	// FrameOptions specifies X-Frame-Options header value
	FrameOptions string
	// ContentSecurityPolicy specifies CSP header value
	ContentSecurityPolicy string
	// HSTSEnabled enables HSTS header
	HSTSEnabled bool
	// HSTSMaxAge specifies HSTS max age
	HSTSMaxAge int
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig(frontendURL string) SecurityConfig {
	return SecurityConfig{
		AllowedOrigins: []string{frontendURL, "http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposedHeaders: []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:         86400, // 24 hours
		FrameOptions:    "DENY",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;",
		HSTSEnabled:    false, // Enable only in production with HTTPS
		HSTSMaxAge:     31536000, // 1 year
	}
}

// SecurityHeaders adds security headers to all responses
func SecurityHeaders(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", config.FrameOptions)

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS filter
		c.Header("X-XSS-Protection", "1; mode=block")

		// Force HTTPS (only in production)
		if config.HSTSEnabled {
			c.Header("Strict-Transport-Security",
				strings.Join([]string{"max-age=", string(rune(config.HSTSMaxAge))}, "; includeSubDomains"))
		}

		// Content Security Policy
		if config.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		// Remove server information
		c.Header("Server", "")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy",
			"geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// EnhancedCORS provides enhanced CORS configuration
func EnhancedCORS(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowedOrigin := ""
		for _, allowed := range config.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = allowed
				break
			}
		}

		// Set CORS headers
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)

			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			c.Header("Access-Control-Allow-Methods",
				strings.Join(config.AllowedMethods, ", "))

			c.Header("Access-Control-Allow-Headers",
				strings.Join(config.AllowedHeaders, ", "))

			if len(config.ExposedHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers",
					strings.Join(config.ExposedHeaders, ", "))
			}

			c.Header("Access-Control-Max-Age",
				string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestSizeLimit limits the size of requests
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{
				"error": "Request entity too large",
				"code":   413,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// StripSlashes removes trailing slashes from URLs
func StripSlashes() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimSuffix(c.Request.URL.Path, "/")
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		c.Next()
	}
}
