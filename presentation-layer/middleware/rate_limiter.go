package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements guideline 15 - Rate Limiting & Throttling
type RateLimiter struct {
	clients map[string]*ClientLimiter
	mu      sync.RWMutex
	rate    int
	burst   int
}

type ClientLimiter struct {
	tokens    int
	lastRefill time.Time
	mu        sync.Mutex
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		rate:    rate,
		burst:   burst,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := rl.getClientID(c)
		
		if !rl.allow(clientID) {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.rate))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
			c.Abort()
			return
		}

		remaining := rl.getRemaining(clientID)
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
		
		c.Next()
	}
}

func (rl *RateLimiter) getClientID(c *gin.Context) string {
	// Use tenant ID if available, otherwise fall back to IP
	if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
		return "tenant:" + tenantID
	}
	return "ip:" + c.ClientIP()
}

func (rl *RateLimiter) allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[clientID]
	if !exists {
		client = &ClientLimiter{
			tokens:    rl.burst,
			lastRefill: time.Now(),
		}
		rl.clients[clientID] = client
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(client.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * rl.rate / 60 // rate per minute

	if tokensToAdd > 0 {
		client.tokens += tokensToAdd
		if client.tokens > rl.burst {
			client.tokens = rl.burst
		}
		client.lastRefill = now
	}

	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) getRemaining(clientID string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if client, exists := rl.clients[clientID]; exists {
		client.mu.Lock()
		defer client.mu.Unlock()
		return client.tokens
	}
	return rl.burst
}