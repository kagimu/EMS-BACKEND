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
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"gorm.io/gorm"
)

const bootstrapTestDB = "db_wholesome_onboarding_test"

// requireOnboardingDB connects and positively verifies the effective database
// is exactly the dedicated disposable onboarding database. It skips (never
// writes) otherwise (TEST11).
func requireOnboardingDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.RequireTestDB(t)
	var current string
	if err := db.Raw("SELECT DATABASE()").Scan(&current).Error; err != nil {
		t.Fatalf("could not read the connected database: %v", err)
	}
	if current != bootstrapTestDB {
		t.Skipf("skipping: connected database is %q, not %q", current, bootstrapTestDB)
	}
	return db
}

func fubRequest(t *testing.T, client *http.Client, base, method, path, token string, payload any) (int, []byte) {
	t.Helper()
	var reader io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, base+path, reader)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

// TestFirstUserBootstrapLaterRegistrantIsNotAdmin covers TEST5: a registrant on
// a non-empty installation is an ordinary user with no roles.
func TestFirstUserBootstrapLaterRegistrantIsNotAdmin(t *testing.T) {
	db := requireOnboardingDB(t)
	t.Setenv("JWT_SECRET", "fub-test-secret")
	t.Setenv("MAIL_HOST", "127.0.0.1")
	t.Setenv("MAIL_PORT", "1")
	t.Setenv("MAIL_FROM", "fub@example.invalid")

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)

	suffix := time.Now().UnixNano() % 1000000000

	// Ensure the installation is non-empty so this test exercises the
	// "an earlier user already exists" rule (the ordinary-registrant case)
	// rather than the fresh-installation bootstrap rule.
	sentinelNumber := fmt.Sprintf("fub-sentinel-%d", suffix)
	if err := db.Create(&authModel.User{Number: sentinelNumber, FirstName: "Fub", SecondName: "Sentinel", Gender: "M", Type: "User", Status: "Active"}).Error; err != nil {
		t.Fatalf("seed sentinel: %v", err)
	}
	t.Cleanup(func() {
		db.Unscoped().Where("number = ?", sentinelNumber).Delete(&authModel.User{})
	})

	email := fmt.Sprintf("fub-%d@example.test", suffix)
	phone := fmt.Sprintf("077%09d", suffix%1000000000)

	status, body := fubRequest(t, server.Client(), server.URL, http.MethodPost, "/auth/register", "", map[string]string{
		"first_name": "Fub", "second_name": "Later", "last_name": "User", "gender": "Male",
		"email": email, "phone": phone,
	})
	if status != http.StatusOK {
		t.Fatalf("register = %d, want 200 (%s)", status, string(body))
	}

	var created struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &created)

	t.Cleanup(func() {
		if created.Data.Number != "" {
			db.Unscoped().Where("user_number = ?", created.Data.Number).Delete(&authModel.Contact{})
			db.Unscoped().Where("user_number = ?", created.Data.Number).Delete(&authModel.PasswordReset{})
			db.Unscoped().Where("number = ?", created.Data.Number).Delete(&authModel.User{})
		}
	})

	if created.Data.Number == "" {
		t.Fatalf("registration returned no user number: %s", string(body))
	}

	var stored authModel.User
	if err := db.Where("number = ?", created.Data.Number).First(&stored).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if stored.Type != "User" {
		t.Fatalf("later registrant type = %q, want User (must never be Admin)", stored.Type)
	}

	var roleCount int64
	db.Raw("SELECT COUNT(*) FROM user_roles WHERE user_number = ?", created.Data.Number).Scan(&roleCount)
	if roleCount != 0 {
		t.Fatalf("later registrant has %d user_roles, want 0", roleCount)
	}
}

// TestFirstUserBootstrapAbandonedRegistrationIsRecoverable covers TEST2/TEST3:
// an unactivated registration can be activated through the existing reset/
// set-password flow (the installation is not permanently locked), without
// granting Admin to any later registrant.
func TestFirstUserBootstrapAbandonedRegistrationIsRecoverable(t *testing.T) {
	db := requireOnboardingDB(t)
	t.Setenv("JWT_SECRET", "fub-test-secret")
	t.Setenv("MAIL_HOST", "127.0.0.1")
	t.Setenv("MAIL_PORT", "1")
	t.Setenv("MAIL_FROM", "fub@example.invalid")

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)

	suffix := time.Now().UnixNano() % 1000000000
	email := fmt.Sprintf("fub-rec-%d@example.test", suffix)
	phone := fmt.Sprintf("078%09d", suffix%1000000000)

	status, body := fubRequest(t, server.Client(), server.URL, http.MethodPost, "/auth/register", "", map[string]string{
		"first_name": "Fub", "second_name": "Recover", "last_name": "User", "gender": "Female",
		"email": email, "phone": phone,
	})
	if status != http.StatusOK {
		t.Fatalf("register = %d, want 200 (%s)", status, string(body))
	}
	var created struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &created)
	number := created.Data.Number
	if number == "" {
		t.Fatalf("no user number: %s", string(body))
	}
	t.Cleanup(func() {
		db.Unscoped().Where("user_number = ?", number).Delete(&authModel.Contact{})
		db.Unscoped().Where("user_number = ?", number).Delete(&authModel.PasswordReset{})
		db.Unscoped().Where("user_number = ?", number).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("number = ?", number).Delete(&authModel.User{})
	})

	// Abandoned: still Inactive.
	var before authModel.User
	db.Where("number = ?", number).First(&before)
	if before.Status != "Inactive" {
		t.Fatalf("pre-activation status = %q, want Inactive", before.Status)
	}

	// Recover: existing reset-password flow re-issues/returns the activation token.
	status, body = fubRequest(t, server.Client(), server.URL, http.MethodPost, "/auth/reset-password", "", map[string]string{"email": email})
	if status != http.StatusOK {
		t.Fatalf("reset-password = %d, want 200 (%s)", status, string(body))
	}
	var reset struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(body, &reset)
	if reset.Token == "" {
		t.Fatalf("reset-password returned no token: %s", string(body))
	}

	status, body = fubRequest(t, server.Client(), server.URL, http.MethodPost, "/auth/set-password", "", map[string]string{
		"token": reset.Token, "password": "Sup3rSecret!23", "confirm_password": "Sup3rSecret!23",
	})
	if status != http.StatusOK {
		t.Fatalf("set-password = %d, want 200 (%s)", status, string(body))
	}

	var after authModel.User
	db.Where("number = ?", number).First(&after)
	if after.Status != "Active" {
		t.Fatalf("post-activation status = %q, want Active", after.Status)
	}
}

// TestFirstUserBootstrapPermissionAndOwnership covers TEST6–TEST9: a role-less
// user is denied stakeholder management; a properly authorised administrator
// can create the organisation, select it and explicitly create the owner
// stakeholder; ownership does not change users.type; and holding ownership does
// not grant an administrator bypass.
func TestFirstUserBootstrapPermissionAndOwnership(t *testing.T) {
	db := requireOnboardingDB(t)
	t.Setenv("JWT_SECRET", "fub-test-secret")
	t.Setenv("MAIL_HOST", "127.0.0.1")
	t.Setenv("MAIL_PORT", "1")
	t.Setenv("MAIL_FROM", "fub@example.invalid")

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	suffix := time.Now().UnixNano() % 1000000000
	adminNumber := fmt.Sprintf("fub-admin-%d", suffix)
	plainNumber := fmt.Sprintf("fub-plain-%d", suffix)
	adminEmail := fmt.Sprintf("fub-admin-%d@example.test", suffix)
	plainEmail := fmt.Sprintf("fub-plain-%d@example.test", suffix)

	var createdEntities []string
	t.Cleanup(func() {
		db.Unscoped().Where("organisation_number IN ?", createdEntities).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("user_number IN ?", []string{adminNumber, plainNumber}).Delete(&stakeholderModel.Stakeholder{})
		for _, en := range createdEntities {
			db.Unscoped().Where("entity_number = ?", en).Delete(&eosModel.EntityUser{})
			db.Unscoped().Where("entity_number = ?", en).Delete(&eosModel.ActiveUserEntity{})
			db.Unscoped().Where("number = ?", en).Delete(&eosModel.Entity{})
		}
		db.Unscoped().Where("user_number IN ?", []string{adminNumber, plainNumber}).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", []string{adminNumber, plainNumber}).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("user_number IN ?", []string{adminNumber, plainNumber}).Delete(&authModel.Contact{})
		db.Unscoped().Where("number IN ?", []string{adminNumber, plainNumber}).Delete(&authModel.User{})
	})

	seed := func(number, email, typ string) {
		u := authModel.User{Number: number, FirstName: "Fub", SecondName: typ, Gender: "M", Type: typ, Status: "Active"}
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("seed user %s: %v", number, err)
		}
		c := authModel.Contact{Contact: email, ContactType: "Email", UserNumber: number}
		if err := db.Omit("User").Create(&c).Error; err != nil {
			t.Fatalf("seed contact %s: %v", email, err)
		}
	}
	seed(adminNumber, adminEmail, "Admin")
	seed(plainNumber, plainEmail, "User")

	token := func(number string) string {
		tok, err := authServices.JWTAuthService().GenerateToken(number)
		if err != nil {
			t.Fatalf("token: %v", err)
		}
		return tok
	}
	adminToken := token(adminNumber)
	plainToken := token(plainNumber)

	// TEST6: role-less user denied stakeholder management.
	if status, _ := fubRequest(t, client, server.URL, http.MethodPost, "/main/ems/stakeholders", plainToken, map[string]string{
		"user_number": plainNumber, "stakeholder_type": "owner",
	}); status != http.StatusForbidden {
		t.Fatalf("role-less stakeholder create = %d, want 403", status)
	}

	// TEST7: authorised administrator bootstraps the organisation.
	status, body := fubRequest(t, client, server.URL, http.MethodPost, "/main/eos/entity/add", adminToken, map[string]string{
		"name": fmt.Sprintf("Fub Org %d", suffix), "phone": "0700000009", "email": fmt.Sprintf("fub-org-%d@example.test", suffix), "address": "Kampala",
	})
	if status != http.StatusOK {
		t.Fatalf("entity/add = %d, want 200 (%s)", status, string(body))
	}

	status, body = fubRequest(t, client, server.URL, http.MethodGet, "/main/eos/my-entities", adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("my-entities = %d, want 200 (%s)", status, string(body))
	}
	var entities []eosModel.Entity
	_ = json.Unmarshal(body, &entities)
	if len(entities) == 0 {
		t.Fatalf("my-entities returned none: %s", string(body))
	}
	entityNumber := entities[len(entities)-1].Number
	createdEntities = append(createdEntities, entityNumber)

	status, body = fubRequest(t, client, server.URL, http.MethodPost, "/main/eos/entity/current?entity_number="+entityNumber, adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("entity/current = %d, want 200 (%s)", status, string(body))
	}

	status, body = fubRequest(t, client, server.URL, http.MethodPost, "/main/ems/stakeholders", adminToken, map[string]string{
		"user_number": adminNumber, "stakeholder_type": "owner",
	})
	if status != http.StatusCreated {
		t.Fatalf("owner stakeholder create = %d, want 201 (%s)", status, string(body))
	}

	// TEST8: ownership does not change users.type.
	var adminAfter authModel.User
	db.Where("number = ?", adminNumber).First(&adminAfter)
	if adminAfter.Type != "Admin" {
		t.Fatalf("admin type changed to %q after owner creation, want Admin", adminAfter.Type)
	}

	// TEST9: an owner stakeholder does NOT gain an administrator bypass — the
	// role-less plain user (also made an owner) still cannot manage stakeholders.
	status, body = fubRequest(t, client, server.URL, http.MethodPost, "/main/ems/stakeholders", adminToken, map[string]string{
		"user_number": plainNumber, "stakeholder_type": "owner",
	})
	if status != http.StatusCreated {
		t.Fatalf("second owner stakeholder create = %d, want 201 (%s)", status, string(body))
	}
	if status, _ := fubRequest(t, client, server.URL, http.MethodPost, "/main/ems/stakeholders", plainToken, map[string]string{
		"user_number": plainNumber, "stakeholder_type": "director",
	}); status != http.StatusForbidden {
		t.Fatalf("owner-role-less user stakeholder create = %d, want 403 (ownership must not grant admin)", status)
	}
}
