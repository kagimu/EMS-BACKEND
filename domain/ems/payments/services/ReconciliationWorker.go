package services

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	paymentEvents "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/events"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
)

type ReconciliationWorkerStatus struct {
	Enabled          bool      `json:"enabled"`
	Running          bool      `json:"running"`
	IntervalSeconds  int       `json:"interval_seconds"`
	MaxAttempts      int       `json:"max_attempts"`
	LastRunAt        time.Time `json:"last_run_at,omitempty"`
	LastSuccessAt    time.Time `json:"last_success_at,omitempty"`
	LastError        string    `json:"last_error,omitempty"`
	LastScannedCount int       `json:"last_scanned_count"`
	LastUpdatedCount int       `json:"last_updated_count"`
}

var (
	reconciliationRunning int32
	reconciliationLastRun atomic.Value
	reconciliationLastOK  atomic.Value
	reconciliationLastErr atomic.Value
	reconciliationScanned atomic.Int64
	reconciliationUpdated atomic.Int64
)

func reconciliationEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("PAYMENTS_RECONCILIATION_ENABLED")))
	return value == "true" || value == "1" || value == "yes"
}

func reconciliationInterval() time.Duration {
	value := strings.TrimSpace(os.Getenv("PAYMENTS_RECONCILIATION_INTERVAL_SECONDS"))
	if value == "" {
		return 2 * time.Minute
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 10 {
		return 2 * time.Minute
	}
	return time.Duration(seconds) * time.Second
}

func reconciliationMaxAttempts() int {
	value := strings.TrimSpace(os.Getenv("PAYMENTS_RECONCILIATION_MAX_ATTEMPTS"))
	if value == "" {
		return 12
	}
	attempts, err := strconv.Atoi(value)
	if err != nil || attempts < 1 {
		return 12
	}
	return attempts
}

func WorkerStatus() ReconciliationWorkerStatus {
	status := ReconciliationWorkerStatus{
		Enabled:          reconciliationEnabled(),
		Running:          atomic.LoadInt32(&reconciliationRunning) == 1,
		IntervalSeconds:  int(reconciliationInterval().Seconds()),
		MaxAttempts:      reconciliationMaxAttempts(),
		LastScannedCount: int(reconciliationScanned.Load()),
		LastUpdatedCount: int(reconciliationUpdated.Load()),
	}
	if value := reconciliationLastRun.Load(); value != nil {
		status.LastRunAt, _ = value.(time.Time)
	}
	if value := reconciliationLastOK.Load(); value != nil {
		status.LastSuccessAt, _ = value.(time.Time)
	}
	if value := reconciliationLastErr.Load(); value != nil {
		status.LastError, _ = value.(string)
	}
	return status
}

// StartReconciliationWorker starts the single in-process worker. It is safe to
// call once from main; duplicate starts are ignored. The worker is intentionally
// process-local; running multiple EMS instances requires an external lease or
// an operational decision to accept duplicate read-only status checks.
func StartReconciliationWorker(ctx context.Context) {
	if !reconciliationEnabled() {
		return
	}
	if !atomic.CompareAndSwapInt32(&reconciliationRunning, 0, 1) {
		return
	}
	go func() {
		defer atomic.StoreInt32(&reconciliationRunning, 0)
		interval := reconciliationInterval()
		log.Printf("payment reconciliation worker started (interval=%s, max_attempts=%d)", interval, reconciliationMaxAttempts())
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("payment reconciliation worker stopped")
				return
			case <-ticker.C:
				if _, err := ReconcileOnce(ctx); err != nil {
					log.Printf("payment reconciliation run failed: %v", err)
				}
			}
		}
	}()
}

// ReconcileOnce performs one bounded reconciliation pass. It never creates a
// new payment request and never converts provider downtime into failure.
func ReconcileOnce(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	reconciliationLastRun.Store(now)
	var payments []paymentModels.MobileMoneyPayment
	cutoff := now.Add(-reconciliationInterval())
	if err := config.DB.Where(
		"status IN ? AND check_count < ? AND (last_checked_at IS NULL OR last_checked_at < ?)",
		[]string{paymentModels.StatusInitiated, paymentModels.StatusPending, paymentModels.StatusUnknown},
		reconciliationMaxAttempts(), cutoff,
	).Order("last_checked_at asc").Limit(100).Find(&payments).Error; err != nil {
		reconciliationLastErr.Store(err.Error())
		return 0, err
	}
	scanned := 0
	updated := 0
	for _, payment := range payments {
		if ctx.Err() != nil {
			return updated, ctx.Err()
		}
		scanned++
		beforeStatus := payment.Status
		attempt := payment.CheckCount + 1
		_ = paymentEvents.Record(config.DB, payment.Number, paymentModels.EventReconciliationStarted,
			payment.Provider, payment.ProviderReference, "", map[string]interface{}{
				"attempt": attempt, "previous_status": beforeStatus,
			})
		refreshed, err := Refresh(payment.Number)
		result := paymentModels.ReconciliationPending
		discrepancy := ""
		providerAmount := int64(0)
		providerCurrency := ""
		if err != nil {
			result = paymentModels.ReconciliationProviderUnavailable
			discrepancy = err.Error()
		} else if refreshed != nil {
			providerAmount = refreshed.ProviderAmount
			if providerAmount == 0 {
				providerAmount = refreshed.AmountUGX
			}
			providerCurrency = refreshed.ProviderCurrency
			if providerCurrency == "" {
				providerCurrency = refreshed.Currency
			}
			if refreshed.Status == paymentModels.StatusSuccessful || refreshed.Status == paymentModels.StatusFailed ||
				refreshed.Status == paymentModels.StatusExpired || refreshed.Status == paymentModels.StatusCancelled {
				result = paymentModels.ReconciliationMatched
			}
			if strings.Contains(strings.ToLower(refreshed.FailureReason), "mismatch") {
				result = paymentModels.ReconciliationMismatch
				discrepancy = refreshed.FailureReason
			}
			if refreshed.Status != beforeStatus {
				updated++
			}
		}
		_ = config.DB.Create(&paymentModels.PaymentReconciliation{
			TransactionNumber: payment.Number,
			ProviderCode:      payment.Provider,
			ProviderReference: payment.ProviderReference,
			ProviderAmount:    providerAmount,
			EMSAmount:         payment.AmountUGX,
			ProviderCurrency:  providerCurrency,
			EMSCurrency:       payment.Currency,
			Result:            result,
			Discrepancy:       truncateReason(discrepancy),
			AttemptNumber:     attempt,
			ReconciledAt:      time.Now().UTC(),
		}).Error
		_ = paymentEvents.Record(config.DB, payment.Number, paymentModels.EventReconciliationCompleted,
			payment.Provider, payment.ProviderReference, "", map[string]interface{}{
				"attempt": attempt, "result": result, "discrepancy": truncateReason(discrepancy),
			})
	}
	reconciliationScanned.Store(int64(scanned))
	reconciliationUpdated.Store(int64(updated))
	reconciliationLastOK.Store(time.Now().UTC())
	reconciliationLastErr.Store("")
	return updated, nil
}
