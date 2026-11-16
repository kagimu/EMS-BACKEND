package services

import (
	"context"
	"testing"
	"time"
)

func TestReconciliationWorkerStopsOnContextCancel(t *testing.T) {
	t.Setenv("PAYMENTS_RECONCILIATION_ENABLED", "true")
	ctx, cancel := context.WithCancel(context.Background())
	StartReconciliationWorker(ctx)
	defer cancel()

	deadline := time.Now().Add(2 * time.Second)
	for !WorkerStatus().Running && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if !WorkerStatus().Running {
		t.Fatal("reconciliation worker did not start")
	}
	cancel()
	deadline = time.Now().Add(2 * time.Second)
	for WorkerStatus().Running && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if WorkerStatus().Running {
		t.Fatal("reconciliation worker did not stop after context cancellation")
	}
}
