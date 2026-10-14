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
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	c3AdminA    = "ems-c3-admin-a"
	c3AdminB    = "ems-c3-admin-b"
	c3Switcher  = "ems-c3-switcher"
	c3NoPerm    = "ems-c3-noperm"
	c3NoCtx     = "ems-c3-nocontext"
	c3User      = "ems-c3-user"        // assignable, no stakeholder
	c3UserAlt   = "ems-c3-user-alt"    // assignable
	c3UserNoLog = "ems-c3-user-noauth" // inactive, may still be assigned
	c3Bare      = "ems-c3-bare"        // used for "stakeholder without assignment"

	c3OrgA = "ems-c3-org-a"
	c3OrgB = "ems-c3-org-b"

	c3UnitA        = "ems-c3-unit-a"
	c3UnitB        = "ems-c3-unit-b"
	c3UnitInactive = "ems-c3-unit-inactive"
	c3UnitGuardAct = "ems-c3-unit-guard-active"
	c3UnitGuardIna = "ems-c3-unit-guard-inactive"
	c3UnitDelAct   = "ems-c3-unit-del-active"
	c3UnitDelIna   = "ems-c3-unit-del-inactive"
	c3UnitDelEnd   = "ems-c3-unit-del-ended"
	c3UnitScope    = "ems-c3-unit-scope"

	c3BasePath     = "/main/ems/user-assignments"
	c3UnitBasePath = "/main/ems/organisation-units"
	c3StakePath    = "/main/ems/stakeholders"
	c3EduURL       = "/main/campus/all"
)

type c3Fixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenSwitch string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupC3Fixture(t *testing.T) *c3Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-c3-test-secret")

	users := []string{
		c3AdminA, c3AdminB, c3Switcher, c3NoPerm, c3NoCtx,
		c3User, c3UserAlt, c3UserNoLog, c3Bare,
	}
	organisations := []string{c3OrgA, c3OrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
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
		{c3AdminA, "Admin", "Active"},
		{c3AdminB, "Admin", "Active"},
		{c3Switcher, "Admin", "Active"},
		{c3NoPerm, "User", "Active"},
		{c3NoCtx, "Admin", "Active"},
		{c3User, "User", "Active"},
		{c3UserAlt, "User", "Active"},
		{c3UserNoLog, "User", "Inactive"},
		{c3Bare, "User", "Active"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "C3",
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
		{c3OrgA, "C3 Test Org A"},
		{c3OrgB, "C3 Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: c3AdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c3AdminA, c3OrgA},
		{c3AdminB, c3OrgB},
		{c3Switcher, c3OrgA},
		{c3NoPerm, c3OrgA},
		// c3NoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{c3Switcher, c3OrgA},
		{c3Switcher, c3OrgB},
	} {
		membership := eosModel.EntityUser{EntityNumber: spec.entityNumber, UserNumber: spec.userNumber, AddedByNumber: c3AdminA}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	for _, spec := range []struct {
		number, organisation, name, status string
	}{
		{c3UnitA, c3OrgA, "Unit A", "ACTIVE"},
		{c3UnitB, c3OrgB, "Unit B", "ACTIVE"},
		{c3UnitInactive, c3OrgA, "Inactive Unit", "INACTIVE"},
		{c3UnitGuardAct, c3OrgA, "Guard Active Unit", "ACTIVE"},
		{c3UnitGuardIna, c3OrgA, "Guard Inactive Unit", "ACTIVE"},
		{c3UnitDelAct, c3OrgA, "Delete Guard Active Unit", "ACTIVE"},
		{c3UnitDelIna, c3OrgA, "Delete Guard Inactive Unit", "ACTIVE"},
		{c3UnitDelEnd, c3OrgA, "Delete Guard Ended Unit", "ACTIVE"},
		{c3UnitScope, c3OrgA, "Scope Isolation Unit", "ACTIVE"},
	} {
		unit := organisationModel.OrganisationUnit{
			Number:             spec.number,
			OrganisationNumber: spec.organisation,
			UnitType:           "department",
			Name:               spec.name,
			Status:             spec.status,
			AddedByNumber:      c3AdminA,
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

	return &c3Fixture{
		db:          db,
		tokenAdminA: token(c3AdminA),
		tokenAdminB: token(c3AdminB),
		tokenSwitch: token(c3Switcher),
		tokenNoPerm: token(c3NoPerm),
		tokenNoCtx:  token(c3NoCtx),
	}
}

type c3Envelope struct {
	Message string `json:"message"`
	Data    struct {
		Number                 string  `json:"number"`
		UserNumber             string  `json:"user_number"`
		OrganisationUnitNumber string  `json:"organisation_unit_number"`
		AssignmentType         string  `json:"assignment_type"`
		PositionTitle          string  `json:"position_title"`
		Status                 string  `json:"status"`
		ValidTo                *string `json:"valid_to"`
	} `json:"data"`
}

// TestUserAssignments covers Foundation C3.
func TestUserAssignments(t *testing.T) {
	fx := setupC3Fixture(t)

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

	create := func(token, userNumber, unitNumber, assignmentType, position string) (int, []byte) {
		return request(http.MethodPost, c3BasePath, token, map[string]any{
			"user_number":              userNumber,
			"organisation_unit_number": unitNumber,
			"assignment_type":          assignmentType,
			"position_title":           position,
		})
	}

	decode := func(t *testing.T, body []byte) c3Envelope {
		t.Helper()
		var envelope c3Envelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding %s: %v", string(body), err)
		}
		return envelope
	}

	countAssignments := func(userNumber, unitNumber string) int64 {
		t.Helper()
		var count int64
		query := fx.db.Model(&assignmentModel.UserAssignment{}).Where("organisation_number = ? AND user_number = ?", c3OrgA, userNumber)
		if unitNumber != "" {
			query = query.Where("organisation_unit_number = ?", unitNumber)
		}
		if err := query.Count(&count).Error; err != nil {
			t.Fatalf("counting assignments: %v", err)
		}
		return count
	}

	t.Run("authentication_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c3BasePath, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("permission_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, c3BasePath, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, c3BasePath, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("no-context = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("create_and_retrieve_assignment", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c3User, c3UnitA, "member", "Analyst")
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decode(t, body)
		if created.Data.Status != "ACTIVE" || created.Data.PositionTitle != "Analyst" {
			t.Fatalf("unexpected assignment: %s", string(body))
		}

		if status, _ = request(http.MethodGet, c3BasePath+"/"+created.Data.Number, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("retrieve = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("update_assignment", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c3User, c3UnitA, "manager", "Manager")
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		created := decode(t, body)

		status, body = request(http.MethodPut, c3BasePath+"/"+created.Data.Number, fx.tokenAdminA, map[string]any{
			"organisation_unit_number": c3UnitA,
			"assignment_type":          "head",
			"position_title":           "Head of Unit A",
		})
		if status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}
		if updated := decode(t, body); updated.Data.PositionTitle != "Head of Unit A" || updated.Data.AssignmentType != "head" {
			t.Fatalf("update not applied: %s", string(body))
		}
	})

	t.Run("lifecycle_inactive_reactivate_end", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c3User, c3UnitA, "staff", "Lifecycle")
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		row := decode(t, body)
		path := c3BasePath + "/" + row.Data.Number

		for _, step := range []struct{ value, want string }{
			{"INACTIVE", "INACTIVE"},
			{"ACTIVE", "ACTIVE"},
		} {
			status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]any{"status": step.value})
			if status != http.StatusOK {
				t.Fatalf("status %s = %d (%s)", step.value, status, string(body))
			}
			if decoded := decode(t, body); decoded.Data.Status != step.want {
				t.Fatalf("status = %q, want %q", decoded.Data.Status, step.want)
			}
		}

		status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})
		if status != http.StatusOK {
			t.Fatalf("end = %d (%s)", status, string(body))
		}
		ended := decode(t, body)
		if ended.Data.Status != "ENDED" || ended.Data.ValidTo == nil {
			t.Fatalf("ending must set ENDED and stamp valid_to: %s", string(body))
		}
	})

	t.Run("invalid_date_range_and_inputs", func(t *testing.T) {
		status, _ := request(http.MethodPost, c3BasePath, fx.tokenAdminA, map[string]any{
			"user_number":              c3UserAlt,
			"organisation_unit_number": c3UnitA,
			"assignment_type":          "member",
			"valid_from":               "2026-06-01T00:00:00Z",
			"valid_to":                 "2026-01-01T00:00:00Z",
		})
		if status != http.StatusBadRequest {
			t.Fatalf("reversed date range = %d, want %d", status, http.StatusBadRequest)
		}

		if status, _ := create(fx.tokenAdminA, c3UserAlt, c3UnitA, "superhero", ""); status != http.StatusBadRequest {
			t.Fatalf("invalid type = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := create(fx.tokenAdminA, "does-not-exist", c3UnitA, "member", ""); status != http.StatusBadRequest {
			t.Fatalf("unknown user = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("multiple_assignments_and_duplicate", func(t *testing.T) {
		if status, body := create(fx.tokenAdminA, c3UserAlt, c3UnitA, "member", ""); status != http.StatusCreated {
			t.Fatalf("first = %d (%s)", status, string(body))
		}
		if status, _ := create(fx.tokenAdminA, c3UserAlt, c3UnitA, "member", ""); status != http.StatusConflict {
			t.Fatalf("duplicate = %d, want %d", status, http.StatusConflict)
		}
		// A different unit is allowed.
		if status, body := create(fx.tokenAdminA, c3UserAlt, c3UnitGuardAct, "manager", ""); status != http.StatusCreated {
			t.Fatalf("second unit = %d (%s)", status, string(body))
		}
		if total := countAssignments(c3UserAlt, ""); total < 2 {
			t.Fatalf("expected multiple assignments for %s, found %d", c3UserAlt, total)
		}
	})

	t.Run("unit_rules", func(t *testing.T) {
		if status, _ := create(fx.tokenAdminA, c3User, c3UnitInactive, "member", ""); status != http.StatusBadRequest {
			t.Fatalf("inactive unit = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := create(fx.tokenAdminA, c3User, c3UnitB, "member", ""); status != http.StatusBadRequest {
			t.Fatalf("cross-organisation unit = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := create(fx.tokenAdminA, c3User, "does-not-exist", "member", ""); status != http.StatusBadRequest {
			t.Fatalf("unknown unit = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("user_rules", func(t *testing.T) {
		// (14) user without a stakeholder can be assigned.
		var stakeholders int64
		if err := fx.db.Model(&stakeholderModel.Stakeholder{}).Where("user_number = ?", c3User).Count(&stakeholders).Error; err != nil {
			t.Fatalf("counting stakeholders: %v", err)
		}
		if stakeholders != 0 {
			t.Fatalf("%s should have no stakeholder relationship for this test", c3User)
		}

		// (16) an inactive user may still hold an assignment.
		if status, body := create(fx.tokenAdminA, c3UserNoLog, c3UnitA, "member", ""); status != http.StatusCreated {
			t.Fatalf("inactive user assignment = %d (%s)", status, string(body))
		}
	})

	t.Run("stakeholder_without_assignment", func(t *testing.T) {
		if status, body := request(http.MethodPost, c3StakePath, fx.tokenAdminA, map[string]any{
			"user_number": c3Bare, "stakeholder_type": "staff",
		}); status != http.StatusCreated {
			t.Fatalf("stakeholder create = %d (%s)", status, string(body))
		}
		if count := countAssignments(c3Bare, ""); count != 0 {
			t.Fatalf("stakeholder %s must have no assignment, found %d", c3Bare, count)
		}
	})

	t.Run("organisation_isolation_and_idor", func(t *testing.T) {
		status, body := create(fx.tokenAdminA, c3User, c3UnitGuardAct, "member", "Isolation")
		if status != http.StatusCreated {
			t.Fatalf("org A create = %d (%s)", status, string(body))
		}
		row := decode(t, body)
		path := c3BasePath + "/" + row.Data.Number

		if status, _ = request(http.MethodGet, path, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodPut, path, fx.tokenAdminB, map[string]any{
			"organisation_unit_number": c3UnitB, "assignment_type": "member",
		}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation update = %d, want %d", status, http.StatusNotFound)
		}

		_, listB := request(http.MethodGet, c3BasePath, fx.tokenAdminB, nil)
		if bytes.Contains(listB, []byte(row.Data.Number)) {
			t.Fatalf("org B list must not contain org A assignment: %s", string(listB))
		}
	})

	t.Run("organisation_switching", func(t *testing.T) {
		statusA, bodyA := create(fx.tokenAdminA, c3User, c3UnitGuardIna, "member", "Switch A")
		if statusA != http.StatusCreated {
			t.Fatalf("org A = %d (%s)", statusA, string(bodyA))
		}
		rowA := decode(t, bodyA)

		statusB, bodyB := create(fx.tokenAdminB, c3User, c3UnitB, "member", "Switch B")
		if statusB != http.StatusCreated {
			t.Fatalf("org B = %d (%s)", statusB, string(bodyB))
		}
		rowB := decode(t, bodyB)

		_, listBefore := request(http.MethodGet, c3BasePath, fx.tokenSwitch, nil)
		if !bytes.Contains(listBefore, []byte(rowA.Data.Number)) || bytes.Contains(listBefore, []byte(rowB.Data.Number)) {
			t.Fatalf("switcher in org A must see only org A assignments: %s", string(listBefore))
		}

		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", c3Switcher).Update("entity_number", c3OrgB).Error; err != nil {
			t.Fatalf("switching active entity: %v", err)
		}

		_, listAfter := request(http.MethodGet, c3BasePath, fx.tokenSwitch, nil)
		if !bytes.Contains(listAfter, []byte(rowB.Data.Number)) || bytes.Contains(listAfter, []byte(rowA.Data.Number)) {
			t.Fatalf("switcher in org B must see only org B assignments: %s", string(listAfter))
		}
	})

	t.Run("unit_deletion_guard", func(t *testing.T) {
		// ACTIVE assignment blocks unit deletion.
		status, body := create(fx.tokenAdminA, c3User, c3UnitDelAct, "head", "Guard")
		if status != http.StatusCreated {
			t.Fatalf("guard assignment = %d (%s)", status, string(body))
		}
		activeAssignment := decode(t, body)

		if status, _ = request(http.MethodDelete, c3UnitBasePath+"/"+c3UnitDelAct, fx.tokenAdminA, nil); status != http.StatusConflict {
			t.Fatalf("delete unit with ACTIVE assignment = %d, want %d", status, http.StatusConflict)
		}

		// Ending the assignment unblocks deletion.
		if status, _ = request(http.MethodPatch, c3BasePath+"/"+activeAssignment.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("ending assignment = %d, want %d", status, http.StatusOK)
		}
		if status, _ = request(http.MethodDelete, c3UnitBasePath+"/"+c3UnitDelAct, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("delete unit after ENDED assignment = %d, want %d", status, http.StatusOK)
		}

		// INACTIVE assignments do not block deletion.
		status, body = create(fx.tokenAdminA, c3User, c3UnitDelIna, "staff", "Inactive guard")
		if status != http.StatusCreated {
			t.Fatalf("inactive-guard assignment = %d (%s)", status, string(body))
		}
		inactiveAssignment := decode(t, body)
		if status, _ = request(http.MethodPatch, c3BasePath+"/"+inactiveAssignment.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivating assignment = %d, want %d", status, http.StatusOK)
		}
		if status, _ = request(http.MethodDelete, c3UnitBasePath+"/"+c3UnitDelIna, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("delete unit with INACTIVE assignment = %d, want %d", status, http.StatusOK)
		}

		// ENDED assignments do not block deletion.
		status, body = create(fx.tokenAdminA, c3User, c3UnitDelEnd, "member", "Ended guard")
		if status != http.StatusCreated {
			t.Fatalf("ended-guard assignment = %d (%s)", status, string(body))
		}
		endedAssignment := decode(t, body)
		if status, _ = request(http.MethodPatch, c3BasePath+"/"+endedAssignment.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("ending guard assignment = %d, want %d", status, http.StatusOK)
		}
		if status, _ = request(http.MethodDelete, c3UnitBasePath+"/"+c3UnitDelEnd, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("delete unit with ENDED assignment = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("no_duplicate_identity", func(t *testing.T) {
		var before int64
		if err := fx.db.Model(&authModel.User{}).Count(&before).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}

		if status, body := create(fx.tokenAdminA, c3Bare, c3UnitA, "member", ""); status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}

		var after int64
		if err := fx.db.Model(&authModel.User{}).Count(&after).Error; err != nil {
			t.Fatalf("counting users: %v", err)
		}
		if before != after {
			t.Fatalf("user count changed from %d to %d; assignments must not create identities", before, after)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, c3EduURL, fx.tokenAdminA, nil)

		status, body := create(fx.tokenAdminA, c3User, c3UnitScope, "member", "Scope check")
		if status != http.StatusCreated {
			t.Fatalf("scope-check create = %d (%s)", status, string(body))
		}
		row := decode(t, body)
		request(http.MethodPatch, c3BasePath+"/"+row.Data.Number+"/status", fx.tokenAdminA, map[string]any{"status": "ENDED"})

		statusAfter, bodyAfter := request(http.MethodGet, c3EduURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS assignment changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS assignment changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}
