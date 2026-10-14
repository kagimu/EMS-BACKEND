package middleware

import (
	"net/http"
	"sync"
	"time"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

const (
	// activationLimitWindow is the fixed window per authenticated user.
	activationLimitWindow = time.Minute
	// activationLimitMax is deliberately generous: normal installations
	// activate a handful of times; the guard mitigates abuse, not usage.
	activationLimitMax = 30
	// activationLimitSweep bounds the in-memory table; pruning is lazy and
	// happens only on request, never in a goroutine.
	activationLimitSweep = 2048
)

type activationWindow struct {
	start time.Time
	count int
}

var (
	activationLimitMu    sync.Mutex
	activationLimitTable = make(map[string]activationWindow)
)

// RequireActivationRateLimit is the in-process, request-driven abuse guard for
// the activation endpoint (Foundation F §22, Increment 5 hardening).
//
// It uses a fixed window per authenticated user and prunes expired entries
// lazily on access. No scheduler, worker, queue, ticker or background process
// is introduced, and it stores no licence state: licence validity remains
// derived at request time.
func RequireActivationRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}
		user, ok := value.(authModels.User)
		if !ok || user.Number == "" {
			c.Next()
			return
		}

		now := time.Now()
		activationLimitMu.Lock()
		entry, found := activationLimitTable[user.Number]
		if !found || now.Sub(entry.start) >= activationLimitWindow {
			entry = activationWindow{start: now}
		}
		entry.count++
		activationLimitTable[user.Number] = entry
		if len(activationLimitTable) > activationLimitSweep {
			for key, candidate := range activationLimitTable {
				if now.Sub(candidate.start) >= activationLimitWindow {
					delete(activationLimitTable, key)
				}
			}
		}
		activationLimitMu.Unlock()

		if entry.count > activationLimitMax {
			c.JSON(http.StatusTooManyRequests,
				commons.Feedback("Too many activation attempts; please try again later"))
			c.Abort()
			return
		}

		c.Next()
	}
}
