package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

// TestSetupProgressAndRegisterConflict exercises the public GET /setup/progress
// endpoint (read-only facts + next_action, minimum disclosure) and the safe 409
// handling when an email already has an account. It writes only to the
// dedicated onboarding database and skips otherwise.
func TestSetupProgressAndRegisterConflict(t *testing.T) {
	db := requireOnboardingDB(t)
	t.Setenv("JWT_SECRET", "setup-progress-test-secret")
	// Fail the SMTP dial fast (the sender swallows dial errors, so registration
	// still succeeds; email delivery is now post-commit).
	t.Setenv("MAIL_HOST", "127.0.0.1")
	t.Setenv("MAIL_PORT", "1")
	t.Setenv("MAIL_FROM", "setup-progress@example.invalid")

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	post := func(path string, payload any) (int, []byte) {
		t.Helper()
		raw, _ := json.Marshal(payload)
		resp, err := client.Post(server.URL+path, "application/json", bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("POST %s: %v", path, err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	// ---- 1. progress endpoint: shape, booleans, minimum disclosure ----
	resp, err := client.Get(server.URL + "/setup/progress")
	if err != nil {
		t.Fatalf("GET /setup/progress: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /setup/progress = %d, want 200", resp.StatusCode)
	}
	rawBody, _ := io.ReadAll(resp.Body)

	var progress struct {
		State      string          `json:"state"`
		NextAction string          `json:"next_action"`
		Steps      map[string]bool `json:"steps"`
	}
	if err := json.Unmarshal(rawBody, &progress); err != nil {
		t.Fatalf("decode progress: %v (%s)", err, string(rawBody))
	}

	// Minimum disclosure: only state, steps (6 booleans), next_action.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(rawBody, &top); err != nil {
		t.Fatalf("decode top: %v", err)
	}
	if len(top) != 3 {
		t.Fatalf("progress exposed %d top-level fields, want 3: %s", len(top), string(rawBody))
	}
	for _, key := range []string{"state", "steps", "next_action"} {
		if _, ok := top[key]; !ok {
			t.Fatalf("progress missing %q: %s", key, string(rawBody))
		}
	}
	for _, key := range []string{"users_exist", "admin_active", "owner_exists", "organisation_exists", "subscription_active", "education_entitled"} {
		if _, ok := progress.Steps[key]; !ok {
			t.Fatalf("progress steps missing %q: %s", key, string(rawBody))
		}
	}
	allowed := map[string]bool{
		"commercial_setup": true, "resume_activation": true, "complete_setup": true,
		"activate_subscription": true, "ready": true,
	}
	if !allowed[progress.NextAction] {
		t.Fatalf("unexpected next_action %q", progress.NextAction)
	}

	// Cross-check the reported facts against the database.
	var users, owners, orgs int64
	db.Model(&authModel.User{}).Count(&users)
	db.Raw("SELECT COUNT(*) FROM ems_stakeholders WHERE stakeholder_type='owner' AND status='ACTIVE'").Scan(&owners)
	db.Raw("SELECT COUNT(*) FROM entities").Scan(&orgs)
	if progress.Steps["users_exist"] != (users > 0) {
		t.Fatalf("users_exist=%v but db users=%d", progress.Steps["users_exist"], users)
	}
	if progress.Steps["owner_exists"] != (owners > 0) {
		t.Fatalf("owner_exists=%v but db owners=%d", progress.Steps["owner_exists"], owners)
	}
	if progress.Steps["organisation_exists"] != (orgs > 0) {
		t.Fatalf("organisation_exists=%v but db orgs=%d", progress.Steps["organisation_exists"], orgs)
	}
	if progress.State == "initialized" && !progress.Steps["owner_exists"] {
		t.Fatal("state initialized requires owner_exists")
	}

	// ---- 2. duplicate email → 409 (safe) ----
	suffix := time.Now().UnixNano() % 1000000000
	email := fmt.Sprintf("prog-%d@example.test", suffix)
	phone := fmt.Sprintf("077%09d", suffix)
	payload := map[string]string{
		"first_name": "Prog", "second_name": "Test", "last_name": "User",
		"gender": "Male", "email": email, "phone": phone,
	}
	status, regBody := post("/auth/register", payload)
	if status != http.StatusOK {
		t.Fatalf("first registration = %d, want 200 (%s)", status, string(regBody))
	}
	var created struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
		TokenInfo string `json:"tokenInfo"`
	}
	_ = json.Unmarshal(regBody, &created)
	t.Cleanup(func() {
		if created.Data.Number != "" {
			db.Unscoped().Where("user_number = ?", created.Data.Number).Delete(&authModel.Contact{})
			db.Unscoped().Where("user_number = ?", created.Data.Number).Delete(&authModel.PasswordReset{})
			db.Unscoped().Where("number = ?", created.Data.Number).Delete(&authModel.User{})
		}
	})
	if created.Data.Number == "" {
		t.Fatalf("registration returned no user number: %s", string(regBody))
	}
	if created.TokenInfo == "" {
		t.Fatalf("registration returned no activation token: %s", string(regBody))
	}

	status, dupBody := post("/auth/register", payload)
	if status != http.StatusConflict {
		t.Fatalf("duplicate email = %d, want 409 (%s)", status, string(dupBody))
	}
	var dup struct {
		NextAction string `json:"next_action"`
	}
	_ = json.Unmarshal(dupBody, &dup)
	if dup.NextAction != "existing_account" {
		t.Fatalf("duplicate email next_action = %q, want existing_account", dup.NextAction)
	}

	// No second account was created for that email.
	var contactCount int64
	db.Model(&authModel.Contact{}).Where("contact = ? AND contact_type = ?", email, "Email").Count(&contactCount)
	// Note: registration stores contact_type "Email" (capitalised).
	if contactCount != 1 {
		var lower int64
		db.Model(&authModel.Contact{}).Where("contact = ? AND contact_type = ?", email, "email").Count(&lower)
		if contactCount+lower != 1 {
			t.Fatalf("email has %d contacts, want 1 (no duplicate account)", contactCount+lower)
		}
	}
}
