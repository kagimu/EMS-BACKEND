package controllers

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/controllers"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

// resetLockTimeoutSeconds bounds how long the reset waits for the bootstrap
// named lock before aborting with a conflict.
const resetLockTimeoutSeconds = 10

// acquireResetLock takes the SAME installation bootstrap named lock used by
// first-owner registration, on one pinned connection, so a reset and a
// registration can never interleave. MySQL named locks are connection-scoped
// and are NOT released by COMMIT/ROLLBACK, so the lock is always released
// explicitly on the connection that acquired it. Returns false when the lock
// could not be acquired within the timeout.
func acquireResetLock() (func(), bool) {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return func() {}, false
	}
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return func() {}, false
	}

	var got sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, ?)", authController.BootstrapLockName, resetLockTimeoutSeconds).Scan(&got); err != nil || !got.Valid || got.Int64 != 1 {
		conn.Close()
		return func() {}, false
	}

	released := false
	return func() {
		if released {
			return
		}
		released = true
		_, _ = conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK(?)", authController.BootstrapLockName)
		conn.Close()
	}, true
}

// Developer-only onboarding reset.
//
// This is a bounded, defence-in-depth developer tool that restores the DEDICATED
// disposable onboarding database to a genuinely fresh-install state. It exists
// so the first-time registration journey (including the first-user scenario) can
// be repeated from a clean slate.
//
// Safety properties (all enforced server-side, independent of the frontend):
//   - The whole feature is disabled unless DEV_RESET_ENABLED=true. When it is
//     not set the route is not even registered, so no production route exists.
//   - JWT authentication is NOT bypassed: the route lives under the authenticated
//     /main group, so a valid bearer token and an existing user are required.
//   - The caller must be an administrator.
//   - The caller must send the exact confirmation phrase.
//   - The ACTUAL connected database is resolved with SELECT DATABASE() and must
//     equal the configured disposable onboarding database name. Protected
//     databases are refused even if misconfigured.
//   - No database name is ever accepted from the client.
//   - The clear is performed in a single transaction with foreign-key checks
//     disabled; the final state is verified before the commit and the whole
//     transaction is rolled back if verification fails, so success is never
//     reported for a partial reset.
const (
	// OnboardingResetConfirmation is the exact phrase an operator must send.
	OnboardingResetConfirmation = "RESET EMS TEST SYSTEM"
	// DefaultOnboardingDatabase is the dedicated disposable onboarding database.
	DefaultOnboardingDatabase = "db_wholesome_onboarding_test"

	devResetEnabledEnv    = "DEV_RESET_ENABLED"
	onboardingDatabaseEnv = "DEV_ONBOARDING_DB_NAME"

	// DefaultDevResetPassword is the password required by the sign-in-free
	// testing reset. Overridable with DEV_RESET_PASSWORD. This is a testing
	// convenience gate, NOT a security boundary: it is only reachable when
	// DEV_RESET_ENABLED=true, and it can only ever clear the disposable
	// onboarding database (protected databases are always refused below).
	DefaultDevResetPassword = "12345"
	devResetPasswordEnv     = "DEV_RESET_PASSWORD"
)

// protectedDatabases must never be reset by this feature, even if misconfigured.
var protectedDatabases = map[string]struct{}{
	"db_wholesome":         {},
	"db_wholesome_e1_test": {},
	"db_wholesome_e2_test": {},
}

// DevResetEnabled reports whether the developer reset feature is enabled. It is
// deliberately opt-in: absent means disabled.
func DevResetEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(devResetEnabledEnv)), "true")
}

// OnboardingDatabaseName returns the configured disposable database name.
func OnboardingDatabaseName() string {
	if name := strings.TrimSpace(os.Getenv(onboardingDatabaseEnv)); name != "" {
		return name
	}
	return DefaultOnboardingDatabase
}

// DevResetPassword returns the password required by the sign-in-free testing
// reset. Absent configuration falls back to DefaultDevResetPassword.
func DevResetPassword() string {
	if password := os.Getenv(devResetPasswordEnv); password != "" {
		return password
	}
	return DefaultDevResetPassword
}

type onboardingResetRequest struct {
	Confirmation string `json:"confirmation"`
	// Password is only used by the sign-in-free testing route; it is ignored by
	// the authenticated administrator route.
	Password string `json:"password"`
}

// Reset-target refusal reasons (exported for tests).
var (
	ErrResetTargetUnknown   = errors.New("the connected database could not be identified")
	ErrResetTargetProtected = errors.New("refusing to reset a protected database")
	ErrResetTargetMismatch  = errors.New("refusing to reset a database other than the configured onboarding database")
)

// evaluateResetTarget decides whether the ACTUAL connected database (from
// SELECT DATABASE()) may be reset. It is deliberately pure so the refusal rules
// can be tested without ever targeting a protected database.
func evaluateResetTarget(current, expected string) error {
	current = strings.TrimSpace(current)
	if current == "" {
		return ErrResetTargetUnknown
	}
	if _, isProtected := protectedDatabases[strings.ToLower(current)]; isProtected {
		return ErrResetTargetProtected
	}
	if current != expected {
		return ErrResetTargetMismatch
	}
	return nil
}

// ResetOnboardingDatabase clears all application data from the dedicated
// disposable onboarding database. See the package-level comment for the full
// safety contract.
func ResetOnboardingDatabase(ctx *gin.Context) {
	// 1. Environment gate (defence in depth; the route is also gated).
	if !DevResetEnabled() {
		ctx.JSON(http.StatusNotFound, commons.Feedback("Not found"))
		return
	}

	// 2. Authentication + developer authorisation (admin caller only).
	value, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	user, ok := value.(authModels.User)
	if !ok || user.Number == "" {
		ctx.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	if !strings.EqualFold(user.Type, "Admin") {
		ctx.JSON(http.StatusForbidden, commons.Feedback("The developer reset requires an administrator account"))
		return
	}

	// 3. Exact confirmation phrase.
	var req onboardingResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Confirmation phrase is required"))
		return
	}
	if req.Confirmation != OnboardingResetConfirmation {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Confirmation phrase does not match"))
		return
	}

	performOnboardingReset(ctx)
}

// ResetOnboardingDatabaseWithPassword clears the dedicated disposable onboarding
// database WITHOUT a signed-in session, for the testing phase only. It is
// registered outside the authenticated /main group and is gated by
// DEV_RESET_ENABLED plus a password (DEV_RESET_PASSWORD, defaulting to
// DefaultDevResetPassword).
//
// Every other server-side safety property still applies unchanged: the exact
// confirmation phrase, the disposable-database identity check (protected
// databases are always refused), the bootstrap lock, and the transactional
// clear-and-verify. Only the caller-identity step (JWT + administrator) is
// replaced by the password, because an operator who has forgotten the admin
// password must still be able to repeat first-time setup.
func ResetOnboardingDatabaseWithPassword(ctx *gin.Context) {
	// 1. Environment gate (defence in depth; the route is also gated).
	if !DevResetEnabled() {
		ctx.JSON(http.StatusNotFound, commons.Feedback("Not found"))
		return
	}

	// 2. Exact confirmation phrase + reset password.
	var req onboardingResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Confirmation phrase is required"))
		return
	}
	if req.Confirmation != OnboardingResetConfirmation {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Confirmation phrase does not match"))
		return
	}
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(DevResetPassword())) != 1 {
		ctx.JSON(http.StatusForbidden, commons.Feedback("The developer reset password is incorrect"))
		return
	}

	performOnboardingReset(ctx)
}

// performOnboardingReset carries out the shared, caller-independent part of the
// reset: verify the ACTUAL connected database, take the bootstrap lock, clear
// every table in a single transaction, verify the result and commit.
func performOnboardingReset(ctx *gin.Context) {
	// 4. Verify the ACTUAL connected database (never trust a client value).
	var current string
	if err := config.DB.Raw("SELECT DATABASE()").Scan(&current).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Could not identify the connected database"))
		return
	}
	expected := OnboardingDatabaseName()
	switch err := evaluateResetTarget(current, expected); {
	case errors.Is(err, ErrResetTargetUnknown):
		ctx.JSON(http.StatusConflict, commons.Feedback("The connected database could not be identified; refusing to reset"))
		return
	case errors.Is(err, ErrResetTargetProtected):
		ctx.JSON(http.StatusConflict, commons.Feedback("Refusing to reset a protected database"))
		return
	case errors.Is(err, ErrResetTargetMismatch):
		ctx.JSON(http.StatusConflict, commons.Feedback(fmt.Sprintf("Refusing to reset %q: this feature only operates on %q", strings.TrimSpace(current), expected)))
		return
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Could not verify the reset target"))
		return
	}

	// 4b. Coordinate with first-owner registration using the SAME bootstrap
	// named lock, held for the critical reset period and released on every exit
	// path via defer. If the lock cannot be acquired, abort without resetting.
	releaseBootstrapLock, locked := acquireResetLock()
	if !locked {
		ctx.JSON(http.StatusConflict, commons.Feedback("Another installation operation (registration or reset) is in progress. Please retry shortly."))
		return
	}
	defer releaseBootstrapLock()

	// 5. Enumerate every application table in the connected database.
	var tables []string
	if err := config.DB.Raw(
		"SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'",
	).Scan(&tables).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Could not enumerate the database tables"))
		return
	}
	if len(tables) == 0 {
		ctx.JSON(http.StatusConflict, commons.Feedback("No tables were found; refusing to report a partial reset"))
		return
	}

	// 6. Clear every table inside one transaction with foreign-key checks off.
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Could not start the reset transaction"))
		return
	}

	fail := func(status int, message string) {
		// Always restore FK checks on the transaction connection before unwinding,
		// so the pooled connection is not left with checks disabled.
		tx.Exec("SET FOREIGN_KEY_CHECKS=1")
		tx.Rollback()
		ctx.JSON(status, commons.Feedback(message))
	}

	if err := tx.Exec("SET FOREIGN_KEY_CHECKS=0").Error; err != nil {
		fail(http.StatusInternalServerError, "Could not disable foreign-key checks")
		return
	}

	for _, table := range tables {
		if err := tx.Exec("DELETE FROM `" + table + "`").Error; err != nil {
			fail(http.StatusInternalServerError, "Reset failed while clearing table "+table)
			return
		}
	}

	// 7. Verify EVERY table that was cleared is now empty BEFORE committing;
	// roll back if anything remains. The same slice that was deleted from is
	// checked, so a partial reset anywhere — including the commercial and
	// custom-field tables — is caught rather than reported as success.
	for _, table := range tables {
		var count int64
		if err := tx.Raw("SELECT COUNT(*) FROM `" + table + "`").Scan(&count).Error; err != nil {
			fail(http.StatusInternalServerError, "Reset verification failed for table "+table)
			return
		}
		if count != 0 {
			fail(http.StatusConflict, fmt.Sprintf("Reset verification failed: %s still contains %d row(s)", table, count))
			return
		}
	}

	// 8. Re-enable FK checks on the same connection, then commit.
	if err := tx.Exec("SET FOREIGN_KEY_CHECKS=1").Error; err != nil {
		fail(http.StatusInternalServerError, "Could not re-enable foreign-key checks")
		return
	}
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to commit the reset"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "The onboarding test database has been reset to a fresh state",
		"database":       strings.TrimSpace(current),
		"tables_cleared": len(tables),
	})
}
