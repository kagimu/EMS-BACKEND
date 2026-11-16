package controllers

import (
	"net/http"
	"sync"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

type rateWindow struct {
	count   int
	resetAt time.Time
}

// RateLimit is a small in-process guard for public payment initiation and
// callback endpoints. It complements, but does not replace, the transaction
// idempotency constraints in the payment service. Multi-instance deployments
// should front this with a shared gateway/rate-limit layer.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 20
	}
	if window <= 0 {
		window = time.Minute
	}
	var mu sync.Mutex
	buckets := map[string]rateWindow{}
	return func(c *gin.Context) {
		key := c.ClientIP() + "|" + c.FullPath()
		now := time.Now()
		mu.Lock()
		bucket := buckets[key]
		if bucket.resetAt.IsZero() || now.After(bucket.resetAt) {
			bucket = rateWindow{resetAt: now.Add(window)}
		}
		bucket.count++
		buckets[key] = bucket
		allowed := bucket.count <= limit
		mu.Unlock()
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, commons.Feedback("too many payment requests; try again later"))
			return
		}
		c.Next()
	}
}
