package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	reportingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	s4OrgA    = "stage4-org-a"
	s4OrgB    = "stage4-org-b"
	s4OrgLife = "stage4-org-life"

	s4AdminA    = "stage4-admin-a"
	s4AdminB    = "stage4-admin-b"
	s4AdminLife = "stage4-admin-life"

	s4UnitA    = "stage4-unit-a"
	s4UnitB    = "stage4-unit-b"
	s4UnitLife = "stage4-unit-life"

	s4ReportingPath = "/main/ems/reporting-lines"
)

type s4Fixture struct {
	db           *gorm.DB
	tokenA       string
	tokenB       string
	tokenLife    string
	assignments  map[string]string
	lifeLineType string
}

func setupS4Fixture(t *testing.T) *s4Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "stage-4-test-secret")

	organisations := []string{s4OrgA, s4OrgB, s4OrgLife}
	users := []string{
		s4AdminA, s4AdminB, s4AdminLife,
		"stage4-person-1", "stage4-person-2", "stage4-person-3", "stage4-person-4",
		"stage4-person-5", "stage4-person-6", "stage4-person-7", "stage4-person-8",
		"stage4-person-9", "stage4-person-10", "stage4-person-11", "stage4-person-12",
		"stage4-person-13", "stage4-person-14", "stage4-person-15", "stage4-person-16",
		"stage4-person-17", "stage4-person-18", "stage4-person-19", "stage4-person-20",
		"stage4-person-21", "stage4-person-22", "stage4-person-23", "stage4-person-24",
		"stage4-person-25", "stage4-person-26", "stage4-person-27", "stage4-person-28",
		"stage4-person-29", "stage4-person-30",
	}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&reportingModel.ReportingLine{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.Licence{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	admins := map[string]string{s4AdminA: s4OrgA, s4AdminB: s4OrgB, s4AdminLife: s4OrgLife}
	for _, number := range users {
		userType := "User"
		if _, isAdmin := admins[number]; isAdmin {
			userType = "Admin"
		}
		user := authModel.User{
			Number: number, FirstName: "Stage4", SecondName: number, Gender: "M",
			Type: userType, Status: "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, organisation := range organisations {
		entity := eosModel.Entity{Number: organisation, Name: "Stage 4 " + organisation, AddedBy: s4AdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", organisation, err)
		}
	}

	for _, spec := range []struct{ user, entity string }{
		{s4AdminA, s4OrgA}, {s4AdminB, s4OrgB}, {s4AdminLife, s4OrgLife},
	} {
		membership := eosModel.EntityUser{EntityNumber: spec.entity, UserNumber: spec.user, AddedByNumber: spec.user}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership: %v", err)
		}
		active := eosModel.ActiveUserEntity{UserNumber: spec.user, EntityNumber: spec.entity}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active organisation: %v", err)
		}
	}

	for _, spec := range []struct{ number, org, name string }{
		{s4UnitA, s4OrgA, "Stage 4 Unit A"},
		{s4UnitB, s4OrgB, "Stage 4 Unit B"},
		{s4UnitLife, s4OrgLife, "Stage 4 Unit Life"},
	} {
		unit := organisationModel.OrganisationUnit{
			Number: spec.number, OrganisationNumber: spec.org, Name: spec.name,
			UnitType: "division", Status: "ACTIVE", AddedByNumber: s4AdminA,
		}
		if err := db.Create(&unit).Error; err != nil {
			t.Fatalf("seeding unit %s: %v", spec.number, err)
		}
	}

	now := time.Now()
	lifeOrg := s4OrgLife
	licence := licensingModel.Licence{
		Number: "stage4-licence-life", LicenceKey: "stage4-licence-key-life",
		OrganisationNumber: s4OrgLife, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(24 * time.Hour),
		GraceDays: licensingModel.DefaultGraceDays, MaxInstallations: 1,
		LicensedAddons: `["education"]`, Status: licensingModel.StatusActive,
		ActiveOrgNumber: &lifeOrg, IssuedBy: s4AdminLife,
	}
	if err := db.Create(&licence).Error; err != nil {
		t.Fatalf("seeding licence: %v", err)
	}

	assignments := []assignmentModel.UserAssignment{
		{Number: "stage4-assign-sub1", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-1", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Analyst", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-mgr1", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-2", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Operations Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-mgr2", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-3", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Finance Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-sub2", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-4", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-sub3", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-5", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-sub4", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-6", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-mgr-b", OrganisationNumber: s4OrgB, UserNumber: "stage4-person-7", OrganisationUnitNumber: s4UnitB, AssignmentType: "manager", PositionTitle: "Manager B", Status: "ACTIVE", AddedByNumber: s4AdminB},
		{Number: "stage4-assign-cycle-a", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-8", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Cycle A", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-cycle-b", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-9", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Cycle B", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-cycle-c", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-10", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Cycle C", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-history", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-11", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-conc", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-12", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage4-assign-life-sub", OrganisationNumber: s4OrgLife, UserNumber: "stage4-person-1", OrganisationUnitNumber: s4UnitLife, AssignmentType: "member", PositionTitle: "Officer", Status: "ACTIVE", AddedByNumber: s4AdminLife},
		{Number: "stage4-assign-life-mgr", OrganisationNumber: s4OrgLife, UserNumber: "stage4-person-2", OrganisationUnitNumber: s4UnitLife, AssignmentType: "manager", PositionTitle: "Manager", Status: "ACTIVE", AddedByNumber: s4AdminLife},

		// Stage 5: isolated assignment pairs for effective-state tests.
		{Number: "stage5-assign-eff-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-13", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Effective Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-eff-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-14", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Effective Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-future-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-15", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Future Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-future-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-16", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Future Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-expired-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-17", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Expired Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-expired-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-18", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Expired Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-subended-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-19", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Sub Ended Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-subended-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-20", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Sub Ended Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-mgrended-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-21", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Mgr Ended Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-mgrended-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-22", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Mgr Ended Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-react-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-23", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "React Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-react-mgr", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-24", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "React Manager", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-dup-sub", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-25", OrganisationUnitNumber: s4UnitA, AssignmentType: "member", PositionTitle: "Dup Sub", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-dup-mgr1", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-26", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Dup Manager 1", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-dup-mgr2", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-27", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Dup Manager 2", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-hcyc-a", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-28", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Historical Cycle A", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-hcyc-b", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-29", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Historical Cycle B", Status: "ACTIVE", AddedByNumber: s4AdminA},
		{Number: "stage5-assign-hcyc-c", OrganisationNumber: s4OrgA, UserNumber: "stage4-person-30", OrganisationUnitNumber: s4UnitA, AssignmentType: "manager", PositionTitle: "Historical Cycle C", Status: "ACTIVE", AddedByNumber: s4AdminA},
	}
	for _, assignment := range assignments {
		if err := db.Omit("User", "OrganisationUnit").Create(&assignment).Error; err != nil {
			t.Fatalf("seeding assignment %s: %v", assignment.Number, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token: %v", err)
		}
		return value
	}

	return &s4Fixture{
		db:        db,
		tokenA:    token(s4AdminA),
		tokenB:    token(s4AdminB),
		tokenLife: token(s4AdminLife),
	}
}

func TestReportingLines(t *testing.T) {
	fx := setupS4Fixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	doRequest := func(method, path, token string, payload any) (int, []byte, error) {
		var reader io.Reader
		if payload != nil {
			raw, err := json.Marshal(payload)
			if err != nil {
				return 0, nil, err
			}
			reader = bytes.NewReader(raw)
		}
		req, err := http.NewRequest(method, server.URL+path, reader)
		if err != nil {
			return 0, nil, err
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body, nil
	}

	request := func(method, path, token string, payload any) (int, []byte) {
		t.Helper()
		status, body, err := doRequest(method, path, token, payload)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		return status, body
	}

	linePayload := func(subordinate, manager string, from, to *time.Time) map[string]any {
		payload := map[string]any{
			"subordinate_assignment_number": subordinate,
			"manager_assignment_number":     manager,
		}
		if from != nil {
			payload["valid_from"] = from.UTC().Format(time.RFC3339)
		}
		if to != nil {
			payload["valid_to"] = to.UTC().Format(time.RFC3339)
		}
		return payload
	}

	createLine := func(token, subordinate, manager string, from, to *time.Time) (int, []byte, string) {
		t.Helper()
		status, body := request(http.MethodPost, s4ReportingPath, token, linePayload(subordinate, manager, from, to))
		if status != http.StatusCreated {
			return status, body, ""
		}
		var parsed struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding reporting line: %v", err)
		}
		return status, body, parsed.Data.Number
	}

	type lineEnvelope struct {
		Data struct {
			Number         string `json:"number"`
			Status         string `json:"status"`
			Effective      bool   `json:"effective"`
			EffectiveState string `json:"effective_state"`
		} `json:"data"`
	}
	decodeLine := func(t *testing.T, body []byte) lineEnvelope {
		t.Helper()
		var envelope lineEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding reporting line: %v", err)
		}
		return envelope
	}

	t.Run("create_and_list", func(t *testing.T) {
		from := time.Now().Add(-time.Hour)
		to := time.Now().Add(24 * time.Hour)
		status, body, number := createLine(fx.tokenA, "stage4-assign-sub1", "stage4-assign-mgr1", &from, &to)
		if status != http.StatusCreated || number == "" {
			t.Fatalf("create = %d (%s)", status, string(body))
		}

		status, body = request(http.MethodGet, s4ReportingPath, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("list = %d (%s)", status, string(body))
		}
		var parsed struct {
			Data []struct {
				Number                      string `json:"number"`
				SubordinateAssignmentNumber string `json:"subordinate_assignment_number"`
				ManagerAssignmentNumber     string `json:"manager_assignment_number"`
				Status                      string `json:"status"`
				Effective                   bool   `json:"effective"`
				EffectiveState              string `json:"effective_state"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding list: %v", err)
		}
		found := false
		for _, item := range parsed.Data {
			if item.Number == number {
				found = true
				if item.SubordinateAssignmentNumber != "stage4-assign-sub1" || item.ManagerAssignmentNumber != "stage4-assign-mgr1" || item.Status != "ACTIVE" {
					t.Fatalf("unexpected line: %+v", item)
				}
				if !item.Effective || item.EffectiveState != "EFFECTIVE" {
					t.Fatalf("unexpected effective state: %+v", item)
				}
			}
		}
		if !found {
			t.Fatalf("created line %s missing from list", number)
		}
	})

	t.Run("self_reporting_rejected", func(t *testing.T) {
		status, body, _ := createLine(fx.tokenA, "stage4-assign-sub2", "stage4-assign-sub2", nil, nil)
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}
	})

	t.Run("cross_organisation_rejected", func(t *testing.T) {
		status, body, _ := createLine(fx.tokenA, "stage4-assign-sub3", "stage4-assign-mgr-b", nil, nil)
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusNotFound, string(body))
		}
	})

	t.Run("duplicate_active_manager_rejected", func(t *testing.T) {
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-sub4", "stage4-assign-mgr1", nil, nil); status != http.StatusCreated {
			t.Fatalf("first create = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-sub4", "stage4-assign-mgr2", nil, nil); status != http.StatusConflict {
			t.Fatalf("second create = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("cycle_rejected", func(t *testing.T) {
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-cycle-a", "stage4-assign-cycle-b", nil, nil); status != http.StatusCreated {
			t.Fatalf("A->B = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-cycle-b", "stage4-assign-cycle-c", nil, nil); status != http.StatusCreated {
			t.Fatalf("B->C = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-cycle-c", "stage4-assign-cycle-a", nil, nil); status != http.StatusConflict {
			t.Fatalf("C->A = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("invalid_date_range_rejected", func(t *testing.T) {
		from := time.Now().Add(24 * time.Hour)
		to := time.Now()
		status, body, _ := createLine(fx.tokenA, "stage4-assign-sub2", "stage4-assign-mgr1", &from, &to)
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}
	})

	t.Run("status_history_and_reactivation", func(t *testing.T) {
		status, body, number := createLine(fx.tokenA, "stage4-assign-history", "stage4-assign-mgr1", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}

		if status, body := request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("end = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodGet, s4ReportingPath+"/"+number, fx.tokenA, nil); status != http.StatusOK {
			t.Fatalf("get ended = %d (%s)", status, string(body))
		}
		// Ending stamps valid_to; reactivation without a new date decision is refused.
		if status, body := request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ACTIVE"}); status != http.StatusBadRequest {
			t.Fatalf("reactivate expired = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}

		from := time.Now().Add(-time.Hour)
		to := time.Now().Add(24 * time.Hour)
		if status, body := request(http.MethodPut, s4ReportingPath+"/"+number, fx.tokenA, map[string]any{
			"manager_assignment_number": "stage4-assign-mgr1",
			"valid_from":                from.Format(time.RFC3339),
			"valid_to":                  to.Format(time.RFC3339),
		}); status != http.StatusOK {
			t.Fatalf("update dates = %d (%s)", status, string(body))
		}
		status, body = request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("reactivate with dates = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); !decoded.Data.Effective || decoded.Data.EffectiveState != "EFFECTIVE" {
			t.Fatalf("unexpected reactivated state: %s", string(body))
		}
		if status, body := request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("end again = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage4-assign-history", "stage4-assign-mgr2", nil, nil); status != http.StatusCreated {
			t.Fatalf("new line after ending = %d (%s)", status, string(body))
		}
	})

	t.Run("effective_state_active", func(t *testing.T) {
		status, body, _ := createLine(fx.tokenA, "stage5-assign-eff-sub", "stage5-assign-eff-mgr", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); !decoded.Data.Effective || decoded.Data.EffectiveState != "EFFECTIVE" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("effective_state_future", func(t *testing.T) {
		from := time.Now().Add(24 * time.Hour)
		status, body, _ := createLine(fx.tokenA, "stage5-assign-future-sub", "stage5-assign-future-mgr", &from, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); decoded.Data.Effective || decoded.Data.EffectiveState != "UPCOMING" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("effective_state_expired", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)
		status, body, _ := createLine(fx.tokenA, "stage5-assign-expired-sub", "stage5-assign-expired-mgr", &from, &to)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); decoded.Data.Effective || decoded.Data.EffectiveState != "EXPIRED" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("ended_subordinate_assignment_is_not_effective", func(t *testing.T) {
		status, body, number := createLine(fx.tokenA, "stage5-assign-subended-sub", "stage5-assign-subended-mgr", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if err := fx.db.Model(&assignmentModel.UserAssignment{}).
			Where("number = ?", "stage5-assign-subended-sub").
			Update("status", "INACTIVE").Error; err != nil {
			t.Fatalf("ending subordinate assignment: %v", err)
		}
		status, body = request(http.MethodGet, s4ReportingPath+"/"+number, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("get = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); decoded.Data.Effective || decoded.Data.EffectiveState != "ASSIGNMENT_INACTIVE" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("ended_manager_assignment_is_not_effective", func(t *testing.T) {
		status, body, number := createLine(fx.tokenA, "stage5-assign-mgrended-sub", "stage5-assign-mgrended-mgr", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if err := fx.db.Model(&assignmentModel.UserAssignment{}).
			Where("number = ?", "stage5-assign-mgrended-mgr").
			Update("status", "INACTIVE").Error; err != nil {
			t.Fatalf("ending manager assignment: %v", err)
		}
		status, body = request(http.MethodGet, s4ReportingPath+"/"+number, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("get = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); decoded.Data.Effective || decoded.Data.EffectiveState != "ASSIGNMENT_INACTIVE" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("reactivation_expired_requires_dates", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)
		status, body, number := createLine(fx.tokenA, "stage5-assign-react-sub", "stage5-assign-react-mgr", &from, &to)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("end = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ACTIVE"}); status != http.StatusBadRequest {
			t.Fatalf("reactivate expired = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}

		newFrom := time.Now().Add(-time.Hour)
		newTo := time.Now().Add(24 * time.Hour)
		if status, body := request(http.MethodPut, s4ReportingPath+"/"+number, fx.tokenA, map[string]any{
			"manager_assignment_number": "stage5-assign-react-mgr",
			"valid_from":                newFrom.Format(time.RFC3339),
			"valid_to":                  newTo.Format(time.RFC3339),
		}); status != http.StatusOK {
			t.Fatalf("update dates = %d (%s)", status, string(body))
		}
		status, body = request(http.MethodPatch, s4ReportingPath+"/"+number+"/status", fx.tokenA, map[string]string{"status": "ACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("reactivate with dates = %d (%s)", status, string(body))
		}
		if decoded := decodeLine(t, body); !decoded.Data.Effective || decoded.Data.EffectiveState != "EFFECTIVE" {
			t.Fatalf("unexpected state: %s", string(body))
		}
	})

	t.Run("duplicate_expired_relationship_does_not_block", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)
		if status, body, _ := createLine(fx.tokenA, "stage5-assign-dup-sub", "stage5-assign-dup-mgr1", &from, &to); status != http.StatusCreated {
			t.Fatalf("expired line = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage5-assign-dup-sub", "stage5-assign-dup-mgr2", nil, nil); status != http.StatusCreated {
			t.Fatalf("new line after expiry = %d (%s)", status, string(body))
		}
	})

	t.Run("historical_cycle_does_not_block", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)
		if status, body, _ := createLine(fx.tokenA, "stage5-assign-hcyc-a", "stage5-assign-hcyc-b", &from, &to); status != http.StatusCreated {
			t.Fatalf("historical A->B = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage5-assign-hcyc-b", "stage5-assign-hcyc-c", nil, nil); status != http.StatusCreated {
			t.Fatalf("B->C = %d (%s)", status, string(body))
		}
		if status, body, _ := createLine(fx.tokenA, "stage5-assign-hcyc-c", "stage5-assign-hcyc-a", nil, nil); status != http.StatusCreated {
			t.Fatalf("C->A should not be blocked by an expired historical edge: %d (%s)", status, string(body))
		}
	})

	t.Run("organisation_isolation", func(t *testing.T) {
		status, body, number := createLine(fx.tokenA, "stage4-assign-sub3", "stage4-assign-mgr1", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodGet, s4ReportingPath+"/"+number, fx.tokenB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d (%s)", status, http.StatusNotFound, string(body))
		}
	})

	t.Run("concurrent_duplicate_create", func(t *testing.T) {
		statuses := make(chan int, 2)
		var wg sync.WaitGroup
		managers := []string{"stage4-assign-mgr1", "stage4-assign-mgr2"}
		for _, manager := range managers {
			wg.Add(1)
			go func(manager string) {
				defer wg.Done()
				status, _, err := doRequest(http.MethodPost, s4ReportingPath, fx.tokenA, linePayload("stage4-assign-conc", manager, nil, nil))
				if err != nil {
					status = 0
				}
				statuses <- status
			}(manager)
		}
		wg.Wait()
		close(statuses)

		created, conflict := 0, 0
		for status := range statuses {
			switch status {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				conflict++
			}
		}
		if created != 1 || conflict != 1 {
			t.Fatalf("concurrent create: %d created, %d conflict; want 1 and 1", created, conflict)
		}
	})

	t.Run("deactivation_retains_history_and_blocks_access", func(t *testing.T) {
		status, body, number := createLine(fx.tokenLife, "stage4-assign-life-sub", "stage4-assign-life-mgr", nil, nil)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}

		if status, body := request(http.MethodPost, "/main/eos/entity/delete?entity_number="+s4OrgLife, fx.tokenLife, nil); status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}

		if status, body := request(http.MethodGet, s4ReportingPath, fx.tokenLife, nil); status != http.StatusConflict {
			t.Fatalf("reporting after deactivation = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}

		var stored reportingModel.ReportingLine
		if err := fx.db.Where("number = ?", number).First(&stored).Error; err != nil {
			t.Fatalf("reporting history was removed: %v", err)
		}
		if stored.Status != "ACTIVE" {
			t.Fatalf("retained reporting status = %q, want ACTIVE", stored.Status)
		}
	})
}
