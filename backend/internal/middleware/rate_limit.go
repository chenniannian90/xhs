package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	clients map[string]*ClientInfo
	mu      sync.RWMutex
	rate    int           // requests per minute
	burst   int           // maximum burst size
}

// ClientInfo tracks request information for each client
type ClientInfo struct {
	Tokens     int
	LastUpdate time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientInfo),
		rate:    rate,
		burst:   burst,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[clientID]
	if !exists {
		// New client, initialize with burst tokens
		rl.clients[clientID] = &ClientInfo{
			Tokens:     rl.burst - 1,
			LastUpdate: time.Now(),
		}
		return true
	}

	// Calculate time passed since last update
	elapsed := time.Since(client.LastUpdate)
	tokensToAdd := int(elapsed.Minutes()) * rl.rate

	// Add tokens based on rate, up to burst capacity
	client.Tokens += tokensToAdd
	if client.Tokens > rl.burst {
		client.Tokens = rl.burst
	}

	// Check if we have tokens available
	if client.Tokens > 0 {
		client.Tokens--
		client.LastUpdate = time.Now()
		return true
	}

	// No tokens available
	return false
}

// cleanup removes stale client entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for id, client := range rl.clients {
			// Remove clients that haven't been seen in 10 minutes
			if now.Sub(client.LastUpdate) > 10*time.Minute {
				delete(rl.clients, id)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
	// SkipRateLimit can be used to bypass rate limiting for certain IPs
	SkipRateLimit func(string) bool
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60, // 60 requests per minute
		Burst:             10,  // Allow burst of 10 requests
		SkipRateLimit: func(ip string) bool {
			// Skip rate limiting for localhost
			return ip == "127.0.0.1" || ip == "::1"
		},
	}
}

// RateLimit creates a rate limiting middleware
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config.RequestsPerMinute, config.Burst)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check if we should skip rate limiting
		if config.SkipRateLimit != nil && config.SkipRateLimit(clientIP) {
			c.Next()
			return
		}

		// Check rate limit
		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"code":   429,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByEndpoint creates rate limiting middleware with different rates per endpoint
func RateLimitByEndpoint(limits map[string]RateLimitConfig) gin.HandlerFunc {
	limiters := make(map[string]*RateLimiter)
	for path, config := range limits {
		limiters[path] = NewRateLimiter(config.RequestsPerMinute, config.Burst)
	}

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Find matching limit config
		var config RateLimitConfig
		var limiter *RateLimiter
		found := false

		for path, limitConfig := range limits {
			if c.Request.URL.Path == path {
				config = limitConfig
				limiter = limiters[path]
				found = true
				break
			}
		}

		// If no specific limit found, use default
		if !found {
			defaultConfig := DefaultRateLimitConfig()
			// Use a default limiter for all other requests
			if _, ok := limiters["default"]; !ok {
				limiters["default"] = NewRateLimiter(defaultConfig.RequestsPerMinute, defaultConfig.Burst)
			}
			config = defaultConfig
			limiter = limiters["default"]
		}

		// Check if we should skip rate limiting
		if config.SkipRateLimit != nil && config.SkipRateLimit(clientIP) {
			c.Next()
			return
		}

		// Check rate limit
		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded for %s. Maximum %d requests per minute.",
					c.Request.URL.Path, config.RequestsPerMinute),
				"code": 429,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
