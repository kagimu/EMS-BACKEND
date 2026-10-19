package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	stage0aOrgA = "stage0a-org-a"
	stage0aOrgB = "stage0a-org-b"

	stage0aAdminA     = "stage0a-admin-a"
	stage0aAdminB     = "stage0a-admin-b"
	stage0aMemberA    = "stage0a-member-a"
	stage0aRoleAdminA = "stage0a-role-admin-a"
	stage0aDelegatorA = "stage0a-delegator-a"
	stage0aTargetA    = "stage0a-target-a"
	stage0aOutsiderB  = "stage0a-outsider-b"
	stage0aNewcomer   = "stage0a-newcomer"

	stage0aRoleAdminANumber     = "stage0a-role-admin-a"
	stage0aRoleDelegatorNumber  = "stage0a-role-delegator-a"
	stage0aRoleOrgBNumber       = "stage0a-role-org-b"
	stage0aMembershipPath       = "/main/eos/add-user-to-entity"
	stage0aRolesPath            = "/main/settings/roles"
	stage0aAssignmentPermission = "update_permission"
)

type stage0aFixture struct {
	db              *gorm.DB
	tokenAdminA     string
	tokenAdminB     string
	tokenMemberA    string
	tokenRoleAdminA string
	tokenDelegatorA string
}

func setupStage0AFixture(t *testing.T) *stage0aFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "stage-0a-test-secret")

	users := []string{
		stage0aAdminA, stage0aAdminB, stage0aMemberA, stage0aRoleAdminA,
		stage0aDelegatorA, stage0aTargetA, stage0aOutsiderB, stage0aNewcomer,
	}
	organisations := []string{stage0aOrgA, stage0aOrgB}

	cleanup := func() {
		// Remove every role attached to the fixture organisations (including
		// roles created by the test itself) and their dependants before the
		// entity rows, so foreign keys never block cleanup.
		var roleNumbers []string
		fxRoles := db.Unscoped().Model(&permissionModel.Role{}).Where("entity_number IN ?", organisations).Pluck("number", &roleNumbers)
		if fxRoles.Error != nil {
			t.Fatalf("collecting fixture roles: %v", fxRoles.Error)
		}
		if len(roleNumbers) > 0 {
			db.Unscoped().Where("role_number IN ?", roleNumbers).Delete(&permissionModel.UserRole{})
			db.Unscoped().Where("role_number IN ?", roleNumbers).Delete(&permissionModel.RolePermission{})
			db.Unscoped().Where("number IN ?", roleNumbers).Delete(&permissionModel.Role{})
		}

		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionModel.UserRole{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.Licence{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	userTypes := map[string]string{
		stage0aAdminA:     "Admin",
		stage0aAdminB:     "Admin",
		stage0aMemberA:    "User",
		stage0aRoleAdminA: "User",
		stage0aDelegatorA: "User",
		stage0aTargetA:    "User",
		stage0aOutsiderB:  "User",
		stage0aNewcomer:   "User",
	}
	for _, number := range users {
		user := authModel.User{
			Number:     number,
			FirstName:  "Stage0A",
			SecondName: number,
			Gender:     "M",
			Type:       userTypes[number],
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{stage0aOrgA, "Stage 0A Organisation A"},
		{stage0aOrgB, "Stage 0A Organisation B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: stage0aAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	memberships := []struct{ user, entity string }{
		{stage0aAdminA, stage0aOrgA},
		{stage0aAdminB, stage0aOrgB},
		{stage0aMemberA, stage0aOrgA},
		{stage0aRoleAdminA, stage0aOrgA},
		{stage0aDelegatorA, stage0aOrgA},
		{stage0aTargetA, stage0aOrgA},
		{stage0aOutsiderB, stage0aOrgB},
	}
	for _, spec := range memberships {
		membership := eosModel.EntityUser{
			EntityNumber:  spec.entity,
			UserNumber:    spec.user,
			AddedByNumber: stage0aAdminA,
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership %s/%s: %v", spec.user, spec.entity, err)
		}

		active := eosModel.ActiveUserEntity{UserNumber: spec.user, EntityNumber: spec.entity}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.user, err)
		}
	}

	now := time.Now()
	for _, organisation := range organisations {
		activeOrg := organisation
		licence := licensingModel.Licence{
			Number:             "stage0a-licence-" + organisation,
			LicenceKey:         "stage0a-licence-key-" + organisation,
			OrganisationNumber: organisation,
			ValidFrom:          now.Add(-time.Hour),
			ValidTo:            now.Add(24 * time.Hour),
			GraceDays:          licensingModel.DefaultGraceDays,
			MaxInstallations:   1,
			LicensedAddons:     `["education"]`,
			Status:             licensingModel.StatusActive,
			ActiveOrgNumber:    &activeOrg,
			IssuedBy:           stage0aAdminA,
		}
		if err := db.Create(&licence).Error; err != nil {
			t.Fatalf("seeding licence for %s: %v", organisation, err)
		}
	}

	seedRole := func(number, entity string, permissions []string) {
		role := permissionModel.Role{
			Number:       number,
			Name:         number,
			EntityNumber: entity,
			Description:  "Stage 0A fixture role",
			Status:       "ACTIVE",
			UserNumber:   stage0aAdminA,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
		for _, permission := range permissions {
			rolePermission := permissionModel.RolePermission{
				RoleNumber:    number,
				Permission:    permission,
				AddedByNumber: stage0aAdminA,
			}
			if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
				t.Fatalf("seeding role permission %s/%s: %v", number, permission, err)
			}
		}
	}

	seedRole(stage0aRoleAdminANumber, stage0aOrgA, []string{stage0aAssignmentPermission, "view_roles"})
	seedRole(stage0aRoleDelegatorNumber, stage0aOrgA, []string{stage0aAssignmentPermission, "view_roles", "view_students"})
	seedRole(stage0aRoleOrgBNumber, stage0aOrgB, []string{stage0aAssignmentPermission, "view_roles"})

	expiry := now.Add(24 * time.Hour)
	for _, spec := range []struct{ user, role string }{
		{stage0aRoleAdminA, stage0aRoleAdminANumber},
		{stage0aDelegatorA, stage0aRoleDelegatorNumber},
	} {
		userRole := permissionModel.UserRole{
			UserNumber: spec.user,
			RoleNumber: spec.role,
			ExpiryDate: expiry,
			AddedBy:    stage0aAdminA,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding user role %s/%s: %v", spec.user, spec.role, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &stage0aFixture{
		db:              db,
		tokenAdminA:     token(stage0aAdminA),
		tokenAdminB:     token(stage0aAdminB),
		tokenMemberA:    token(stage0aMemberA),
		tokenRoleAdminA: token(stage0aRoleAdminA),
		tokenDelegatorA: token(stage0aDelegatorA),
	}
}

func TestOrganisationSecurityStage0A(t *testing.T) {
	fx := setupStage0AFixture(t)

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

	membershipPayload := func(userNumber, entityNumber string) map[string]string {
		return map[string]string{"user_number": userNumber, "entity_number": entityNumber}
	}

	rolePayload := func(name, entityNumber string, permissions []string) map[string]any {
		return map[string]any{
			"name":          name,
			"description":   "Stage 0A test role",
			"entity_number": entityNumber,
			"permissions":   permissions,
		}
	}

	assignPayload := func(userNumber, roleNumber string) map[string]any {
		return map[string]any{
			"user_number": userNumber,
			"roles": []map[string]string{{
				"role_number": roleNumber,
				"expiry_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			}},
		}
	}

	membershipCount := func(userNumber, entityNumber string) int64 {
		var count int64
		if err := fx.db.Model(&eosModel.EntityUser{}).
			Where("user_number = ? AND entity_number = ?", userNumber, entityNumber).
			Count(&count).Error; err != nil {
			t.Fatalf("counting membership: %v", err)
		}
		return count
	}

	userRoleCount := func(userNumber, roleNumber string) int64 {
		var count int64
		if err := fx.db.Model(&permissionModel.UserRole{}).
			Where("user_number = ? AND role_number = ?", userNumber, roleNumber).
			Count(&count).Error; err != nil {
			t.Fatalf("counting user role: %v", err)
		}
		return count
	}

	t.Run("membership_outsider_cannot_add_self", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminB, membershipPayload(stage0aOutsiderB, stage0aOrgA))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
		if memberships := membershipCount(stage0aOutsiderB, stage0aOrgA); memberships != 0 {
			t.Fatalf("outsider gained membership: %d rows", memberships)
		}
	})

	t.Run("membership_admin_of_other_organisation_cannot_add", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminB, membershipPayload(stage0aNewcomer, stage0aOrgA))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("membership_member_without_permission_cannot_add", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenMemberA, membershipPayload(stage0aNewcomer, stage0aOrgA))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("membership_role_admin_without_entity_permission_cannot_add", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenRoleAdminA, membershipPayload(stage0aNewcomer, stage0aOrgA))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("membership_authorised_admin_adds_user", func(t *testing.T) {
		status, body := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminA, membershipPayload(stage0aNewcomer, stage0aOrgA))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if memberships := membershipCount(stage0aNewcomer, stage0aOrgA); memberships != 1 {
			t.Fatalf("membership rows = %d, want 1", memberships)
		}

		var membership eosModel.EntityUser
		if err := fx.db.Where("user_number = ? AND entity_number = ?", stage0aNewcomer, stage0aOrgA).First(&membership).Error; err != nil {
			t.Fatalf("loading membership: %v", err)
		}
		if membership.AddedByNumber != stage0aAdminA {
			t.Fatalf("added_by = %q, want %q", membership.AddedByNumber, stage0aAdminA)
		}
	})

	t.Run("membership_duplicate_is_idempotent", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminA, membershipPayload(stage0aNewcomer, stage0aOrgA))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
		if memberships := membershipCount(stage0aNewcomer, stage0aOrgA); memberships != 1 {
			t.Fatalf("membership rows = %d, want 1", memberships)
		}
	})

	t.Run("membership_invalid_entity_rejected", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminA, membershipPayload(stage0aNewcomer, "stage0a-missing-org"))
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("membership_invalid_user_rejected", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aMembershipPath, fx.tokenAdminA, membershipPayload("stage0a-missing-user", stage0aOrgA))
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("role_create_requires_permission", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath, fx.tokenMemberA, rolePayload("stage0a-denied", stage0aOrgA, []string{"view_students"}))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_create_for_other_organisation_denied", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath, fx.tokenAdminA, rolePayload("stage0a-cross-org", stage0aOrgB, []string{"view_students"}))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_create_unknown_permission_denied", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath, fx.tokenAdminA, rolePayload("stage0a-unknown", stage0aOrgA, []string{"totally_unknown_permission"}))
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("role_create_cannot_delegate_unheld_permission", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath, fx.tokenRoleAdminA, rolePayload("stage0a-escalation", stage0aOrgA, []string{"view_students"}))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	var adminCreatedRole string
	t.Run("role_create_authorised_by_admin", func(t *testing.T) {
		status, body := request(http.MethodPost, stage0aRolesPath, fx.tokenAdminA, rolePayload("stage0a-admin-created", stage0aOrgA, []string{"view_students"}))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}

		var role permissionModel.Role
		if err := fx.db.Where("name = ?", "stage0a-admin-created").First(&role).Error; err != nil {
			t.Fatalf("loading created role: %v", err)
		}
		adminCreatedRole = role.Number
		if role.EntityNumber != stage0aOrgA {
			t.Fatalf("role entity = %q, want %q", role.EntityNumber, stage0aOrgA)
		}
	})

	t.Run("role_create_authorised_by_delegator", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath, fx.tokenDelegatorA, rolePayload("stage0a-delegator-created", stage0aOrgA, []string{"view_students"}))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("role_list_requires_permission", func(t *testing.T) {
		if status, _ := request(http.MethodGet, stage0aRolesPath, fx.tokenMemberA, nil); status != http.StatusForbidden {
			t.Fatalf("member status = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, stage0aRolesPath, fx.tokenRoleAdminA, nil); status != http.StatusOK {
			t.Fatalf("role admin status = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("role_details_cross_organisation_returns_404", func(t *testing.T) {
		status, _ := request(http.MethodGet, stage0aRolesPath+"/"+stage0aRoleOrgBNumber, fx.tokenAdminA, nil)
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("role_assign_requires_permission", func(t *testing.T) {
		status, _ := request(http.MethodPut, stage0aRolesPath+"/assign", fx.tokenMemberA, assignPayload(stage0aTargetA, stage0aRoleDelegatorNumber))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_assign_cross_organisation_returns_404", func(t *testing.T) {
		status, _ := request(http.MethodPut, stage0aRolesPath+"/assign", fx.tokenAdminA, assignPayload(stage0aTargetA, stage0aRoleOrgBNumber))
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("role_assign_target_must_be_member", func(t *testing.T) {
		status, _ := request(http.MethodPut, stage0aRolesPath+"/assign", fx.tokenAdminA, assignPayload(stage0aOutsiderB, stage0aRoleDelegatorNumber))
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("role_assign_cannot_delegate_unheld_permission", func(t *testing.T) {
		status, _ := request(http.MethodPut, stage0aRolesPath+"/assign", fx.tokenRoleAdminA, assignPayload(stage0aTargetA, adminCreatedRole))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_assign_authorised", func(t *testing.T) {
		status, body := request(http.MethodPut, stage0aRolesPath+"/assign", fx.tokenDelegatorA, assignPayload(stage0aTargetA, stage0aRoleDelegatorNumber))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if count := userRoleCount(stage0aTargetA, stage0aRoleDelegatorNumber); count != 1 {
			t.Fatalf("assignment rows = %d, want 1", count)
		}
	})

	t.Run("role_remove_requires_permission", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/remove", fx.tokenMemberA, map[string]string{
			"user_number": stage0aTargetA, "role_number": stage0aRoleDelegatorNumber,
		})
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_remove_cross_organisation_returns_404", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/remove", fx.tokenAdminA, map[string]string{
			"user_number": stage0aTargetA, "role_number": stage0aRoleOrgBNumber,
		})
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("role_remove_authorised", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/remove", fx.tokenDelegatorA, map[string]string{
			"user_number": stage0aTargetA, "role_number": stage0aRoleDelegatorNumber,
		})
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
		if count := userRoleCount(stage0aTargetA, stage0aRoleDelegatorNumber); count != 0 {
			t.Fatalf("assignment rows = %d, want 0", count)
		}
	})

	t.Run("role_update_cross_organisation_returns_404", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/update/"+stage0aRoleOrgBNumber, fx.tokenAdminA, map[string]any{
			"name": "cross-org update", "description": "denied", "permissions": []string{"view_roles"},
		})
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("role_update_unknown_permission_denied", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/update/"+adminCreatedRole, fx.tokenAdminA, map[string]any{
			"name": "stage0a-admin-created", "description": "updated", "permissions": []string{"totally_unknown_permission"},
		})
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("role_update_authorised", func(t *testing.T) {
		status, _ := request(http.MethodPost, stage0aRolesPath+"/update/"+adminCreatedRole, fx.tokenAdminA, map[string]any{
			"name": "stage0a-admin-created", "description": "updated", "permissions": []string{"view_roles"},
		})
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("role_delete_cross_organisation_returns_404", func(t *testing.T) {
		status, _ := request(http.MethodDelete, stage0aRolesPath+"/delete/"+stage0aRoleOrgBNumber, fx.tokenAdminA, nil)
		if status != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("role_delete_requires_permission", func(t *testing.T) {
		status, _ := request(http.MethodDelete, stage0aRolesPath+"/delete/"+adminCreatedRole, fx.tokenMemberA, nil)
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("role_delete_authorised", func(t *testing.T) {
		status, _ := request(http.MethodDelete, stage0aRolesPath+"/delete/"+adminCreatedRole, fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
		var count int64
		if err := fx.db.Model(&permissionModel.Role{}).Where("number = ?", adminCreatedRole).Count(&count).Error; err != nil {
			t.Fatalf("counting deleted role: %v", err)
		}
		if count != 0 {
			t.Fatalf("role rows = %d, want 0", count)
		}
	})
}
