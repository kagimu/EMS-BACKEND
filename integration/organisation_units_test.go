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
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	orgUnitUserAdminA   = "ems-orgunit-admin-a"
	orgUnitUserAdminB   = "ems-orgunit-admin-b"
	orgUnitUserNoPerm   = "ems-orgunit-noperm"
	orgUnitUserNoCtx    = "ems-orgunit-nocontext"
	orgUnitOrgA         = "ems-orgunit-org-a"
	orgUnitOrgB         = "ems-orgunit-org-b"
	orgUnitBasePath     = "/main/ems/organisation-units"
	orgUnitEducationURL = "/main/campus/all"
)

type orgUnitFixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAdminB string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupOrgUnitFixture(t *testing.T) *orgUnitFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-b-test-secret")

	users := []string{orgUnitUserAdminA, orgUnitUserAdminB, orgUnitUserNoPerm, orgUnitUserNoCtx}
	organisations := []string{orgUnitOrgA, orgUnitOrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct {
		number   string
		userType string
	}{
		{orgUnitUserAdminA, "Admin"},
		{orgUnitUserAdminB, "Admin"},
		{orgUnitUserNoPerm, "User"},
		{orgUnitUserNoCtx, "Admin"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "OrgUnit",
			SecondName: spec.userType,
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{orgUnitOrgA, "Org Unit Test A"},
		{orgUnitOrgB, "Org Unit Test B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: orgUnitUserAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{orgUnitUserAdminA, orgUnitOrgA},
		{orgUnitUserAdminB, orgUnitOrgB},
		{orgUnitUserNoPerm, orgUnitOrgA},
		// orgUnitUserNoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&eosModel.EntityUser{
			EntityNumber:  spec.entityNumber,
			UserNumber:    spec.userNumber,
			AddedByNumber: spec.userNumber,
		}).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &orgUnitFixture{
		db:          db,
		tokenAdminA: token(orgUnitUserAdminA),
		tokenAdminB: token(orgUnitUserAdminB),
		tokenNoPerm: token(orgUnitUserNoPerm),
		tokenNoCtx:  token(orgUnitUserNoCtx),
	}
}

type orgUnitEnvelope struct {
	Message string `json:"message"`
	Data    struct {
		Number       string `json:"number"`
		Name         string `json:"name"`
		UnitType     string `json:"unit_type"`
		Status       string `json:"status"`
		ParentNumber string `json:"parent_number"`
	} `json:"data"`
}

func decodeUnit(t *testing.T, body []byte) orgUnitEnvelope {
	t.Helper()
	var envelope orgUnitEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decoding response %s: %v", string(body), err)
	}
	return envelope
}

// TestOrganisationUnits covers Foundation B end to end. It shares one fixture so
// the (slow) migration and seeding run once.
func TestOrganisationUnits(t *testing.T) {
	fx := setupOrgUnitFixture(t)

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

	createUnit := func(token, name, unitType, code, parent string) (int, []byte) {
		return request(http.MethodPost, orgUnitBasePath, token, map[string]string{
			"name":          name,
			"unit_type":     unitType,
			"code":          code,
			"parent_number": parent,
		})
	}

	t.Run("unauthenticated_returns_401", func(t *testing.T) {
		if status, _ := request(http.MethodGet, orgUnitBasePath, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated_without_permission_returns_403", func(t *testing.T) {
		if status, _ := request(http.MethodGet, orgUnitBasePath, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("no_active_organisation_returns_409", func(t *testing.T) {
		if status, _ := request(http.MethodGet, orgUnitBasePath, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("status = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("create_retrieve_update_status_delete", func(t *testing.T) {
		status, body := createUnit(fx.tokenAdminA, "Operations", "department", "OPS-1", "")
		if status != http.StatusCreated {
			t.Fatalf("create status = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		created := decodeUnit(t, body)
		if created.Data.Number == "" || created.Data.Status != "ACTIVE" {
			t.Fatalf("unexpected created unit: %s", string(body))
		}

		path := orgUnitBasePath + "/" + created.Data.Number

		if status, _ = request(http.MethodGet, path, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("get status = %d, want %d", status, http.StatusOK)
		}

		status, body = request(http.MethodPut, path, fx.tokenAdminA, map[string]string{
			"name": "Operations Updated", "unit_type": "division", "code": "OPS-1", "parent_number": "",
		})
		if status != http.StatusOK {
			t.Fatalf("update status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if updated := decodeUnit(t, body); updated.Data.Name != "Operations Updated" {
			t.Fatalf("update did not apply: %s", string(body))
		}

		status, body = request(http.MethodPatch, path+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("status update = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if updated := decodeUnit(t, body); updated.Data.Status != "INACTIVE" {
			t.Fatalf("status not applied: %s", string(body))
		}

		if status, _ = request(http.MethodDelete, path, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("delete status = %d, want %d", status, http.StatusOK)
		}
		if status, _ = request(http.MethodGet, path, fx.tokenAdminA, nil); status != http.StatusNotFound {
			t.Fatalf("get after delete = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("hierarchy_rules", func(t *testing.T) {
		status, body := createUnit(fx.tokenAdminA, "Head Office", "head_office", "HO-1", "")
		if status != http.StatusCreated {
			t.Fatalf("parent create = %d (%s)", status, string(body))
		}
		parent := decodeUnit(t, body)
		parentPath := orgUnitBasePath + "/" + parent.Data.Number

		status, body = createUnit(fx.tokenAdminA, "Finance", "department", "FIN-1", parent.Data.Number)
		if status != http.StatusCreated {
			t.Fatalf("child create = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		child := decodeUnit(t, body)
		if child.Data.ParentNumber != parent.Data.Number {
			t.Fatalf("child parent = %q, want %q", child.Data.ParentNumber, parent.Data.Number)
		}

		if status, _ = createUnit(fx.tokenAdminA, "Orphan", "unit", "ORPH-1", "does-not-exist"); status != http.StatusBadRequest {
			t.Fatalf("invalid parent = %d, want %d", status, http.StatusBadRequest)
		}

		status, _ = request(http.MethodPut, orgUnitBasePath+"/"+child.Data.Number, fx.tokenAdminA, map[string]string{
			"name": "Finance", "unit_type": "department", "code": "FIN-1", "parent_number": child.Data.Number,
		})
		if status != http.StatusBadRequest {
			t.Fatalf("self-parent = %d, want %d", status, http.StatusBadRequest)
		}

		status, _ = request(http.MethodPut, parentPath, fx.tokenAdminA, map[string]string{
			"name": "Head Office", "unit_type": "head_office", "code": "HO-1", "parent_number": child.Data.Number,
		})
		if status != http.StatusBadRequest {
			t.Fatalf("circular parent = %d, want %d", status, http.StatusBadRequest)
		}

		status, body = createUnit(fx.tokenAdminA, "Inactive Parent", "division", "INACT-1", "")
		if status != http.StatusCreated {
			t.Fatalf("inactive parent create = %d (%s)", status, string(body))
		}
		inactiveParent := decodeUnit(t, body)
		status, _ = request(http.MethodPatch, orgUnitBasePath+"/"+inactiveParent.Data.Number+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"})
		if status != http.StatusOK {
			t.Fatalf("deactivate parent = %d, want %d", status, http.StatusOK)
		}
		if status, _ = createUnit(fx.tokenAdminA, "Under Inactive", "unit", "UNDER-1", inactiveParent.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("inactive parent child = %d, want %d", status, http.StatusBadRequest)
		}

		if status, _ = request(http.MethodDelete, parentPath, fx.tokenAdminA, nil); status != http.StatusConflict {
			t.Fatalf("delete with children = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("organisation_isolation_and_idor", func(t *testing.T) {
		status, body := createUnit(fx.tokenAdminA, "Org A Only Unit", "department", "ORGA-1", "")
		if status != http.StatusCreated {
			t.Fatalf("org A create = %d (%s)", status, string(body))
		}
		orgAUnit := decodeUnit(t, body)
		path := orgUnitBasePath + "/" + orgAUnit.Data.Number

		_, listA := request(http.MethodGet, orgUnitBasePath, fx.tokenAdminA, nil)
		if !bytes.Contains(listA, []byte(orgAUnit.Data.Number)) {
			t.Fatalf("org A list must contain its own unit: %s", string(listA))
		}

		_, listB := request(http.MethodGet, orgUnitBasePath, fx.tokenAdminB, nil)
		if bytes.Contains(listB, []byte(orgAUnit.Data.Number)) {
			t.Fatalf("org B list must not contain org A unit: %s", string(listB))
		}

		if status, _ = request(http.MethodGet, path, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodPut, path, fx.tokenAdminB, map[string]string{
			"name": "Hijacked", "unit_type": "unit", "code": "ORGA-1", "parent_number": "",
		}); status != http.StatusNotFound {
			t.Fatalf("cross-organisation update = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ = request(http.MethodDelete, path, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation delete = %d, want %d", status, http.StatusNotFound)
		}

		if status, _ = createUnit(fx.tokenAdminB, "Cross Parent", "unit", "CROSS-1", orgAUnit.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("cross-organisation parent = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, orgUnitEducationURL, fx.tokenAdminA, nil)

		// Mutate EMS organisation units heavily.
		status, body := createUnit(fx.tokenAdminA, "Scope Check", "division", "SCOPE-1", "")
		if status != http.StatusCreated {
			t.Fatalf("scope-check create = %d (%s)", status, string(body))
		}
		unit := decodeUnit(t, body)
		request(http.MethodPatch, orgUnitBasePath+"/"+unit.Data.Number+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"})
		request(http.MethodDelete, orgUnitBasePath+"/"+unit.Data.Number, fx.tokenAdminA, nil)

		statusAfter, bodyAfter := request(http.MethodGet, orgUnitEducationURL, fx.tokenAdminA, nil)

		if statusBefore != statusAfter {
			t.Fatalf("Education status changed from %d to %d after EMS unit changes", statusBefore, statusAfter)
		}
		if !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS unit changes:\nbefore: %s\nafter:  %s", string(bodyBefore), string(bodyAfter))
		}
	})
}
