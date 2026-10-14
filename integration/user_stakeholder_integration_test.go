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
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	c2AdminA       = "ems-c2-admin-a"
	c2AdminB       = "ems-c2-admin-b"
	c2Switcher     = "ems-c2-switcher"
	c2NoPerm       = "ems-c2-noperm"
	c2NoCtx        = "ems-c2-nocontext"
	c2PersonLogin  = "ems-c2-person-login"
	c2PersonNoAuth = "ems-c2-person-noauth"
	c2PersonMulti  = "ems-c2-person-multi"
	c2PersonBare   = "ems-c2-person-bare"
	c2OrgA         = "ems-c2-org-a"
	c2OrgB         = "ems-c2-org-b"
	c2BasePath     = "/main/ems/stakeholders"
	c2EduURL       = "/main/campus/all"
	// Placeholder bcrypt-shaped value. It is never used to authenticate in this
	// test; only the existence of the password row matters for can_authenticate.
	c2PasswordHash = "$2a$10$abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcd"
)

type c2Fixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenSwitch string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupC2Fixture(t *testing.T) *c2Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-c2-test-secret")

	users := []string{
		c2AdminA, c2AdminB, c2Switcher, c2NoPerm, c2NoCtx,
		c2PersonLogin, c2PersonNoAuth, c2PersonMulti, c2PersonBare,
	}
	organisations := []string{c2OrgA, c2OrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.Contact{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct{ number, userType, status string }{
		{c2AdminA, "Admin", "Active"},
		{c2AdminB, "Admin", "Active"},
		{c2Switcher, "Admin", "Active"},
		{c2NoPerm, "User", "Active"},
		{c2NoCtx, "Admin", "Active"},
		{c2PersonLogin, "User", "Active"},
		{c2PersonNoAuth, "User", "Inactive"},
		{c2PersonMulti, "User", "Active"},
		{c2PersonBare, "User", "Active"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "C2",
			SecondName: spec.userType,
			Gender:     "M",
			Type:       spec.userType,
			Status:     spec.status,
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{c2OrgA, "C2 Test Org A"},
		{c2OrgB, "C2 Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: c2AdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	// Active organisation context (membership-backed in real usage).
	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c2AdminA, c2OrgA},
		{c2AdminB, c2OrgB},
		{c2Switcher, c2OrgA},
		{c2NoPerm, c2OrgA},
		// c2NoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	// Membership: only the switcher belongs to two organisations.
	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c2Switcher, c2OrgA},
		{c2Switcher, c2OrgB},
	} {
		membership := eosModel.EntityUser{
			EntityNumber:  spec.entityNumber,
			UserNumber:    spec.userNumber,
			AddedByNumber: c2AdminA,
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	// Login credentials: only the two "login" persons get an Email contact and a password.
	for _, userNumber := range []string{c2PersonLogin, c2PersonMulti} {
		contact := authModel.Contact{Contact: userNumber + "@example.com", ContactType: "Email", UserNumber: userNumber}
		if err := db.Omit("User").Create(&contact).Error; err != nil {
			t.Fatalf("seeding contact for %s: %v", userNumber, err)
		}
		password := authModel.UserPassword{UserNumber: userNumber, UserPassword: c2PasswordHash}
		if err := db.Omit("User").Create(&password).Error; err != nil {
			t.Fatalf("seeding password for %s: %v", userNumber, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &c2Fixture{
		db:          db,
		tokenAdminA: token(c2AdminA),
		tokenAdminB: token(c2AdminB),
		tokenSwitch: token(c2Switcher),
		tokenNoPerm: token(c2NoPerm),
		tokenNoCtx:  token(c2NoCtx),
	}
}

type c2Envelope struct {
	Message string `json:"message"`
	Data    struct {
		Number          string  `json:"number"`
		UserNumber      string  `json:"user_number"`
		StakeholderType string  `json:"stakeholder_type"`
		Status          string  `json:"status"`
		CanAuthenticate bool    `json:"can_authenticate"`
		ValidTo         *string `json:"valid_to"`
	} `json:"data"`
}

// TestUserStakeholderIntegration covers Foundation C2.
func TestUserStakeholderIntegration(t *testing.T) {
	fx := setupC2Fixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	request := func(method, path, token string, payload any) (int, []byte) {
		t.Helper()

		var reader io.Reader
		if payload != nil {
			raw, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("encoding payload: %v", err)
			}
			reader = bytes.NewReader(raw)
		}

		req, err := http.NewRequest(method, server.URL+path, reader)
		if err != nil {
			t.Fatalf("building request %s %s: %v", method, path, err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	create := func(token, userNumber, stakeholderType string) (int, []byte) {
		return request(http.MethodPost, c2BasePath, token, map[string]any{
			"user_number": userNumber, "stakeholder_type": stakeholderType,
		})
	}

	decode := func(t *testing.T, body []byte) c2Envelope {
		t.Helper()
		var envelope c2Envelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding %s: %v", string(body), err)
		}
		return envelope
	}

	countUsers := func(number string) int64 {
		t.Helper()
		var count int64
		if err := fx.db.Model(&authModel.User{}).Where("number = ?", number).Count(&count).Error; err != nil {
			t.Fatalf("counting users for %s: %v", number, err)
		}
		return count
	}

	t.Run("authentication_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c2BasePath, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("permission_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c2BasePath, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission status = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, c2BasePath, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("no-context status = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("user_with_stakeholder_and_can_authenticate", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c2PersonLogin, "director")
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		if !created.Data.CanAuthenticate {
			t.Fatalf("a user with an Email contact, a password and Active status must report can_authenticate=true: %s", string(body))
		}
	})

	t.Run("stakeholder_without_login_credentials", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c2PersonNoAuth, "shareholder")
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		if created.Data.Status != "ACTIVE" {
			t.Fatalf("relationship status = %q, want ACTIVE", created.Data.Status)
		}
		if created.Data.CanAuthenticate {
			t.Fatalf("a user with no contact/password must report can_authenticate=false: %s", string(body))
		}

		var passwords int64
		if err := fx.db.Model(&authModel.UserPassword{}).Where("user_number = ?", c2PersonNoAuth).Count(&passwords).Error; err != nil {
			t.Fatalf("counting passwords: %v", err)
		}
		if passwords != 0 {
			t.Fatalf("stakeholder operations must not create a password for %s", c2PersonNoAuth)
		}
	})

	t.Run("user_without_stakeholder", func(t *testing.T) {
		var count int64
		if err := fx.db.Model(&stakeholderModel.Stakeholder{}).
			Where("organisation_number = ? AND user_number = ?", c2OrgA, c2PersonBare).Count(&count).Error; err != nil {
			t.Fatalf("counting stakeholders: %v", err)
		}
		if count != 0 {
			t.Fatalf("user %s must have no stakeholder relationship", c2PersonBare)
		}
	})

	t.Run("user_with_multiple_stakeholder_types", func(t *testing.T) {
		for _, stakeholderType := range []string{"director", "shareholder", "board_member"} {
			if status, body := create(fx.tokenAdminA, c2PersonMulti, stakeholderType); status != http.StatusCreated {
				t.Fatalf("create %s = %d (%s)", stakeholderType, status, string(body))
			}
		}
		if users := countUsers(c2PersonMulti); users != 1 {
			t.Fatalf("person %s must remain exactly one user row, found %d", c2PersonMulti, users)
		}
	})

	t.Run("owner_with_active_user", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c2PersonLogin, "owner")
		if status != http.StatusCreated {
			t.Fatalf("create owner = %d (%s)", status, string(body))
		}
		if decoded := decode(t, body); decoded.Data.StakeholderType != "owner" {
			t.Fatalf("expected owner relationship: %s", string(body))
		}
	})

	t.Run("inactive_authentication_status", func(t *testing.T) {
		var user authModel.User
		if err := fx.db.Where("number = ?", c2PersonNoAuth).First(&user).Error; err != nil {
			t.Fatalf("loading user: %v", err)
		}
		if user.Status != "Inactive" {
			t.Fatalf("authentication status = %q, want Inactive", user.Status)
		}

		// The relationship exists and is ACTIVE, but the person cannot authenticate.
		var active int64
		if err := fx.db.Model(&stakeholderModel.Stakeholder{}).
			Where("organisation_number = ? AND user_number = ? AND status = ?", c2OrgA, c2PersonNoAuth, "ACTIVE").
			Count(&active).Error; err != nil {
			t.Fatalf("counting active relationships: %v", err)
		}
		if active != 1 {
			t.Fatalf("expected exactly one ACTIVE relationship for %s", c2PersonNoAuth)
		}
	})

	t.Run("inactive_and_ended_stakeholder_status", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c2PersonMulti, "staff")
		if status != http.StatusCreated {
			t.Fatalf("create staff = %d (%s)", status, string(body))
		}
		staff := decode(t, body)

		status, body = request(http.MethodPatch, c2BasePath+"/"+staff.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
		if decoded := decode(t, body); decoded.Data.Status != "INACTIVE" {
			t.Fatalf("expected INACTIVE: %s", string(body))
		}

		status, body = create(fx.tokenAdminA, c2PersonMulti, "consultant")
		if status != http.StatusCreated {
			t.Fatalf("create consultant = %d (%s)", status, string(body))
		}
		consultant := decode(t, body)

		status, body = request(http.MethodPatch, c2BasePath+"/"+consultant.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})
		if status != http.StatusOK {
			t.Fatalf("end = %d (%s)", status, string(body))
		}
		ended := decode(t, body)
		if ended.Data.Status != "ENDED" || ended.Data.ValidTo == nil {
			t.Fatalf("ending must set ENDED and stamp valid_to: %s", string(body))
		}

		// A new ACTIVE relationship of the same type is allowed after ENDED.
		if status, body = create(fx.tokenAdminA, c2PersonMulti, "consultant"); status != http.StatusCreated {
			t.Fatalf("re-create after ENDED = %d (%s)", status, string(body))
		}
	})

	t.Run("organisation_membership_mismatch", func(t *testing.T) {
		var memberships int64
		if err := fx.db.Model(&eosModel.EntityUser{}).Where("user_number = ?", c2PersonLogin).Count(&memberships).Error; err != nil {
			t.Fatalf("counting memberships: %v", err)
		}
		if memberships != 0 {
			t.Fatalf("%s must have no EntityUser membership for this test", c2PersonLogin)
		}

		// A stakeholder relationship is permitted without membership (C2 V4).
		if status, body := create(fx.tokenAdminA, c2PersonLogin, "partner"); status != http.StatusCreated {
			t.Fatalf("membership-less stakeholder = %d (%s)", status, string(body))
		}
	})

	t.Run("cross_organisation_identity_reference", func(t *testing.T) {
		statusA, bodyA := create(fx.tokenAdminA, c2PersonLogin, "consultant")
		if statusA != http.StatusCreated {
			t.Fatalf("org A create = %d (%s)", statusA, string(bodyA))
		}
		orgARow := decode(t, bodyA)

		statusB, bodyB := create(fx.tokenAdminB, c2PersonLogin, "consultant")
		if statusB != http.StatusCreated {
			t.Fatalf("org B create = %d (%s)", statusB, string(bodyB))
		}
		orgBRow := decode(t, bodyB)

		_, listA := request(http.MethodGet, c2BasePath, fx.tokenAdminA, nil)
		if !bytes.Contains(listA, []byte(orgARow.Data.Number)) {
			t.Fatalf("org A list must contain its own row")
		}
		if bytes.Contains(listA, []byte(orgBRow.Data.Number)) {
			t.Fatalf("org A list must not contain org B's row: %s", string(listA))
		}
	})

	t.Run("organisation_switching", func(t *testing.T) {
		statusA, bodyA := create(fx.tokenAdminA, c2PersonLogin, "board_member")
		if statusA != http.StatusCreated {
			t.Fatalf("org A board_member = %d (%s)", statusA, string(bodyA))
		}
		orgARow := decode(t, bodyA)

		statusB, bodyB := create(fx.tokenAdminB, c2PersonLogin, "board_member")
		if statusB != http.StatusCreated {
			t.Fatalf("org B partner = %d (%s)", statusB, string(bodyB))
		}
		orgBRow := decode(t, bodyB)

		_, listBefore := request(http.MethodGet, c2BasePath, fx.tokenSwitch, nil)
		if !bytes.Contains(listBefore, []byte(orgARow.Data.Number)) || bytes.Contains(listBefore, []byte(orgBRow.Data.Number)) {
			t.Fatalf("switcher in org A must see only org A rows: %s", string(listBefore))
		}

		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", c2Switcher).Update("entity_number", c2OrgB).Error; err != nil {
			t.Fatalf("switching active entity: %v", err)
		}

		_, listAfter := request(http.MethodGet, c2BasePath, fx.tokenSwitch, nil)
		if !bytes.Contains(listAfter, []byte(orgBRow.Data.Number)) || bytes.Contains(listAfter, []byte(orgARow.Data.Number)) {
			t.Fatalf("switcher in org B must see only org B rows: %s", string(listAfter))
		}
	})

	t.Run("no_duplicate_identity", func(t *testing.T) {
		var before int64
		if err := fx.db.Model(&authModel.User{}).Count(&before).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}

		if status, body := create(fx.tokenAdminA, c2PersonMulti, "executive"); status != http.StatusCreated {
			t.Fatalf("create executive = %d (%s)", status, string(body))
		}

		var after int64
		if err := fx.db.Model(&authModel.User{}).Count(&after).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}
		if before != after {
			t.Fatalf("user count changed from %d to %d; stakeholder operations must not create identities", before, after)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, c2EduURL, fx.tokenAdminA, nil)

		status, body := create(fx.tokenAdminA, c2PersonBare, "director")
		if status != http.StatusCreated {
			t.Fatalf("scope-check create = %d (%s)", status, string(body))
		}
		row := decode(t, body)
		request(http.MethodPatch, c2BasePath+"/"+row.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})

		statusAfter, bodyAfter := request(http.MethodGet, c2EduURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS stakeholder changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS stakeholder changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}
