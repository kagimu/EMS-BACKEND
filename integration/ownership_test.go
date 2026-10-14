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
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	c4AdminA   = "ems-c4-admin-a"
	c4AdminB   = "ems-c4-admin-b"
	c4Switcher = "ems-c4-switcher"
	c4NoPerm   = "ems-c4-noperm"
	c4NoCtx    = "ems-c4-nocontext"
	c4Owner1   = "ems-c4-owner-1" // has a login AND Admin type, so it can act on itself
	c4Owner2   = "ems-c4-owner-2" // no login
	c4Owner3   = "ems-c4-owner-3" // inactive user
	c4Owner4   = "ems-c4-owner-4"

	c4OrgA = "ems-c4-org-a"
	c4OrgB = "ems-c4-org-b"

	c4UnitA    = "ems-c4-unit-a"
	c4BasePath = "/main/ems/stakeholders"
	c4EduURL   = "/main/campus/all"
)

type c4Fixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenSwitch string
	tokenNoPerm string
	tokenNoCtx  string
	tokenOwner1 string
}

func setupC4Fixture(t *testing.T) *c4Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-c4-test-secret")

	users := []string{c4AdminA, c4AdminB, c4Switcher, c4NoPerm, c4NoCtx, c4Owner1, c4Owner2, c4Owner3, c4Owner4}
	organisations := []string{c4OrgA, c4OrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
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
		{c4AdminA, "Admin", "Active"},
		{c4AdminB, "Admin", "Active"},
		{c4Switcher, "Admin", "Active"},
		{c4NoPerm, "User", "Active"},
		{c4NoCtx, "Admin", "Active"},
		{c4Owner1, "Admin", "Active"},
		{c4Owner2, "User", "Active"},
		{c4Owner3, "User", "Inactive"},
		{c4Owner4, "User", "Active"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "C4",
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
		{c4OrgA, "C4 Test Org A"},
		{c4OrgB, "C4 Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: c4AdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c4AdminA, c4OrgA},
		{c4AdminB, c4OrgB},
		{c4Switcher, c4OrgA},
		{c4NoPerm, c4OrgA},
		// c4NoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c4Switcher, c4OrgA},
		{c4Switcher, c4OrgB},
	} {
		membership := eosModel.EntityUser{EntityNumber: spec.entityNumber, UserNumber: spec.userNumber, AddedByNumber: c4AdminA}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	// Owner 1 has a usable login (Email contact + password + Active status).
	contact := authModel.Contact{Contact: c4Owner1 + "@example.com", ContactType: "Email", UserNumber: c4Owner1}
	if err := db.Omit("User").Create(&contact).Error; err != nil {
		t.Fatalf("seeding contact: %v", err)
	}
	password := authModel.UserPassword{UserNumber: c4Owner1, UserPassword: "$2a$10$abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcd"}
	if err := db.Omit("User").Create(&password).Error; err != nil {
		t.Fatalf("seeding password: %v", err)
	}

	unit := organisationModel.OrganisationUnit{
		Number: c4UnitA, OrganisationNumber: c4OrgA, UnitType: "department",
		Name: "Unit A", Status: "ACTIVE", AddedByNumber: c4AdminA,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("seeding unit: %v", err)
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &c4Fixture{
		db:          db,
		tokenAdminA: token(c4AdminA),
		tokenAdminB: token(c4AdminB),
		tokenSwitch: token(c4Switcher),
		tokenNoPerm: token(c4NoPerm),
		tokenNoCtx:  token(c4NoCtx),
		tokenOwner1: token(c4Owner1),
	}
}

type c4Envelope struct {
	Message string `json:"message"`
	Data    struct {
		Number          string   `json:"number"`
		UserNumber      string   `json:"user_number"`
		StakeholderType string   `json:"stakeholder_type"`
		Status          string   `json:"status"`
		CanAuthenticate bool     `json:"can_authenticate"`
		SharePercentage *float64 `json:"share_percentage"`
	} `json:"data"`
}

// TestOwnership covers Foundation C4.
func TestOwnership(t *testing.T) {
	fx := setupC4Fixture(t)

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
			t.Fatalf("building request: %v", err)
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

	createOwner := func(token, userNumber string) (int, []byte) {
		return request(http.MethodPost, c4BasePath, token, map[string]any{
			"user_number": userNumber, "stakeholder_type": "owner",
		})
	}

	decode := func(t *testing.T, body []byte) c4Envelope {
		t.Helper()
		var envelope c4Envelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding %s: %v", string(body), err)
		}
		return envelope
	}

	setStatus := func(token, number, status string) (int, []byte) {
		return request(http.MethodPatch, c4BasePath+"/"+number+"/status", token, map[string]any{"status": status})
	}

	var owner1Number, owner2Number, owner3Number string

	t.Run("authentication_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("permission_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("no-context = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("create_and_retrieve_owner", func(t *testing.T) {
		status, body := createOwner(fx.tokenAdminA, c4Owner1)
		if status != http.StatusCreated {
			t.Fatalf("create owner = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		if created.Data.StakeholderType != "owner" || created.Data.Status != "ACTIVE" {
			t.Fatalf("unexpected owner: %s", string(body))
		}
		if !created.Data.CanAuthenticate {
			t.Fatalf("owner with a login must report can_authenticate=true: %s", string(body))
		}
		owner1Number = created.Data.Number

		if status, _ = request(http.MethodGet, c4BasePath+"/"+owner1Number, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("retrieve owner = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("owner_without_entity_membership", func(t *testing.T) {
		var memberships int64
		if err := fx.db.Model(&eosModel.EntityUser{}).Where("user_number = ?", c4Owner1).Count(&memberships).Error; err != nil {
			t.Fatalf("counting memberships: %v", err)
		}
		if memberships != 0 {
			t.Fatalf("%s must have no EntityUser membership for this test", c4Owner1)
		}
	})

	t.Run("multiple_active_owners_and_owner_without_login", func(t *testing.T) {
		status, body := createOwner(fx.tokenAdminA, c4Owner2)
		if status != http.StatusCreated {
			t.Fatalf("second owner = %d (%s)", status, string(body))
		}
		second := decode(t, body)
		owner2Number = second.Data.Number
		if second.Data.CanAuthenticate {
			t.Fatalf("owner without a login must report can_authenticate=false: %s", string(body))
		}

		_, list := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", fx.tokenAdminA, nil)
		if !bytes.Contains(list, []byte(owner1Number)) || !bytes.Contains(list, []byte(owner2Number)) {
			t.Fatalf("both active owners must be listed: %s", string(list))
		}
	})

	t.Run("inactive_user_can_be_owner", func(t *testing.T) {
		status, body := createOwner(fx.tokenAdminA, c4Owner3)
		if status != http.StatusCreated {
			t.Fatalf("inactive-user owner = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		owner3Number = created.Data.Number
		if created.Data.CanAuthenticate {
			t.Fatalf("an inactive user must not report can_authenticate=true")
		}
	})

	t.Run("duplicate_active_owner_rejected", func(t *testing.T) {
		if status, _ := createOwner(fx.tokenAdminA, c4Owner1); status != http.StatusConflict {
			t.Fatalf("duplicate active owner = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("inactive_owner_does_not_count", func(t *testing.T) {
		if status, body := setStatus(fx.tokenAdminA, owner3Number, "INACTIVE"); status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
	})

	t.Run("end_non_final_owner", func(t *testing.T) {
		// Two ACTIVE owners exist (owner1, owner2), so ending owner2 is allowed.
		status, body := setStatus(fx.tokenAdminA, owner2Number, "ENDED")
		if status != http.StatusOK {
			t.Fatalf("end non-final owner = %d (%s)", status, string(body))
		}
		if ended := decode(t, body); ended.Data.Status != "ENDED" {
			t.Fatalf("expected ENDED: %s", string(body))
		}
	})

	t.Run("ownership_history", func(t *testing.T) {
		status, body := request(http.MethodGet, c4BasePath+"/"+owner2Number, fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("history retrieve = %d, want %d", status, http.StatusOK)
		}
		if historical := decode(t, body); historical.Data.Status != "ENDED" {
			t.Fatalf("historical owner must remain retrievable as ENDED: %s", string(body))
		}

		_, list := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner&status=ENDED", fx.tokenAdminA, nil)
		if !bytes.Contains(list, []byte(owner2Number)) {
			t.Fatalf("ended ownership must remain in history: %s", string(list))
		}
	})

	t.Run("ownership_percentage_update", func(t *testing.T) {
		status, body := request(http.MethodPut, c4BasePath+"/"+owner1Number, fx.tokenAdminA, map[string]any{
			"title": "Founder", "share_percentage": 30,
		})
		if status != http.StatusOK {
			t.Fatalf("percentage update = %d (%s)", status, string(body))
		}
		updated := decode(t, body)
		if updated.Data.SharePercentage == nil || *updated.Data.SharePercentage != 30 {
			t.Fatalf("share percentage not updated: %s", string(body))
		}
	})

	t.Run("final_owner_protection", func(t *testing.T) {
		// At this point owner1 is the ONLY ACTIVE owner (owner2 ENDED, owner3 INACTIVE).

		if status, body := setStatus(fx.tokenAdminA, owner1Number, "INACTIVE"); status != http.StatusConflict {
			t.Fatalf("deactivate final owner = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
		if status, body := setStatus(fx.tokenAdminA, owner1Number, "ENDED"); status != http.StatusConflict {
			t.Fatalf("end final owner = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}

		// Actor-independent: the final owner cannot end themselves either.
		if status, body := setStatus(fx.tokenOwner1, owner1Number, "ENDED"); status != http.StatusConflict {
			t.Fatalf("self-removal of final owner = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}

		// Still an ACTIVE owner.
		status, body := request(http.MethodGet, c4BasePath+"/"+owner1Number, fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("final owner retrieve = %d", status)
		}
		if final := decode(t, body); final.Data.Status != "ACTIVE" {
			t.Fatalf("final owner must remain ACTIVE: %s", string(body))
		}
	})

	t.Run("ownership_transfer_sequence", func(t *testing.T) {
		// Create the new owner first, then the previous owner can be ended.
		status, body := createOwner(fx.tokenAdminA, c4Owner4)
		if status != http.StatusCreated {
			t.Fatalf("new owner = %d (%s)", status, string(body))
		}
		if status, body = setStatus(fx.tokenAdminA, owner1Number, "ENDED"); status != http.StatusOK {
			t.Fatalf("end previous owner once another active owner exists = %d (%s)", status, string(body))
		}
	})

	t.Run("organisation_switching", func(t *testing.T) {
		status, body := createOwner(fx.tokenAdminB, c4Owner2)
		if status != http.StatusCreated {
			t.Fatalf("org B owner = %d (%s)", status, string(body))
		}
		orgBOwner := decode(t, body)

		_, listA := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", fx.tokenSwitch, nil)
		if bytes.Contains(listA, []byte(orgBOwner.Data.Number)) {
			t.Fatalf("switcher in org A must not see org B owners: %s", string(listA))
		}

		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", c4Switcher).Update("entity_number", c4OrgB).Error; err != nil {
			t.Fatalf("switching active entity: %v", err)
		}

		_, listB := request(http.MethodGet, c4BasePath+"?stakeholder_type=owner", fx.tokenSwitch, nil)
		if !bytes.Contains(listB, []byte(orgBOwner.Data.Number)) {
			t.Fatalf("switcher in org B must see org B owners: %s", string(listB))
		}
	})

	t.Run("cross_organisation_idor", func(t *testing.T) {
		path := c4BasePath + "/" + owner1Number

		if status, _ := request(http.MethodGet, path, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ := request(http.MethodPut, path, fx.tokenAdminB, map[string]any{"title": "Hijacked"}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation update = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ := request(http.MethodPatch, path+"/status", fx.tokenAdminB, map[string]any{"status": "ENDED"}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, c4EduURL, fx.tokenAdminA, nil)

		status, body := createOwner(fx.tokenAdminA, c4Owner3)
		if status != http.StatusCreated {
			t.Fatalf("scope-check owner = %d (%s)", status, string(body))
		}
		row := decode(t, body)
		setStatus(fx.tokenAdminA, row.Data.Number, "INACTIVE")

		statusAfter, bodyAfter := request(http.MethodGet, c4EduURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS ownership changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS ownership changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}
