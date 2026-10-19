package integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	s0bOrgA = "stage0b-org-a"
	s0bOrgB = "stage0b-org-b"

	s0bAdminA     = "stage0b-admin-a"
	s0bAdminNoOrg = "stage0b-admin-no-org"
	s0bUserBoth   = "stage0b-user-both"
	s0bUserExpiry = "stage0b-user-expired"
	s0bUserInact  = "stage0b-user-inactive"
	s0bUserCross  = "stage0b-user-cross"
	s0bUserPermB  = "stage0b-user-perm-b"
	s0bUserNoOrg  = "stage0b-user-no-org"
	s0bRegistry   = "stage0b-registry"
	s0bOwnerUser  = "stage0b-owner-user"

	s0bRoleAValid    = "stage0b-role-a-valid"
	s0bRoleAExpired  = "stage0b-role-a-expired"
	s0bRoleAInactive = "stage0b-role-a-inactive"
	s0bRoleBValid    = "stage0b-role-b-valid"
	s0bRoleRegistry  = "stage0b-role-registry"

	s0bUnitB    = "stage0b-unit-b"
	s0bUnits    = "/main/ems/organisation-units"
	s0bAuditReg = "/main/ems/audit-registry"
)

type stage0bFixture struct {
	db            *gorm.DB
	tokenAdminA   string
	tokenAdminNo  string
	tokenBoth     string
	tokenExpired  string
	tokenInactive string
	tokenCross    string
	tokenPermB    string
	tokenNoOrg    string
	tokenRegistry string
	tokenOwner    string
}

func setupStage0BFixture(t *testing.T) *stage0bFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "stage-0b-test-secret")

	users := []string{
		s0bAdminA, s0bAdminNoOrg, s0bUserBoth, s0bUserExpiry, s0bUserInact,
		s0bUserCross, s0bUserPermB, s0bUserNoOrg, s0bRegistry, s0bOwnerUser,
	}
	organisations := []string{s0bOrgA, s0bOrgB}
	roles := []string{s0bRoleAValid, s0bRoleAExpired, s0bRoleAInactive, s0bRoleBValid, s0bRoleRegistry}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionModel.UserRole{})
		db.Unscoped().Where("role_number IN ?", roles).Delete(&permissionModel.RolePermission{})
		db.Unscoped().Where("number IN ?", roles).Delete(&permissionModel.Role{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	userTypes := map[string]string{
		s0bAdminA:     "Admin",
		s0bAdminNoOrg: "Admin",
	}
	for _, number := range users {
		userType := "User"
		if value, ok := userTypes[number]; ok {
			userType = value
		}
		user := authModel.User{
			Number:     number,
			FirstName:  "Stage0B",
			SecondName: number,
			Gender:     "M",
			Type:       userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{s0bOrgA, "Stage 0B Organisation A"},
		{s0bOrgB, "Stage 0B Organisation B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: s0bAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	memberships := []struct{ user, entity string }{
		{s0bAdminA, s0bOrgA},
		{s0bAdminNoOrg, s0bOrgA},
		{s0bUserBoth, s0bOrgA},
		{s0bUserBoth, s0bOrgB},
		{s0bUserExpiry, s0bOrgA},
		{s0bUserInact, s0bOrgA},
		{s0bUserCross, s0bOrgA},
		{s0bUserCross, s0bOrgB},
		{s0bUserPermB, s0bOrgA},
		{s0bUserPermB, s0bOrgB},
		{s0bUserNoOrg, s0bOrgA},
		{s0bRegistry, s0bOrgA},
		{s0bOwnerUser, s0bOrgA},
	}
	for _, spec := range memberships {
		membership := eosModel.EntityUser{
			EntityNumber:  spec.entity,
			UserNumber:    spec.user,
			AddedByNumber: s0bAdminA,
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership %s/%s: %v", spec.user, spec.entity, err)
		}
	}

	for _, spec := range []struct{ user, entity string }{
		{s0bAdminA, s0bOrgA},
		{s0bUserBoth, s0bOrgA},
		{s0bUserExpiry, s0bOrgA},
		{s0bUserInact, s0bOrgA},
		{s0bUserCross, s0bOrgB},
		{s0bUserPermB, s0bOrgA},
		{s0bOwnerUser, s0bOrgA},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.user, EntityNumber: spec.entity}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.user, err)
		}
	}

	for _, spec := range []struct{ number, user string }{
		{"stage0b-owner-admin", s0bAdminA},
		{"stage0b-owner-user", s0bOwnerUser},
	} {
		stakeholder := stakeholderModel.Stakeholder{
			Number:             spec.number,
			OrganisationNumber: s0bOrgA,
			UserNumber:         spec.user,
			StakeholderType:    stakeholderModel.StakeholderTypeOwner,
			Status:             "ACTIVE",
			AddedByNumber:      s0bAdminA,
		}
		if err := db.Create(&stakeholder).Error; err != nil {
			t.Fatalf("seeding owner stakeholder: %v", err)
		}
	}

	unit := organisationModel.OrganisationUnit{
		Number:             s0bUnitB,
		OrganisationNumber: s0bOrgB,
		Name:               "Organisation B Branch",
		UnitType:           "branch",
		Status:             "ACTIVE",
		AddedByNumber:      s0bAdminA,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("seeding organisation B unit: %v", err)
	}

	seedRole := func(number, organisation, status string, permissions []string) {
		role := permissionModel.Role{
			Number:       number,
			Name:         number,
			EntityNumber: organisation,
			Status:       status,
			UserNumber:   s0bAdminA,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
		for _, permission := range permissions {
			rolePermission := permissionModel.RolePermission{
				RoleNumber:    number,
				Permission:    permission,
				AddedByNumber: s0bAdminA,
			}
			if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
				t.Fatalf("seeding permission %s: %v", permission, err)
			}
		}
	}

	seedRole(s0bRoleAValid, s0bOrgA, "ACTIVE", []string{"view_ems_organisation"})
	seedRole(s0bRoleAExpired, s0bOrgA, "ACTIVE", []string{"view_ems_organisation"})
	seedRole(s0bRoleAInactive, s0bOrgA, "INACTIVE", []string{"view_ems_organisation"})
	seedRole(s0bRoleBValid, s0bOrgB, "ACTIVE", []string{"view_ems_organisation"})
	seedRole(s0bRoleRegistry, s0bOrgA, "ACTIVE", []string{"manage_ems_audit"})

	now := time.Now()
	assign := func(user, role string, expiry time.Time) {
		userRole := permissionModel.UserRole{
			UserNumber: user,
			RoleNumber: role,
			ExpiryDate: expiry,
			AddedBy:    s0bAdminA,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding user role %s/%s: %v", user, role, err)
		}
	}

	assign(s0bUserBoth, s0bRoleAValid, now.Add(24*time.Hour))
	assign(s0bUserExpiry, s0bRoleAExpired, now.Add(-time.Hour))
	assign(s0bUserInact, s0bRoleAInactive, now.Add(24*time.Hour))
	assign(s0bUserCross, s0bRoleAExpired, now.Add(-time.Hour))
	assign(s0bUserCross, s0bRoleBValid, now.Add(24*time.Hour))
	assign(s0bUserPermB, s0bRoleBValid, now.Add(24*time.Hour))
	assign(s0bUserNoOrg, s0bRoleAValid, now.Add(24*time.Hour))
	assign(s0bRegistry, s0bRoleRegistry, now.Add(24*time.Hour))

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &stage0bFixture{
		db:            db,
		tokenAdminA:   token(s0bAdminA),
		tokenAdminNo:  token(s0bAdminNoOrg),
		tokenBoth:     token(s0bUserBoth),
		tokenExpired:  token(s0bUserExpiry),
		tokenInactive: token(s0bUserInact),
		tokenCross:    token(s0bUserCross),
		tokenPermB:    token(s0bUserPermB),
		tokenNoOrg:    token(s0bUserNoOrg),
		tokenRegistry: token(s0bRegistry),
		tokenOwner:    token(s0bOwnerUser),
	}
}

func TestPermissionResolutionStage0B(t *testing.T) {
	fx := setupStage0BFixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	request := func(method, path, token string) (int, []byte) {
		t.Helper()

		var reader io.Reader
		req, err := http.NewRequest(method, server.URL+path, reader)
		if err != nil {
			t.Fatalf("building request %s %s: %v", method, path, err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	switchOrganisation := func(userNumber, organisation string) {
		t.Helper()
		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", userNumber).
			Update("entity_number", organisation).Error; err != nil {
			t.Fatalf("switching active organisation for %s: %v", userNumber, err)
		}
	}

	t.Run("A_permission_in_A_does_not_authorise_B", func(t *testing.T) {
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenBoth); status != http.StatusOK {
			t.Fatalf("organisation A status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		switchOrganisation(s0bUserBoth, s0bOrgB)
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenBoth); status != http.StatusForbidden {
			t.Fatalf("organisation B status = %d, want %d (%s)", status, http.StatusForbidden, string(body))
		}
	})

	t.Run("B_expired_assignment_does_not_grant", func(t *testing.T) {
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenExpired); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("C_inactive_role_does_not_grant", func(t *testing.T) {
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenInactive); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("D_other_organisation_invalid_role_does_not_affect_valid_org", func(t *testing.T) {
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenCross); status != http.StatusOK {
			t.Fatalf("organisation B status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		switchOrganisation(s0bUserCross, s0bOrgA)
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenCross); status != http.StatusForbidden {
			t.Fatalf("organisation A status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("E_missing_active_organisation_fails_safely", func(t *testing.T) {
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenNoOrg); status != http.StatusConflict {
			t.Fatalf("user status = %d, want %d", status, http.StatusConflict)
		}
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenAdminNo); status != http.StatusConflict {
			t.Fatalf("admin status = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("F_client_supplied_organisation_is_ignored", func(t *testing.T) {
		if status, _ := request(http.MethodGet, s0bUnits+"?entity_number="+s0bOrgB, fx.tokenPermB); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
		switchOrganisation(s0bUserPermB, s0bOrgB)
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenPermB); status != http.StatusOK {
			t.Fatalf("organisation B status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("G_installation_level_authority_still_works", func(t *testing.T) {
		if status, body := request(http.MethodGet, s0bAuditReg, fx.tokenRegistry); status != http.StatusOK {
			t.Fatalf("registry status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if status, _ := request(http.MethodGet, s0bAuditReg, fx.tokenAdminNo); status != http.StatusForbidden {
			t.Fatalf("admin without registry authority = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("H_admin_is_bounded_to_active_organisation", func(t *testing.T) {
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenAdminA); status != http.StatusOK {
			t.Fatalf("admin in organisation A = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if status, _ := request(http.MethodGet, s0bUnits+"/"+s0bUnitB, fx.tokenAdminA); status != http.StatusNotFound {
			t.Fatalf("cross-organisation unit = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenAdminNo); status != http.StatusConflict {
			t.Fatalf("admin without organisation = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("I_owner_retains_intended_access_and_ownership_is_not_permission", func(t *testing.T) {
		if status, body := request(http.MethodGet, s0bUnits, fx.tokenAdminA); status != http.StatusOK {
			t.Fatalf("owner admin status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if status, _ := request(http.MethodGet, s0bUnits, fx.tokenOwner); status != http.StatusForbidden {
			t.Fatalf("non-admin owner without a role = %d, want %d", status, http.StatusForbidden)
		}
	})
}
