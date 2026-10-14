package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
)

const onboardingResetDatabase = "db_wholesome_onboarding_test"

// TestDevResetRouteIsOptIn verifies that the developer reset route only exists
// when DEV_RESET_ENABLED=true, so no production route is created. This test does
// not need a database.
func TestDevResetRouteIsOptIn(t *testing.T) {
	routePath := "POST /main/dev/reset-onboarding"

	t.Setenv("DEV_RESET_ENABLED", "false")
	if routesAreRegistered(routePath) {
		t.Fatalf("dev reset route registered while DEV_RESET_ENABLED=false")
	}

	t.Setenv("DEV_RESET_ENABLED", "true")
	if !routesAreRegistered(routePath) {
		t.Fatalf("dev reset route not registered while DEV_RESET_ENABLED=true")
	}
}

func routesAreRegistered(want string) bool {
	for _, route := range routes.SetupRoutes().Routes() {
		if route.Method+" "+route.Path == want {
			return true
		}
	}
	return false
}

// TestOnboardingResetDedicatedDatabase exercises the ACTUAL reset behaviour
// against the dedicated disposable onboarding database only. It hard-skips
// (before any destructive action) unless the connected database is exactly the
// designated onboarding database, so it can never reset a protected database.
func TestOnboardingResetDedicatedDatabase(t *testing.T) {
	db := testutil.RequireTestDB(t)

	// HARD SAFETY GUARD — skip (never reset) unless we are on the dedicated DB.
	var current string
	if err := db.Raw("SELECT DATABASE()").Scan(&current).Error; err != nil {
		t.Fatalf("could not read the connected database: %v", err)
	}
	if current != onboardingResetDatabase {
		t.Skipf("destructive reset test skipped: connected database is %q, not %q", current, onboardingResetDatabase)
	}

	t.Setenv("JWT_SECRET", "onboarding-reset-test-secret")
	t.Setenv("DEV_RESET_ENABLED", "true")
	t.Setenv("DEV_ONBOARDING_DB_NAME", onboardingResetDatabase)
	// Keep registration fast: make any email dial fail immediately (the sender
	// swallows dial errors, so registration still succeeds).
	t.Setenv("MAIL_HOST", "127.0.0.1")
	t.Setenv("MAIL_PORT", "1")
	t.Setenv("MAIL_FROM", "onboarding-test@example.invalid")

	const adminNumber = "onboarding-reset-admin"
	const plainNumber = "onboarding-reset-plain"

	cleanup := func() {
		db.Unscoped().Where("user_number IN ?", []string{adminNumber, plainNumber}).Delete(&authModel.Contact{})
		db.Unscoped().Where("number IN ?", []string{adminNumber, plainNumber}).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, u := range []authModel.User{
		{Number: adminNumber, FirstName: "Reset", SecondName: "Admin", Gender: "M", Type: "Admin", Status: "Active"},
		{Number: plainNumber, FirstName: "Reset", SecondName: "Plain", Gender: "M", Type: "User", Status: "Active"},
	} {
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("seeding user %s: %v", u.Number, err)
		}
	}

	tokenFor := func(number string) string {
		tok, err := authServices.JWTAuthService().GenerateToken(number)
		if err != nil {
			t.Fatalf("generating token: %v", err)
		}
		return tok
	}

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	call := func(token, confirmation string) (int, []byte) {
		t.Helper()
		body, _ := json.Marshal(map[string]string{"confirmation": confirmation})
		req, err := http.NewRequest(http.MethodPost, server.URL+"/main/dev/reset-onboarding", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("building request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, raw
	}

	// --- Unauthenticated is refused (JWT still enforced). ---
	if status, _ := call("", resetPhrase); status != http.StatusUnauthorized {
		t.Fatalf("unauthenticated reset = %d, want 401", status)
	}

	// --- Non-admin is refused. ---
	if status, _ := call(tokenFor(plainNumber), resetPhrase); status != http.StatusForbidden {
		t.Fatalf("non-admin reset = %d, want 403", status)
	}

	// --- Incorrect confirmation is refused. ---
	if status, _ := call(tokenFor(adminNumber), "wrong phrase"); status != http.StatusBadRequest {
		t.Fatalf("wrong confirmation = %d, want 400", status)
	}

	// --- Database mismatch is refused (protected/mismatched target). ---
	t.Setenv("DEV_ONBOARDING_DB_NAME", "some_other_database")
	if status, body := call(tokenFor(adminNumber), resetPhrase); status != http.StatusConflict {
		t.Fatalf("mismatched target = %d, want 409 (%s)", status, string(body))
	}
	t.Setenv("DEV_ONBOARDING_DB_NAME", onboardingResetDatabase)

	// --- Disabled feature is refused even though the route exists. ---
	t.Setenv("DEV_RESET_ENABLED", "false")
	if status, _ := call(tokenFor(adminNumber), resetPhrase); status != http.StatusNotFound {
		t.Fatalf("disabled feature = %d, want 404", status)
	}
	t.Setenv("DEV_RESET_ENABLED", "true")

	// Capture a sentinel table count before the destructive reset.
	var tablesBefore int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE()").Scan(&tablesBefore).Error; err != nil {
		t.Fatalf("could not count tables: %v", err)
	}

	// --- Authorised reset succeeds. ---
	status, body := call(tokenFor(adminNumber), resetPhrase)
	if status != http.StatusOK {
		t.Fatalf("authorised reset = %d, want 200 (%s)", status, string(body))
	}

	// --- Application data removed. ---
	for _, table := range []string{"users", "contacts", "entities", "entity_users", "active_user_entities", "ems_licences"} {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM `" + table + "`").Scan(&count).Error; err != nil {
			t.Fatalf("counting %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("after reset, %s still has %d row(s)", table, count)
		}
	}

	// --- EVERY enumerated base table is empty (full verification, B2). ---
	var allTables []string
	if err := db.Raw("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'").Scan(&allTables).Error; err != nil {
		t.Fatalf("could not enumerate base tables: %v", err)
	}
	if len(allTables) == 0 {
		t.Fatal("no base tables found to verify")
	}
	for _, table := range allTables {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM `" + table + "`").Scan(&count).Error; err != nil {
			t.Fatalf("counting %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("after reset, table %s still has %d row(s)", table, count)
		}
	}

	// --- Schema preserved. ---
	var tablesAfter int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE()").Scan(&tablesAfter).Error; err != nil {
		t.Fatalf("could not count tables after reset: %v", err)
	}
	if tablesAfter != tablesBefore || tablesAfter == 0 {
		t.Fatalf("schema not preserved: tables before=%d after=%d", tablesBefore, tablesAfter)
	}
	for _, table := range []string{"users", "contacts", "entities", "ems_licences", "ems_modules"} {
		var exists int64
		if err := db.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", table).Scan(&exists).Error; err != nil {
			t.Fatalf("checking table %s: %v", table, err)
		}
		if exists != 1 {
			t.Fatalf("table %s missing after reset", table)
		}
	}

	// --- A new first user can register through the existing flow. ---
	regBody, _ := json.Marshal(map[string]string{
		"first_name": "First", "second_name": "Admin", "last_name": "Test",
		"gender": "Male", "email": "onboarding-first-admin@example.test", "phone": "0700000001",
	})
	regResp, err := client.Post(server.URL+"/auth/register", "application/json", bytes.NewReader(regBody))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	regRaw, _ := io.ReadAll(regResp.Body)
	regResp.Body.Close()
	if regResp.StatusCode != http.StatusOK {
		t.Fatalf("post-reset registration = %d, want 200 (%s)", regResp.StatusCode, string(regRaw))
	}
	var regUser struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	_ = json.Unmarshal(regRaw, &regUser)

	// Clean up the freshly registered user so the onboarding DB is left empty.
	if regUser.Data.Number != "" {
		db.Unscoped().Where("user_number = ?", regUser.Data.Number).Delete(&authModel.Contact{})
		db.Unscoped().Where("user_number = ?", regUser.Data.Number).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("user_number = ?", regUser.Data.Number).Delete(&authModel.PasswordReset{})
		db.Unscoped().Where("number = ?", regUser.Data.Number).Delete(&authModel.User{})
	}

	// --- Education authentication remains enforced (unauthenticated → 401). ---
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/main/curriculum/all", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("education probe failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated education route = %d, want 401", resp.StatusCode)
	}
}

// resetPhrase keeps the exact phrase in one place for the test.
const resetPhrase = "RESET EMS TEST SYSTEM"
