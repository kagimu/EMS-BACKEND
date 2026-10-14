package controllers

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
)

// TestBootstrapLockIsMutuallyExclusive verifies the concurrency guarantee behind
// the first-user determination (TEST4): while one registration holds the
// installation bootstrap lock, a second holder cannot also acquire it, so two
// simultaneous first-user registrations cannot both observe an empty
// installation. After release the lock becomes acquirable again.
//
// It runs only against the dedicated disposable onboarding database and skips
// otherwise (TEST11).
func TestBootstrapLockIsMutuallyExclusive(t *testing.T) {
	db := testutil.RequireTestDB(t)

	var current string
	if err := db.Raw("SELECT DATABASE()").Scan(&current).Error; err != nil {
		t.Fatalf("could not read the connected database: %v", err)
	}
	if current != "db_wholesome_onboarding_test" {
		t.Skipf("skipping: connected database is %q, not the dedicated onboarding database", current)
	}

	release, ok := acquireBootstrapLock()
	if !ok {
		t.Fatal("first acquireBootstrapLock() returned not-acquired")
	}

	sqlDB, err := db.DB()
	if err != nil {
		release()
		t.Fatalf("db handle: %v", err)
	}
	// A SECOND, independent connection — named locks are re-entrant for the
	// same session, so mutual exclusion must be proven from another connection.
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		release()
		t.Fatalf("second connection: %v", err)
	}
	defer conn.Close()

	var got sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, 0)", bootstrapLockName).Scan(&got); err != nil {
		release()
		t.Fatalf("second GET_LOCK: %v", err)
	}
	if got.Valid && got.Int64 == 1 {
		release()
		t.Fatal("bootstrap lock granted to a second holder while the first still held it")
	}

	// Release, then the lock must become acquirable.
	release()
	if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, 5)", bootstrapLockName).Scan(&got); err != nil {
		t.Fatalf("post-release GET_LOCK: %v", err)
	}
	if !got.Valid || got.Int64 != 1 {
		t.Fatal("bootstrap lock not acquirable after release")
	}
	_, _ = conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK(?)", bootstrapLockName)
}
