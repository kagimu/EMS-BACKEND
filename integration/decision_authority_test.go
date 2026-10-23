package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	dAdminA    = "ems-d-admin-a"
	dAdminB    = "ems-d-admin-b"
	dSwitcher  = "ems-d-switcher"
	dNoPerm    = "ems-d-noperm"
	dNoCtx     = "ems-d-nocontext"
	dOwnerUser = "ems-d-owner-user"

	dOrgA = "ems-d-org-a"
	dOrgB = "ems-d-org-b"

	dRoleA     = "ems-d-role-a"
	dRoleB     = "ems-d-role-b"
	dUnitA     = "ems-d-unit-a"
	dUnitCh    = "ems-d-unit-child"
	dUnitOt    = "ems-d-unit-other"
	dUnitIn    = "ems-d-unit-inactive"
	dUnitB     = "ems-d-unit-b"
	dUnitDeact = "ems-d-unit-deact"

	dAuthBase   = "/main/ems/decision-authorities"
	dResolveURL = "/main/ems/authority-resolution"
	dUnitBase   = "/main/ems/organisation-units"
	dEduURL     = "/main/campus/all"
)

type dFixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenSwitch string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupDFixture(t *testing.T) *dFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-d-test-secret")

	users := []string{dAdminA, dAdminB, dSwitcher, dNoPerm, dNoCtx, dOwnerUser}
	organisations := []string{dOrgA, dOrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&authorityModel.DecisionAuthority{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct{ number, userType, status string }{
		{dAdminA, "Admin", "Active"},
		{dAdminB, "Admin", "Active"},
		{dSwitcher, "Admin", "Active"},
		{dNoPerm, "User", "Active"},
		{dNoCtx, "Admin", "Active"},
		{dOwnerUser, "User", "Active"},
	} {
		user := authModel.User{
			Number: spec.number, FirstName: "D", SecondName: spec.userType,
			Gender: "M", Type: spec.userType, Status: spec.status,
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{dOrgA, "D Test Org A"},
		{dOrgB, "D Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: dAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{dAdminA, dOrgA}, {dAdminB, dOrgB}, {dSwitcher, dOrgA}, {dNoPerm, dOrgA},
		// dNoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{dSwitcher, dOrgA}, {dSwitcher, dOrgB},
	} {
		membership := eosModel.EntityUser{EntityNumber: spec.entityNumber, UserNumber: spec.userNumber, AddedByNumber: dAdminA}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	for _, spec := range []struct{ number, organisation, name string }{
		{dRoleA, dOrgA, "Finance Manager"},
		{dRoleB, dOrgB, "Org B Role"},
	} {
		role := permissionsModel.Role{
			Number: spec.number, Name: spec.name, EntityNumber: spec.organisation,
			Status: "ACTIVE", UserNumber: dAdminA,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", spec.number, err)
		}
	}

	type unitSpec struct{ number, organisation, parent, status string }
	for _, spec := range []unitSpec{
		{dUnitA, dOrgA, "", "ACTIVE"},
		{dUnitCh, dOrgA, dUnitA, "ACTIVE"},
		{dUnitOt, dOrgA, "", "ACTIVE"},
		{dUnitIn, dOrgA, "", "INACTIVE"},
		{dUnitDeact, dOrgA, "", "ACTIVE"},
		{dUnitB, dOrgB, "", "ACTIVE"},
	} {
		unit := organisationModel.OrganisationUnit{
			Number: spec.number, OrganisationNumber: spec.organisation, ParentNumber: spec.parent,
			UnitType: "department", Name: spec.number, Status: spec.status, AddedByNumber: dAdminA,
		}
		if err := db.Create(&unit).Error; err != nil {
			t.Fatalf("seeding unit %s: %v", spec.number, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &dFixture{
		db:          db,
		tokenAdminA: token(dAdminA),
		tokenAdminB: token(dAdminB),
		tokenSwitch: token(dSwitcher),
		tokenNoPerm: token(dNoPerm),
		tokenNoCtx:  token(dNoCtx),
	}
}

type dAuthorityEnvelope struct {
	Message string `json:"message"`
	Data    struct {
		Number          string   `json:"number"`
		DecisionType    string   `json:"decision_type"`
		ApprovalLevel   int      `json:"approval_level"`
		MinAmount       *float64 `json:"min_amount"`
		MaxAmount       *float64 `json:"max_amount"`
		RoleNumber      string   `json:"role_number"`
		StakeholderType string   `json:"stakeholder_type"`
		Status          string   `json:"status"`
	} `json:"data"`
}

type dResolutionEnvelope struct {
	Message string `json:"message"`
	Data    struct {
		Configured               bool    `json:"configured"`
		DecisionType             string  `json:"decision_type"`
		Amount                   float64 `json:"amount"`
		OrganisationUnitNumber   string  `json:"organisation_unit_number"`
		ExceedsDeclaredAuthority bool    `json:"exceeds_declared_authority"`
		RequiredPath             []struct {
			ApprovalLevel int `json:"approval_level"`
			Anchors       []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"anchors"`
		} `json:"required_path"`
	} `json:"data"`
}

// TestDecisionAuthority covers Foundation D.
func TestDecisionAuthority(t *testing.T) {
	fx := setupDFixture(t)

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

	createAuthority := func(token string, payload map[string]any) (int, []byte) {
		return request(http.MethodPost, dAuthBase, token, payload)
	}

	decodeAuthority := func(t *testing.T, body []byte) dAuthorityEnvelope {
		t.Helper()
		var envelope dAuthorityEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding %s: %v", string(body), err)
		}
		return envelope
	}

	resolve := func(token, decisionType string, amount float64, unit string) (int, dResolutionEnvelope, []byte) {
		t.Helper()
		path := dResolveURL + "?decision_type=" + decisionType + "&amount=" + formatAmount(amount)
		if unit != "" {
			path += "&organisation_unit=" + unit
		}
		status, body := request(http.MethodGet, path, token, nil)
		var envelope dResolutionEnvelope
		if status == http.StatusOK {
			if err := json.Unmarshal(body, &envelope); err != nil {
				t.Fatalf("decoding resolution %s: %v", string(body), err)
			}
		}
		return status, envelope, body
	}

	levels := func(envelope dResolutionEnvelope) []int {
		result := make([]int, 0, len(envelope.Data.RequiredPath))
		for _, level := range envelope.Data.RequiredPath {
			result = append(result, level.ApprovalLevel)
		}
		return result
	}

	rolePayload := func(decisionType string, level int, max *float64) map[string]any {
		payload := map[string]any{
			"decision_type": decisionType, "approval_level": level, "role_number": dRoleA,
		}
		if max != nil {
			payload["max_amount"] = *max
		}
		return payload
	}

	f := func(v float64) *float64 { return &v }

	t.Run("authentication_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, dAuthBase, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := request(http.MethodGet, dResolveURL+"?decision_type=x&amount=1", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("resolver status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("permission_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, dAuthBase, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, dAuthBase, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("no-context = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("create_retrieve_update_lifecycle", func(t *testing.T) {
		status, body := createAuthority(fx.tokenAdminA, rolePayload("purchase_request", 1, f(500000)))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decodeAuthority(t, body)
		if created.Data.Status != "ACTIVE" || created.Data.ApprovalLevel != 1 {
			t.Fatalf("unexpected authority: %s", string(body))
		}
		path := dAuthBase + "/" + created.Data.Number

		if status, _ = request(http.MethodGet, path, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("retrieve = %d, want %d", status, http.StatusOK)
		}

		status, body = request(http.MethodPut, path, fx.tokenAdminA, rolePayload("purchase_request", 1, f(600000)))
		if status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}
		if updated := decodeAuthority(t, body); updated.Data.MaxAmount == nil || *updated.Data.MaxAmount != 600000 {
			t.Fatalf("update not applied: %s", string(body))
		}

		for _, step := range []string{"INACTIVE", "ACTIVE", "ENDED"} {
			if status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]any{"status": step}); status != http.StatusOK {
				t.Fatalf("status %s = %d (%s)", step, status, string(body))
			}
		}
	})

	t.Run("anchor_and_validation_rules", func(t *testing.T) {
		// Neither anchor.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 1,
		}); status != http.StatusBadRequest {
			t.Fatalf("no anchor = %d, want %d", status, http.StatusBadRequest)
		}
		// Both anchors.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 1, "role_number": dRoleA, "stakeholder_type": "owner",
		}); status != http.StatusBadRequest {
			t.Fatalf("both anchors = %d, want %d", status, http.StatusBadRequest)
		}
		// Overlong stakeholder type.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 1, "stakeholder_type": "stakeholder-type-over-thirty-chars",
		}); status != http.StatusBadRequest {
			t.Fatalf("overlong stakeholder type = %d, want %d", status, http.StatusBadRequest)
		}
		// Unknown role.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 1, "role_number": "does-not-exist",
		}); status != http.StatusBadRequest {
			t.Fatalf("unknown role = %d, want %d", status, http.StatusBadRequest)
		}
		// Level < 1.
		if status, _ := createAuthority(fx.tokenAdminA, rolePayload("purchase_request", 0, f(1000))); status != http.StatusBadRequest {
			t.Fatalf("level 0 = %d, want %d", status, http.StatusBadRequest)
		}
		// Negative amount.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 2, "role_number": dRoleA, "max_amount": -5,
		}); status != http.StatusBadRequest {
			t.Fatalf("negative amount = %d, want %d", status, http.StatusBadRequest)
		}
		// min > max.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 2, "role_number": dRoleA,
			"min_amount": 5000, "max_amount": 1000,
		}); status != http.StatusBadRequest {
			t.Fatalf("min > max = %d, want %d", status, http.StatusBadRequest)
		}
		// Missing decision type.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"approval_level": 1, "role_number": dRoleA,
		}); status != http.StatusBadRequest {
			t.Fatalf("missing decision type = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("unit_scope_validation", func(t *testing.T) {
		for _, unit := range []string{dUnitB, dUnitIn, "does-not-exist"} {
			if status, body := createAuthority(fx.tokenAdminA, map[string]any{
				"decision_type": "unit_request", "approval_level": 1, "role_number": dRoleA,
				"max_amount": 1000000, "organisation_unit_number": unit,
			}); status != http.StatusBadRequest {
				t.Fatalf("unit %s = %d, want %d (%s)", unit, status, http.StatusBadRequest, string(body))
			}
		}
	})

	t.Run("configure_purchase_request_ladder", func(t *testing.T) {
		// Level 1 (already created in lifecycle test is ENDED) — create fresh rows.
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("purchase_request", 1, f(500000))); status != http.StatusCreated {
			t.Fatalf("level 1 = %d (%s)", status, string(body))
		}
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("purchase_request", 2, f(5000000))); status != http.StatusCreated {
			t.Fatalf("level 2 = %d (%s)", status, string(body))
		}
		if status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "purchase_request", "approval_level": 3, "stakeholder_type": "owner",
		}); status != http.StatusCreated {
			t.Fatalf("level 3 = %d (%s)", status, string(body))
		}
	})

	t.Run("resolver_required_path_and_boundaries", func(t *testing.T) {
		cases := []struct {
			amount float64
			want   string
		}{
			{0, "[1]"},
			{100000, "[1]"},
			{500000, "[1]"},
			{500001, "[1 2]"},
			{2000000, "[1 2]"},
			{5000000, "[1 2]"},
			{5000001, "[1 2 3]"},
		}
		for _, testCase := range cases {
			status, envelope, body := resolve(fx.tokenAdminA, "purchase_request", testCase.amount, "")
			if status != http.StatusOK {
				t.Fatalf("resolve %v = %d (%s)", testCase.amount, status, string(body))
			}
			if !envelope.Data.Configured {
				t.Fatalf("purchase_request must be configured (%s)", string(body))
			}
			if envelope.Data.ExceedsDeclaredAuthority {
				t.Fatalf("purchase_request must not exceed declared authority (%s)", string(body))
			}
			if got := formatLevels(levels(envelope)); got != testCase.want {
				t.Fatalf("amount %v path = %s, want %s", testCase.amount, got, testCase.want)
			}
		}

		// Anchors are the exact contract Foundation E consumes.
		status, envelope, body := resolve(fx.tokenAdminA, "purchase_request", 5000001, "")
		if status != http.StatusOK {
			t.Fatalf("resolve 5000001 = %d (%s)", status, string(body))
		}
		if len(envelope.Data.RequiredPath) != 3 {
			t.Fatalf("expected a 3-level path, got %s", string(body))
		}
		first := envelope.Data.RequiredPath[0]
		if first.ApprovalLevel != 1 || len(first.Anchors) != 1 ||
			first.Anchors[0].Type != "role" || first.Anchors[0].Value != dRoleA {
			t.Fatalf("level 1 anchors must be the configured role: %s", string(body))
		}
		third := envelope.Data.RequiredPath[2]
		if third.ApprovalLevel != 3 || len(third.Anchors) != 1 ||
			third.Anchors[0].Type != "stakeholder_type" || third.Anchors[0].Value != "owner" {
			t.Fatalf("level 3 anchors must be the owner stakeholder type: %s", string(body))
		}
	})

	t.Run("configured_states_bd2", func(t *testing.T) {
		// State 1: nothing configured for this decision type anywhere.
		_, envelope, _ := resolve(fx.tokenAdminA, "never_configured_decision", 10, "")
		if envelope.Data.Configured || len(envelope.Data.RequiredPath) != 0 {
			t.Fatal("State 1: an unconfigured decision type must report configured=false with an empty path")
		}

		// State 2: configured and applicable to this request.
		_, envelope, _ = resolve(fx.tokenAdminA, "purchase_request", 100000, "")
		if !envelope.Data.Configured || len(envelope.Data.RequiredPath) == 0 {
			t.Fatal("State 2: a configured and applicable decision type must report configured=true with a path")
		}

		// State 3: governed, but nothing applies to this scope. A unit-scoped-only
		// configuration is created here so the assertion is self-contained.
		if status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "state3_request", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000000, "organisation_unit_number": dUnitA,
		}); status != http.StatusCreated {
			t.Fatalf("State 3 config = %d (%s)", status, string(body))
		}

		_, envelope, body := resolve(fx.tokenAdminA, "state3_request", 500000, dUnitOt)
		if !envelope.Data.Configured {
			t.Fatalf("State 3: configured must be true (%s)", string(body))
		}
		if len(envelope.Data.RequiredPath) != 0 {
			t.Fatalf("State 3: the path must be empty for a non-applicable scope (%s)", string(body))
		}
		if envelope.Data.ExceedsDeclaredAuthority {
			t.Fatalf("State 3: exceeds_declared_authority must be false when nothing applies to this scope (%s)", string(body))
		}
	})

	t.Run("unit_deactivation_does_not_revoke_authority_bd3", func(t *testing.T) {
		if status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "deact_request", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000000, "organisation_unit_number": dUnitDeact,
		}); status != http.StatusCreated {
			t.Fatalf("deact config = %d (%s)", status, string(body))
		}

		_, envelope, _ := resolve(fx.tokenAdminA, "deact_request", 500000, dUnitDeact)
		if !envelope.Data.Configured || len(envelope.Data.RequiredPath) != 1 {
			t.Fatal("authority must resolve while the unit is ACTIVE")
		}

		if status, body := request(http.MethodPatch, dUnitBase+"/"+dUnitDeact+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivating unit = %d (%s)", status, string(body))
		}

		// BD-3: deactivating the unit must not silently revoke the authority.
		_, envelope, body := resolve(fx.tokenAdminA, "deact_request", 500000, dUnitDeact)
		if !envelope.Data.Configured || len(envelope.Data.RequiredPath) != 1 {
			t.Fatalf("deactivating the unit must not revoke configured authority (BD-3): %s", string(body))
		}

		// But new authority may not be configured against the inactive unit.
		if status, _ := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "deact_request_two", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000, "organisation_unit_number": dUnitDeact,
		}); status != http.StatusBadRequest {
			t.Fatalf("configuring against an inactive unit = %d, want %d", status, http.StatusBadRequest)
		}

		// An unknown unit is still invalid for the resolver.
		if status, _, _ := resolve(fx.tokenAdminA, "deact_request", 500000, "does-not-exist"); status != http.StatusBadRequest {
			t.Fatalf("resolving with an unknown unit = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("resolver_exceeds_and_not_configured", func(t *testing.T) {
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("budget_change", 1, f(100000))); status != http.StatusCreated {
			t.Fatalf("budget level 1 = %d (%s)", status, string(body))
		}

		status, envelope, body := resolve(fx.tokenAdminA, "budget_change", 50000, "")
		if status != http.StatusOK || !envelope.Data.Configured || envelope.Data.ExceedsDeclaredAuthority {
			t.Fatalf("budget 50000 unexpected: %s", string(body))
		}

		status, envelope, body = resolve(fx.tokenAdminA, "budget_change", 100001, "")
		if status != http.StatusOK {
			t.Fatalf("budget 100001 = %d", status)
		}
		if !envelope.Data.ExceedsDeclaredAuthority {
			t.Fatalf("budget 100001 must exceed declared authority: %s", string(body))
		}
		if got := formatLevels(levels(envelope)); got != "[1]" {
			t.Fatalf("exceeds path = %s, want [1] (all configured levels)", got)
		}

		status, envelope, _ = resolve(fx.tokenAdminA, "travel_request", 10, "")
		if status != http.StatusOK {
			t.Fatalf("unconfigured resolve = %d", status)
		}
		if envelope.Data.Configured || len(envelope.Data.RequiredPath) != 0 {
			t.Fatalf("travel_request must be unconfigured with an empty path")
		}
	})

	t.Run("resolver_excludes_inactive_and_expired", func(t *testing.T) {
		// Inactive.
		status, body := createAuthority(fx.tokenAdminA, rolePayload("inactive_request", 1, f(1000)))
		if status != http.StatusCreated {
			t.Fatalf("inactive config = %d (%s)", status, string(body))
		}
		row := decodeAuthority(t, body)
		request(http.MethodPatch, dAuthBase+"/"+row.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"})

		_, envelope, _ := resolve(fx.tokenAdminA, "inactive_request", 10, "")
		if envelope.Data.Configured {
			t.Fatal("an INACTIVE authority row must not configure the decision type")
		}

		// Expired (valid_to in the past).
		past := "2020-01-01T00:00:00Z"
		if status, body = createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "expired_request", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000, "valid_from": "2019-01-01T00:00:00Z", "valid_to": past,
		}); status != http.StatusCreated {
			t.Fatalf("expired config = %d (%s)", status, string(body))
		}

		_, envelope, _ = resolve(fx.tokenAdminA, "expired_request", 10, "")
		if envelope.Data.Configured {
			t.Fatal("an expired authority row must not configure the decision type")
		}
	})

	t.Run("resolver_unit_subtree_scope", func(t *testing.T) {
		if status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "unit_request", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000000, "organisation_unit_number": dUnitA,
		}); status != http.StatusCreated {
			t.Fatalf("unit config = %d (%s)", status, string(body))
		}

		// BD-2: `configured` is scope-independent, so every case below reports
		// configured = true; only the applicable path length changes.
		for _, testCase := range []struct {
			unit    string
			pathLen int
		}{
			{dUnitA, 1},  // the scoped unit itself
			{dUnitCh, 1}, // a descendant (subtree)
			{dUnitOt, 0}, // unrelated unit: governed, but nothing applies here
			{"", 0},      // unit-less request: unit-scoped rows do not apply (BD-4)
		} {
			_, envelope, body := resolve(fx.tokenAdminA, "unit_request", 500000, testCase.unit)
			if !envelope.Data.Configured {
				t.Fatalf("unit %q must report configured=true (%s)", testCase.unit, string(body))
			}
			if len(envelope.Data.RequiredPath) != testCase.pathLen {
				t.Fatalf("unit %q path length = %d, want %d (%s)", testCase.unit, len(envelope.Data.RequiredPath), testCase.pathLen, string(body))
			}
		}

		// An organisation-wide row applies to any unit.
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("orgwide_request", 1, f(1000000))); status != http.StatusCreated {
			t.Fatalf("orgwide config = %d (%s)", status, string(body))
		}
		_, envelope, _ := resolve(fx.tokenAdminA, "orgwide_request", 500000, dUnitOt)
		if !envelope.Data.Configured {
			t.Fatal("an organisation-wide row must apply to any unit")
		}
	})

	t.Run("organisation_isolation_and_idor", func(t *testing.T) {
		_, listB := request(http.MethodGet, dAuthBase, fx.tokenAdminB, nil)
		if bytes.Contains(listB, []byte("purchase_request")) {
			t.Fatalf("org B must not see org A authority rows: %s", string(listB))
		}

		_, envelope, _ := resolve(fx.tokenAdminB, "purchase_request", 100000, "")
		if envelope.Data.Configured {
			t.Fatal("org B must not resolve org A authority")
		}

		// Cross-organisation IDOR on a known org A row.
		status, body := request(http.MethodGet, dAuthBase+"?decision_type=budget_change", fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("org A list = %d", status)
		}
		var listEnvelope struct {
			Data []struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &listEnvelope); err != nil || len(listEnvelope.Data) == 0 {
			t.Fatalf("could not read org A authority rows: %s", string(body))
		}
		number := listEnvelope.Data[0].Number

		if status, _ = request(http.MethodGet, dAuthBase+"/"+number, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodPatch, dAuthBase+"/"+number+"/status", fx.tokenAdminB, map[string]any{"status": "ENDED"}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("organisation_switching", func(t *testing.T) {
		if status, body := createAuthority(fx.tokenAdminB, map[string]any{
			"decision_type": "orgb_request", "approval_level": 1, "role_number": dRoleB, "max_amount": 100,
		}); status != http.StatusCreated {
			t.Fatalf("org B authority = %d (%s)", status, string(body))
		}

		_, envelopeA, _ := resolve(fx.tokenSwitch, "orgb_request", 10, "")
		if envelopeA.Data.Configured {
			t.Fatal("switcher in org A must not resolve org B authority")
		}

		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", dSwitcher).Update("entity_number", dOrgB).Error; err != nil {
			t.Fatalf("switching active entity: %v", err)
		}

		_, envelopeB, _ := resolve(fx.tokenSwitch, "orgb_request", 10, "")
		if !envelopeB.Data.Configured {
			t.Fatal("switcher in org B must resolve org B authority")
		}
	})

	t.Run("owner_and_permission_separation", func(t *testing.T) {
		// An owner STAKEHOLDER does not create authority: no authority row exists
		// for this decision type, so the resolver stays unconfigured.
		if status, body := request(http.MethodPost, "/main/ems/stakeholders", fx.tokenAdminA, map[string]any{
			"user_number": dOwnerUser, "stakeholder_type": "owner",
		}); status != http.StatusCreated {
			t.Fatalf("owner stakeholder = %d (%s)", status, string(body))
		}

		_, envelope, body := resolve(fx.tokenAdminA, "owner_decision", 1, "")
		if envelope.Data.Configured {
			t.Fatalf("ownership must not imply authority: %s", string(body))
		}

		// Authority management requires the permission; a user without it is refused.
		if status, _ := createAuthority(fx.tokenNoPerm, rolePayload("purchase_request", 9, f(1))); status != http.StatusForbidden {
			t.Fatalf("non-permitted authority creation = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("future_dated_authority_is_excluded", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

		status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "future_request", "approval_level": 1, "role_number": dRoleA,
			"max_amount": 1000000, "valid_from": future,
		})
		if status != http.StatusCreated {
			t.Fatalf("future-dated config = %d (%s)", status, string(body))
		}

		_, envelope, body := resolve(fx.tokenAdminA, "future_request", 1000, "")
		if envelope.Data.Configured || len(envelope.Data.RequiredPath) != 0 {
			t.Fatalf("a future-dated authority must not participate in resolution before its validity begins: %s", string(body))
		}
	})

	t.Run("cross_organisation_role_anchor_rejected", func(t *testing.T) {
		// dRoleB belongs to organisation B; adminA operates in organisation A.
		status, body := createAuthority(fx.tokenAdminA, map[string]any{
			"decision_type": "cross_role_request", "approval_level": 1, "role_number": dRoleB,
		})
		if status != http.StatusBadRequest {
			t.Fatalf("cross-organisation role anchor = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}
	})

	t.Run("duplicate_same_level_any_coverage", func(t *testing.T) {
		// Two ACTIVE rows at the same level; the second allows a larger amount.
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("dup_request", 1, f(100))); status != http.StatusCreated {
			t.Fatalf("duplicate row 1 = %d (%s)", status, string(body))
		}
		if status, body := createAuthority(fx.tokenAdminA, rolePayload("dup_request", 1, f(1000))); status != http.StatusCreated {
			t.Fatalf("duplicate row 2 = %d (%s)", status, string(body))
		}

		// ANY-coverage: the level covers 500 because one of its rows allows it.
		_, envelope, body := resolve(fx.tokenAdminA, "dup_request", 500, "")
		if !envelope.Data.Configured || len(envelope.Data.RequiredPath) != 1 {
			t.Fatalf("a configurable level must resolve: %s", string(body))
		}
		if envelope.Data.ExceedsDeclaredAuthority {
			t.Fatalf("ANY-coverage must cover 500 via the larger row: %s", string(body))
		}
		// Identical anchors on rows at the same level are de-duplicated.
		if anchors := envelope.Data.RequiredPath[0].Anchors; len(anchors) != 1 || anchors[0].Type != "role" || anchors[0].Value != dRoleA {
			t.Fatalf("duplicate rows must merge into one level with de-duplicated anchors: %s", string(body))
		}

		// Above every row at the level, with no further level configured.
		_, envelope, body = resolve(fx.tokenAdminA, "dup_request", 1001, "")
		if !envelope.Data.ExceedsDeclaredAuthority || len(envelope.Data.RequiredPath) != 1 {
			t.Fatalf("1001 must exceed the level's rows: %s", string(body))
		}
	})

	t.Run("unknown_unit_during_resolution", func(t *testing.T) {
		status, _, body := resolve(fx.tokenAdminA, "purchase_request", 100000, "does-not-exist-unit")
		if status != http.StatusBadRequest {
			t.Fatalf("unknown unit = %d, want %d; it must not be treated as organisation-wide (%s)", status, http.StatusBadRequest, string(body))
		}
	})

	t.Run("non_finite_amounts_rejected", func(t *testing.T) {
		for _, amount := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			status, _, _ := resolve(fx.tokenAdminA, "purchase_request", amount, "")
			if status != http.StatusBadRequest {
				t.Fatalf("amount %v = %d, want %d", amount, status, http.StatusBadRequest)
			}
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, dEduURL, fx.tokenAdminA, nil)

		status, body := createAuthority(fx.tokenAdminA, rolePayload("scope_check", 1, f(1000)))
		if status != http.StatusCreated {
			t.Fatalf("scope-check authority = %d (%s)", status, string(body))
		}
		row := decodeAuthority(t, body)
		request(http.MethodPatch, dAuthBase+"/"+row.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})

		statusAfter, bodyAfter := request(http.MethodGet, dEduURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS authority changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS authority changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}

func formatAmount(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatLevels(levels []int) string {
	parts := make([]string, 0, len(levels))
	for _, level := range levels {
		parts = append(parts, strconv.Itoa(level))
	}
	return "[" + strings.Join(parts, " ") + "]"
}
