package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	tokens         float64
	capacity       float64
	refillRate     float64
	lastRefillTime time.Time
	mu             sync.Mutex
}

type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

func (rl *RateLimiter) getBucket(key string, capacity, refillRate float64) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:         capacity,
			capacity:       capacity,
			refillRate:     refillRate,
			lastRefillTime: time.Now(),
		}
		rl.buckets[key] = bucket
	}
	return bucket
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+(duration*tb.refillRate))
	tb.lastRefillTime = now
}

func (tb *TokenBucket) tryConsume(tokens float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware creates a rate limiting middleware
// requestsPerMinute: number of requests allowed per minute
func RateLimitMiddleware(requestsPerMinute float64) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		// Get client identifier (IP or user ID if authenticated)
		clientID := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			clientID = userID.(string)
		}

		// Higher rate limit for authenticated users
		rateLimit := requestsPerMinute
		if _, exists := c.Get("user_id"); exists {
			rateLimit *= 5 // 5x higher limit for authenticated users
		}

		bucket := limiter.getBucket(clientID, rateLimit, rateLimit/60.0)
		if !bucket.tryConsume(1.0) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", rateLimit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%.0f", bucket.tokens))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", bucket.lastRefillTime.Add(time.Minute).Unix()))

		c.Next()
	}
}
