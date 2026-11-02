package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
)

func stage1Request(t *testing.T, client *http.Client, serverURL, method, path, token string, body any) (int, []byte) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encoding request body: %v", err)
		}
		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, serverURL+path, reader)
	if err != nil {
		t.Fatalf("building request %s %s: %v", method, path, err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request %s %s failed: %v", method, path, err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, responseBody
}

func cleanupStage1Role(t *testing.T, fx *stage0bFixture, name string) {
	t.Helper()
	t.Cleanup(func() {
		var role permissionModel.Role
		if err := fx.db.Unscoped().Where("name = ?", name).First(&role).Error; err != nil {
			return
		}
		fx.db.Unscoped().Where("role_number = ?", role.Number).Delete(&permissionModel.RolePermission{})
		fx.db.Unscoped().Where("role_number = ?", role.Number).Delete(&permissionModel.UserRole{})
		fx.db.Unscoped().Where("number = ?", role.Number).Delete(&permissionModel.Role{})
	})
}

func seedStage1Licence(t *testing.T, fx *stage0bFixture, organisation, number string) {
	t.Helper()
	activeOrg := organisation
	licence := licensingModel.Licence{
		Number:             number,
		OrganisationNumber: organisation,
		LicenceKey:         number + "-key",
		ValidFrom:          time.Now().Add(-time.Hour),
		ValidTo:            time.Now().Add(24 * time.Hour),
		GraceDays:          14,
		MaxInstallations:   10,
		LicensedAddons:     `["education"]`,
		Status:             licensingModel.StatusActive,
		ActiveOrgNumber:    &activeOrg,
	}
	if err := fx.db.Create(&licence).Error; err != nil {
		t.Fatalf("seeding licence for %s: %v", organisation, err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("number = ?", number).Delete(&licensingModel.Licence{})
	})
}

func TestStage1UserDirectoryAndRoleDelegation(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage1-licence-a")
	seedStage1Licence(t, fx, s0bOrgB, "stage1-licence-b")
	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}

	t.Run("unauthenticated_user_listing_is_unauthorised", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/users", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := req(http.MethodGet, "/main/users/all", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("all-users status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("organisation_listing_requires_view_users", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/users", fx.tokenInactive, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("organisation_listing_is_scoped_to_active_organisation", func(t *testing.T) {
		status, body := req(http.MethodGet, "/main/users?per_page=100", fx.tokenBoth, nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d (%s)", status, string(body))
		}
		if !strings.Contains(string(body), s0bAdminA) {
			t.Fatalf("organisation A user missing from response: %s", string(body))
		}
		if strings.Contains(string(body), s0bUserOnlyB) {
			t.Fatalf("organisation B user leaked into organisation A response: %s", string(body))
		}

		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", s0bUserBoth).
			Update("entity_number", s0bOrgB).Error; err != nil {
			t.Fatalf("switching active organisation: %v", err)
		}
		if status, _ := req(http.MethodGet, "/main/users", fx.tokenBoth, nil); status != http.StatusForbidden {
			t.Fatalf("org B status = %d, want %d", status, http.StatusForbidden)
		}
		if err := fx.db.Model(&eosModel.ActiveUserEntity{}).
			Where("user_number = ?", s0bUserBoth).
			Update("entity_number", s0bOrgA).Error; err != nil {
			t.Fatalf("restoring active organisation: %v", err)
		}
	})

	t.Run("single_user_lookup_is_organisation_scoped", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/users/"+s0bAdminA, fx.tokenBoth, nil); status != http.StatusOK {
			t.Fatalf("same-organisation user = %d, want %d", status, http.StatusOK)
		}
		if status, _ := req(http.MethodGet, "/main/users/"+s0bUserOnlyB, fx.tokenBoth, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation user = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("all_users_requires_installation_authority", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/users/all", fx.tokenAdminA, nil); status != http.StatusForbidden {
			t.Fatalf("admin without authority = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("client_admin_cannot_grant_control_plane_permissions", func(t *testing.T) {
		payload := map[string]any{
			"name":          "Stage1 Blocked Owner Role",
			"description":   "must not be created",
			"permissions":   []string{"manage_licence_registry"},
			"entity_number": s0bOrgA,
		}
		if status, body := req(http.MethodPost, "/main/settings/roles", fx.tokenAdminA, payload); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusForbidden, string(body))
		}

		payload["permissions"] = []string{"view_all_users"}
		if status, _ := req(http.MethodPost, "/main/settings/roles", fx.tokenAdminA, payload); status != http.StatusForbidden {
			t.Fatalf("view_all_users status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("client_admin_can_still_create_ordinary_roles", func(t *testing.T) {
		name := "Stage1 Allowed Client Role"
		cleanupStage1Role(t, fx, name)
		payload := map[string]any{
			"name":          name,
			"description":   "ordinary client role",
			"permissions":   []string{"view_users"},
			"entity_number": s0bOrgA,
		}
		if status, body := req(http.MethodPost, "/main/settings/roles", fx.tokenAdminA, payload); status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("client_admin_cannot_mutate_control_plane_role", func(t *testing.T) {
		if status, _ := req(http.MethodPost, "/main/settings/roles/update/"+s0bRoleLegacyOwner, fx.tokenAdminA, map[string]any{
			"name":        "Stage1 Legacy Owner",
			"description": "attempted update",
			"permissions": []string{"view_users"},
		}); status != http.StatusForbidden {
			t.Fatalf("update status = %d, want %d", status, http.StatusForbidden)
		}

		if status, _ := req(http.MethodPut, "/main/settings/roles/assign", fx.tokenAdminA, map[string]any{
			"user_number": s0bUserBoth,
			"roles":       []map[string]string{{"role_number": s0bRoleLegacyOwner, "expiry_date": "2027-01-01"}},
		}); status != http.StatusForbidden {
			t.Fatalf("assign status = %d, want %d", status, http.StatusForbidden)
		}

		if status, _ := req(http.MethodPost, "/main/settings/roles/remove", fx.tokenAdminA, map[string]any{
			"user_number": s0bUserBoth,
			"role_number": s0bRoleLegacyOwner,
		}); status != http.StatusForbidden {
			t.Fatalf("remove status = %d, want %d", status, http.StatusForbidden)
		}

		if status, _ := req(http.MethodDelete, "/main/settings/roles/delete/"+s0bRoleLegacyOwner, fx.tokenAdminA, nil); status != http.StatusForbidden {
			t.Fatalf("delete status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("create_user_cannot_assign_control_plane_role", func(t *testing.T) {
		payload := map[string]any{
			"first_name":    "Stage1",
			"second_name":   "Blocked",
			"last_name":     "User",
			"gender":        "M",
			"email":         "stage1-blocked@example.com",
			"phone":         "0771234567",
			"entity_number": s0bOrgA,
			"roles":         []map[string]string{{"role_number": s0bRoleLegacyOwner, "expiry_date": "2027-01-01"}},
		}
		if status, body := req(http.MethodPost, "/main/users", fx.tokenAdminA, payload); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusForbidden, string(body))
		}
	})
}

func TestStage1SoftwareOwnerAuthority(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage1-owner-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	req := func(method, path string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, fx.tokenAdminA, body)
	}

	t.Run("owner_can_read_installation_user_directory", func(t *testing.T) {
		status, body := req(http.MethodGet, "/main/users/all", nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d (%s)", status, string(body))
		}
		if !strings.Contains(string(body), s0bUserOnlyB) {
			t.Fatalf("installation directory must include all organisations: %s", string(body))
		}
	})

	t.Run("owner_can_create_control_plane_role", func(t *testing.T) {
		name := "Stage1 Owner Role"
		cleanupStage1Role(t, fx, name)
		status, body := req(http.MethodPost, "/main/settings/roles", map[string]any{
			"name":          name,
			"description":   "owner role",
			"permissions":   []string{"manage_licence_registry"},
			"entity_number": s0bOrgA,
		})
		if status != http.StatusOK {
			t.Fatalf("status = %d (%s)", status, string(body))
		}
	})
}
