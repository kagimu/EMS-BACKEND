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
	stakeUserAdminA = "ems-stake-admin-a"
	stakeUserAdminB = "ems-stake-admin-b"
	stakeUserNoPerm = "ems-stake-noperm"
	stakeUserNoCtx  = "ems-stake-nocontext"
	stakePersonA    = "ems-stake-person-a"
	stakePersonB    = "ems-stake-person-b"
	stakePersonC    = "ems-stake-person-c"
	stakeOrgA       = "ems-stake-org-a"
	stakeOrgB       = "ems-stake-org-b"
	stakeBasePath   = "/main/ems/stakeholders"
	stakeEduURL     = "/main/campus/all"
)

type stakeholderFixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupStakeholderFixture(t *testing.T) *stakeholderFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-c1-test-secret")

	users := []string{
		stakeUserAdminA, stakeUserAdminB, stakeUserNoPerm, stakeUserNoCtx,
		stakePersonA, stakePersonB, stakePersonC,
	}
	organisations := []string{stakeOrgA, stakeOrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct{ number, userType, status string }{
		{stakeUserAdminA, "Admin", "Active"},
		{stakeUserAdminB, "Admin", "Active"},
		{stakeUserNoPerm, "User", "Active"},
		{stakeUserNoCtx, "Admin", "Active"},
		// Persons: no contact and no password, so they cannot log in.
		{stakePersonA, "User", "Inactive"},
		{stakePersonB, "User", "Inactive"},
		{stakePersonC, "User", "Inactive"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "Stake",
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
		{stakeOrgA, "Stakeholder Test Org A"},
		{stakeOrgB, "Stakeholder Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: stakeUserAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{stakeUserAdminA, stakeOrgA},
		{stakeUserAdminB, stakeOrgB},
		{stakeUserNoPerm, stakeOrgA},
		// stakeUserNoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &stakeholderFixture{
		db:          db,
		tokenAdminA: token(stakeUserAdminA),
		tokenAdminB: token(stakeUserAdminB),
		tokenNoPerm: token(stakeUserNoPerm),
		tokenNoCtx:  token(stakeUserNoCtx),
	}
}

type stakeholderEnvelope struct {
	Message string `json:"message"`
	Data    struct {
		Number          string   `json:"number"`
		UserNumber      string   `json:"user_number"`
		StakeholderType string   `json:"stakeholder_type"`
		Title           string   `json:"title"`
		Status          string   `json:"status"`
		SharePercentage *float64 `json:"share_percentage"`
		ValidTo         *string  `json:"valid_to"`
	} `json:"data"`
}

// TestStakeholders covers Foundation C1 end to end with one shared fixture.
func TestStakeholders(t *testing.T) {
	fx := setupStakeholderFixture(t)

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

	create := func(token, userNumber, stakeholderType string, extra map[string]any) (int, []byte) {
		payload := map[string]any{"user_number": userNumber, "stakeholder_type": stakeholderType}
		for key, value := range extra {
			payload[key] = value
		}
		return request(http.MethodPost, stakeBasePath, token, payload)
	}

	decode := func(t *testing.T, body []byte) stakeholderEnvelope {
		t.Helper()
		var envelope stakeholderEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding %s: %v", string(body), err)
		}
		return envelope
	}

	t.Run("unauthenticated_returns_401", func(t *testing.T) {
		if status, _ := request(http.MethodGet, stakeBasePath, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated_without_permission_returns_403", func(t *testing.T) {
		if status, _ := request(http.MethodGet, stakeBasePath, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("no_active_organisation_returns_409", func(t *testing.T) {
		if status, _ := request(http.MethodGet, stakeBasePath, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("status = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("owner_is_a_normal_generic_stakeholder", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, stakePersonA, "owner", nil)
		if status != http.StatusCreated {
			t.Fatalf("create owner = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		created := decode(t, body)
		if created.Data.StakeholderType != "owner" || created.Data.Status != "ACTIVE" {
			t.Fatalf("unexpected owner stakeholder: %s", string(body))
		}

		// The referenced person has no password, proving a stakeholder needs no login.
		var passwords int64
		if err := fx.db.Model(&authModel.UserPassword{}).Where("user_number = ?", stakePersonA).Count(&passwords).Error; err != nil {
			t.Fatalf("counting passwords: %v", err)
		}
		if passwords != 0 {
			t.Fatalf("person %s unexpectedly has a password", stakePersonA)
		}
	})

	t.Run("all_stakeholder_types_supported_for_one_person", func(t *testing.T) {
		types := []string{"shareholder", "board_member", "director", "executive", "staff", "consultant", "partner", "other"}
		for _, stakeholderType := range types {
			extra := map[string]any{}
			if stakeholderType == "shareholder" {
				extra["share_percentage"] = 25.5
				extra["shares"] = 1000
			}
			status, body := create(fx.tokenAdminA, stakePersonB, stakeholderType, extra)
			if status != http.StatusCreated {
				t.Fatalf("create %s = %d, want %d (%s)", stakeholderType, status, http.StatusCreated, string(body))
			}
			if decoded := decode(t, body); decoded.Data.StakeholderType != stakeholderType {
				t.Fatalf("created type = %q, want %q", decoded.Data.StakeholderType, stakeholderType)
			}
		}

		_, listBody := request(http.MethodGet, stakeBasePath+"?stakeholder_type=director", fx.tokenAdminA, nil)
		if !bytes.Contains(listBody, []byte(stakePersonB)) {
			t.Fatalf("filtered list must contain the person: %s", string(listBody))
		}
	})

	t.Run("create_retrieve_update_status_lifecycle", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, stakePersonA, "director", map[string]any{"title": "Operations Director"})
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		path := stakeBasePath + "/" + created.Data.Number

		if status, _ = request(http.MethodGet, path, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("get = %d, want %d", status, http.StatusOK)
		}

		status, body = request(http.MethodPut, path, fx.tokenAdminA, map[string]any{"title": "Group Operations Director", "valid_from": nil, "valid_to": nil})
		if status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}
		if updated := decode(t, body); updated.Data.Title != "Group Operations Director" {
			t.Fatalf("update did not apply: %s", string(body))
		}

		status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
		if updated := decode(t, body); updated.Data.Status != "INACTIVE" {
			t.Fatalf("status not applied: %s", string(body))
		}

		status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})
		if status != http.StatusOK {
			t.Fatalf("end = %d (%s)", status, string(body))
		}
		ended := decode(t, body)
		if ended.Data.Status != "ENDED" {
			t.Fatalf("expected ENDED, got %s", string(body))
		}
		if ended.Data.ValidTo == nil || *ended.Data.ValidTo == "" {
			t.Fatalf("ending a relationship must stamp valid_to: %s", string(body))
		}
	})

	t.Run("duplicate_active_and_history_handling", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, stakePersonC, "staff", nil)
		if status != http.StatusCreated {
			t.Fatalf("first create = %d (%s)", status, string(body))
		}
		first := decode(t, body)

		if status, _ = create(fx.tokenAdminA, stakePersonC, "staff", nil); status != http.StatusConflict {
			t.Fatalf("duplicate ACTIVE = %d, want %d", status, http.StatusConflict)
		}

		// Deactivate, then a new ACTIVE relationship of the same type is allowed.
		if status, _ = request(http.MethodPatch, stakeBasePath+"/"+first.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate = %d, want %d", status, http.StatusOK)
		}
		if status, body = create(fx.tokenAdminA, stakePersonC, "staff", nil); status != http.StatusCreated {
			t.Fatalf("re-create after deactivate = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
	})

	t.Run("invalid_inputs_rejected", func(t *testing.T) {
		if status, _ := create(fx.tokenAdminA, stakePersonC, "superhero", nil); status != http.StatusBadRequest {
			t.Fatalf("invalid type = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := create(fx.tokenAdminA, "does-not-exist", "staff", nil); status != http.StatusBadRequest {
			t.Fatalf("unknown user = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := create(fx.tokenAdminA, stakePersonC, "shareholder", map[string]any{"share_percentage": 150}); status != http.StatusBadRequest {
			t.Fatalf("invalid share percentage = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("no_duplicate_user_created", func(t *testing.T) {
		var before int64
		if err := fx.db.Model(&authModel.User{}).Count(&before).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}

		status, body := create(fx.tokenAdminA, stakePersonA, "consultant", nil)
		if status != http.StatusCreated {
			t.Fatalf("create consultant = %d (%s)", status, string(body))
		}

		var after int64
		if err := fx.db.Model(&authModel.User{}).Count(&after).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}
		if before != after {
			t.Fatalf("user count changed from %d to %d; stakeholders must not create users", before, after)
		}
	})

	t.Run("organisation_isolation_and_idor", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, stakePersonC, "partner", nil)
		if status != http.StatusCreated {
			t.Fatalf("org A create = %d (%s)", status, string(body))
		}
		orgAStakeholder := decode(t, body)
		path := stakeBasePath + "/" + orgAStakeholder.Data.Number

		_, listB := request(http.MethodGet, stakeBasePath, fx.tokenAdminB, nil)
		if bytes.Contains(listB, []byte(orgAStakeholder.Data.Number)) {
			t.Fatalf("org B list must not contain org A stakeholder: %s", string(listB))
		}

		if status, _ = request(http.MethodGet, path, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodPut, path, fx.tokenAdminB, map[string]any{"title": "Hijacked"}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation update = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodPatch, path+"/status", fx.tokenAdminB, map[string]any{"status": "ENDED"}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, stakeEduURL, fx.tokenAdminA, nil)

		status, body := create(fx.tokenAdminA, stakePersonC, "consultant", nil)
		if status != http.StatusCreated {
			t.Fatalf("scope-check create = %d (%s)", status, string(body))
		}
		unit := decode(t, body)
		request(http.MethodPatch, stakeBasePath+"/"+unit.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})

		statusAfter, bodyAfter := request(http.MethodGet, stakeEduURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS stakeholder changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS stakeholder changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}
