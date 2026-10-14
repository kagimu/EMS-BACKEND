package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

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

type onboardingResetRequest struct {
	Confirmation string `json:"confirmation"`
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

	// 7. Verify the final state BEFORE committing; roll back if anything remains.
	verifyTables := []string{
		"users", "contacts", "user_passwords", "password_resets",
		"roles", "role_permissions", "user_roles", "scopes",
		"entities", "entity_users", "active_user_entities",
		"ems_licences", "ems_installations", "ems_licence_events",
		"ems_stakeholders", "ems_user_assignments", "ems_decision_authorities",
		"ems_approval_requests", "ems_approval_steps", "ems_approval_assignments",
		"ems_notifications", "ems_audit_logs",
	}
	for _, table := range verifyTables {
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
