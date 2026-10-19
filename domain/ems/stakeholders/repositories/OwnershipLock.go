package repositories

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"gorm.io/gorm"
)

// ownershipLockTimeoutSeconds bounds how long an ownership mutation waits for
// the per-organisation named lock before failing without changing anything.
const ownershipLockTimeoutSeconds = 5

// ErrOwnershipChangeBusy is returned when the per-organisation ownership lock
// could not be acquired. The caller must not have modified ownership.
var ErrOwnershipChangeBusy = errors.New("another ownership change is in progress for this organisation; please retry")

// OwnershipLockName returns the deterministic MySQL named-lock name used to
// serialise ownership mutations for one organisation. SHA-1 keeps the name
// below MySQL's 64-character limit regardless of the organisation identifier.
func OwnershipLockName(organisationNumber string) string {
	sum := sha1.Sum([]byte(organisationNumber))
	return fmt.Sprintf("ems_ownership_%x", sum)
}

// acquireOwnershipLock takes the per-organisation named lock on a single pinned
// connection (MySQL named locks are connection-scoped and are not released by
// COMMIT/ROLLBACK). It returns an idempotent release function and whether the
// lock was obtained within the timeout.
func acquireOwnershipLock(organisationNumber string) (func(), bool) {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return func() {}, false
	}
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return func() {}, false
	}

	lockName := OwnershipLockName(organisationNumber)
	var got sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, ?)", lockName, ownershipLockTimeoutSeconds).Scan(&got); err != nil || !got.Valid || got.Int64 != 1 {
		conn.Close()
		return func() {}, false
	}

	released := false
	return func() {
		if released {
			return
		}
		released = true
		_, _ = conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK(?)", lockName)
		conn.Close()
	}, true
}

// WithOwnershipTransaction serialises ownership mutations for one organisation
// and runs fn inside a transaction.
//
// The lock is acquired BEFORE the transaction begins and released on every exit
// path (success, error, rollback or panic). If the lock cannot be acquired, fn
// is not called, so no ownership change can occur.
func WithOwnershipTransaction(organisationNumber string, fn func(tx *gorm.DB) error) error {
	if organisationNumber == "" {
		return ErrOwnershipChangeBusy
	}

	release, acquired := acquireOwnershipLock(organisationNumber)
	if !acquired {
		return ErrOwnershipChangeBusy
	}
	defer release()

	return config.DB.Transaction(fn)
}
