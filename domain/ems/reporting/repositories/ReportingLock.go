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

const reportingLockTimeoutSeconds = 5

// ErrReportingChangeBusy is returned when the per-organisation reporting lock
// could not be acquired. The caller must not have modified any reporting line.
var ErrReportingChangeBusy = errors.New("another reporting change is in progress for this organisation; please retry")

// ReportingLockName returns the deterministic MySQL named-lock name used to
// serialise reporting-line mutations for one organisation.
func ReportingLockName(organisationNumber string) string {
	sum := sha1.Sum([]byte(organisationNumber))
	return fmt.Sprintf("ems_reporting_%x", sum)
}

func acquireReportingLock(organisationNumber string) (func(), bool) {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return func() {}, false
	}
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return func() {}, false
	}

	lockName := ReportingLockName(organisationNumber)
	var got sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, ?)", lockName, reportingLockTimeoutSeconds).Scan(&got); err != nil || !got.Valid || got.Int64 != 1 {
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

// WithReportingTransaction serialises reporting changes for one organisation and
// runs fn in a transaction. The lock is acquired before the transaction and
// released on every exit path. If the lock cannot be acquired, fn is not called.
func WithReportingTransaction(organisationNumber string, fn func(tx *gorm.DB) error) error {
	if organisationNumber == "" {
		return ErrReportingChangeBusy
	}

	release, acquired := acquireReportingLock(organisationNumber)
	if !acquired {
		return ErrReportingChangeBusy
	}
	defer release()

	return config.DB.Transaction(fn)
}
