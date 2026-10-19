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
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	stakeholderServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/services"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	s0c2OrgA = "stage0c2-org-a"
	s0c2OrgB = "stage0c2-org-b"
	s0c2OrgC = "stage0c2-org-c"
	s0c2OrgD = "stage0c2-org-d"
	s0c2OrgE = "stage0c2-org-e"
	s0c2OrgF = "stage0c2-org-f"

	s0c2AdminA = "stage0c2-admin-a"
	s0c2AdminB = "stage0c2-admin-b"
	s0c2AdminC = "stage0c2-admin-c"
	s0c2AdminD = "stage0c2-admin-d"
	s0c2AdminE = "stage0c2-admin-e"
	s0c2AdminF = "stage0c2-admin-f"
	s0c2Both   = "stage0c2-both"
	s0c2Stale  = "stage0c2-stale"
	s0c2OwnerF = "stage0c2-owner-f"

	s0c2RoleBoth   = "stage0c2-role-both"
	s0c2EntityPath = "/main/eos/entity"
)

type s0c2Fixture struct {
	db     *gorm.DB
	tokens map[string]string
}

func setupStage0C2Fixture(t *testing.T) *s0c2Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	t.Setenv("JWT_SECRET", "stage-0c2-test-secret")

	users := []string{
		s0c2AdminA, s0c2AdminB, s0c2AdminC, s0c2AdminD, s0c2AdminE, s0c2AdminF,
		s0c2Both, s0c2Stale, s0c2OwnerF,
	}
	organisations := []string{s0c2OrgA, s0c2OrgB, s0c2OrgC, s0c2OrgD, s0c2OrgE, s0c2OrgF}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&notificationModel.Notification{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionModel.UserRole{})
		db.Unscoped().Where("role_number = ?", s0c2RoleBoth).Delete(&permissionModel.RolePermission{})
		db.Unscoped().Where("number = ?", s0c2RoleBoth).Delete(&permissionModel.Role{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.Licence{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	admins := map[string]string{
		s0c2AdminA: s0c2OrgA,
		s0c2AdminB: s0c2OrgB,
		s0c2AdminC: s0c2OrgC,
		s0c2AdminD: s0c2OrgD,
		s0c2AdminE: s0c2OrgE,
		s0c2AdminF: s0c2OrgF,
	}
	for _, number := range users {
		userType := "User"
		if _, isAdmin := admins[number]; isAdmin {
			userType = "Admin"
		}
		user := authModel.User{
			Number:     number,
			FirstName:  "Stage0C2",
			SecondName: number,
			Gender:     "M",
			Type:       userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, organisation := range organisations {
		entity := eosModel.Entity{Number: organisation, Name: "Stage 0C2 " + organisation, AddedBy: s0c2AdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", organisation, err)
		}
	}

	memberships := []struct{ user, entity string }{
		{s0c2AdminA, s0c2OrgA},
		{s0c2AdminB, s0c2OrgB},
		{s0c2AdminC, s0c2OrgC},
		{s0c2AdminD, s0c2OrgD},
		{s0c2AdminE, s0c2OrgE},
		{s0c2AdminF, s0c2OrgF},
		{s0c2Both, s0c2OrgA},
		{s0c2Both, s0c2OrgB},
		{s0c2Stale, s0c2OrgE},
		{s0c2OwnerF, s0c2OrgF},
	}
	for _, spec := range memberships {
		membership := eosModel.EntityUser{
			EntityNumber:  spec.entity,
			UserNumber:    spec.user,
			AddedByNumber: s0c2AdminA,
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
			t.Fatalf("seeding membership %s/%s: %v", spec.user, spec.entity, err)
		}
	}

	activeRows := []struct{ user, entity string }{
		{s0c2AdminA, s0c2OrgA},
		{s0c2AdminB, s0c2OrgB},
		{s0c2AdminC, s0c2OrgC},
		{s0c2AdminD, s0c2OrgD},
		{s0c2AdminE, s0c2OrgE},
		{s0c2AdminF, s0c2OrgF},
		{s0c2Both, s0c2OrgA},
		{s0c2Stale, s0c2OrgE},
	}
	for _, spec := range activeRows {
		active := eosModel.ActiveUserEntity{UserNumber: spec.user, EntityNumber: spec.entity}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.user, err)
		}
	}

	now := time.Now()
	for _, organisation := range organisations {
		activeOrg := organisation
		licence := licensingModel.Licence{
			Number:             "stage0c2-licence-" + organisation,
			LicenceKey:         "stage0c2-licence-key-" + organisation,
			OrganisationNumber: organisation,
			ValidFrom:          now.Add(-time.Hour),
			ValidTo:            now.Add(24 * time.Hour),
			GraceDays:          licensingModel.DefaultGraceDays,
			MaxInstallations:   1,
			LicensedAddons:     `["education"]`,
			Status:             licensingModel.StatusActive,
			ActiveOrgNumber:    &activeOrg,
			IssuedBy:           s0c2AdminA,
		}
		if err := db.Create(&licence).Error; err != nil {
			t.Fatalf("seeding licence for %s: %v", organisation, err)
		}
	}

	role := permissionModel.Role{
		Number:       s0c2RoleBoth,
		Name:         s0c2RoleBoth,
		EntityNumber: s0c2OrgA,
		Status:       "ACTIVE",
		UserNumber:   s0c2AdminA,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding role: %v", err)
	}
	for _, permission := range []string{"can_update_entity", "can_delete_entity"} {
		rolePermission := permissionModel.RolePermission{
			RoleNumber:    s0c2RoleBoth,
			Permission:    permission,
			AddedByNumber: s0c2AdminA,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding role permission %s: %v", permission, err)
		}
	}
	userRole := permissionModel.UserRole{
		UserNumber: s0c2Both,
		RoleNumber: s0c2RoleBoth,
		ExpiryDate: now.Add(24 * time.Hour),
		AddedBy:    s0c2AdminA,
	}
	if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
		t.Fatalf("seeding user role: %v", err)
	}

	owner := stakeholderModel.Stakeholder{
		Number:             "stage0c2-stakeholder-f",
		OrganisationNumber: s0c2OrgF,
		UserNumber:         s0c2OwnerF,
		StakeholderType:    stakeholderModel.StakeholderTypeOwner,
		Status:             "ACTIVE",
		AddedByNumber:      s0c2AdminA,
	}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("seeding owner: %v", err)
	}

	tokens := make(map[string]string, len(users))
	for _, number := range users {
		value, err := authServices.JWTAuthService().GenerateToken(number)
		if err != nil {
			t.Fatalf("generating token for %s: %v", number, err)
		}
		tokens[number] = value
	}

	return &s0c2Fixture{db: db, tokens: tokens}
}

func TestEntityOperationAuthorisationStage0C2(t *testing.T) {
	fx := setupStage0C2Fixture(t)

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

	must := func(method, path, token string, payload any) (int, []byte) {
		t.Helper()
		status, body, err := doRequest(method, path, token, payload)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		return status, body
	}

	updatePayload := func(entityNumber, name string) map[string]any {
		return map[string]any{
			"entity_number": entityNumber,
			"name":          name,
			"phone":         "+256700000000",
			"email":         "entity@example.com",
			"address":       "Kampala",
		}
	}

	t.Run("A_cross_organisation_update_is_forbidden", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/update", fx.tokens[s0c2Both], updatePayload(s0c2OrgB, "Hijacked B"))
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusForbidden, string(body))
		}
		var stored eosModel.Entity
		if err := fx.db.Where("number = ?", s0c2OrgB).First(&stored).Error; err != nil {
			t.Fatalf("loading organisation B: %v", err)
		}
		if stored.Name == "Hijacked B" {
			t.Fatalf("organisation B was modified")
		}
	})

	t.Run("B_cross_organisation_delete_is_forbidden", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/delete?entity_number="+s0c2OrgB, fx.tokens[s0c2Both], nil)
		if status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusForbidden, string(body))
		}
		var stored eosModel.Entity
		if err := fx.db.Where("number = ?", s0c2OrgB).First(&stored).Error; err != nil {
			t.Fatalf("organisation B must remain live: %v", err)
		}
	})

	t.Run("C_authorised_update_works", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/update", fx.tokens[s0c2AdminC], updatePayload(s0c2OrgC, "Updated C"))
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		var stored eosModel.Entity
		if err := fx.db.Where("number = ?", s0c2OrgC).First(&stored).Error; err != nil {
			t.Fatalf("loading organisation C: %v", err)
		}
		if stored.Name != "Updated C" {
			t.Fatalf("organisation C name = %q, want Updated C", stored.Name)
		}
	})

	t.Run("D_authorised_deactivation_preserves_membership_history", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/delete?entity_number="+s0c2OrgD, fx.tokens[s0c2AdminD], nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}

		var memberships int64
		if err := fx.db.Unscoped().Model(&eosModel.EntityUser{}).Where("entity_number = ?", s0c2OrgD).Count(&memberships).Error; err != nil {
			t.Fatalf("counting memberships: %v", err)
		}
		if memberships == 0 {
			t.Fatalf("membership history was deleted")
		}

		var active int64
		if err := fx.db.Unscoped().Model(&eosModel.ActiveUserEntity{}).Where("entity_number = ?", s0c2OrgD).Count(&active).Error; err != nil {
			t.Fatalf("counting active rows: %v", err)
		}
		if active != 0 {
			t.Fatalf("active selections remain for deactivated organisation: %d", active)
		}

		var entity eosModel.Entity
		if err := fx.db.Unscoped().Where("number = ?", s0c2OrgD).First(&entity).Error; err != nil {
			t.Fatalf("loading deactivated entity: %v", err)
		}
		if !entity.DeletedAt.Valid {
			t.Fatalf("entity is not soft-deleted")
		}
	})

	t.Run("F_deactivated_organisation_cannot_be_selected", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/current?entity_number="+s0c2OrgD, fx.tokens[s0c2AdminD], nil)
		if status != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}
	})

	t.Run("G_stale_active_row_does_not_grant_access", func(t *testing.T) {
		if err := fx.db.Unscoped().
			Where("user_number = ? AND entity_number = ?", s0c2Stale, s0c2OrgE).
			Delete(&eosModel.EntityUser{}).Error; err != nil {
			t.Fatalf("removing stale membership: %v", err)
		}

		status, body := must(http.MethodGet, "/main/ems/organisation-units", fx.tokens[s0c2Stale], nil)
		if status != http.StatusConflict {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("J_last_owner_survives_organisation_deactivation", func(t *testing.T) {
		status, body := must(http.MethodPost, s0c2EntityPath+"/delete?entity_number="+s0c2OrgF, fx.tokens[s0c2AdminF], nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}

		var owner stakeholderModel.Stakeholder
		if err := fx.db.Where("organisation_number = ?", s0c2OrgF).First(&owner).Error; err != nil {
			t.Fatalf("owner record was removed: %v", err)
		}
		if owner.Status != "ACTIVE" {
			t.Fatalf("owner status = %q, want ACTIVE", owner.Status)
		}

		if _, err := stakeholderServices.SetStakeholderStatus(s0c2OrgF, owner.Number, "INACTIVE"); err == nil {
			t.Fatalf("final owner was deactivated after organisation lifecycle change")
		}
	})
}

func TestSetupStatusIntegrityStage0C2(t *testing.T) {
	db := testutil.RequireTestDB(t)

	const (
		orphanOwner = "stage0c2-orphan-owner"
		orphanOrg   = "stage0c2-orphan-org"
		liveOwner   = "stage0c2-live-owner"
		liveOrg     = "stage0c2-setup-org"
		setupUser   = "stage0c2-setup-user"
	)

	cleanup := func() {
		db.Unscoped().Where("number IN ?", []string{orphanOwner, liveOwner}).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("entity_number = ?", liveOrg).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number = ?", setupUser).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", []string{liveOrg, orphanOrg}).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number = ?", setupUser).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	// The setup endpoints report installation-wide state. Refuse to run the
	// negative assertions if the shared test database already contains live
	// owners from another fixture.
	var externalOwners int64
	if err := db.Table("ems_stakeholders AS s").
		Joins("JOIN entities AS e ON e.number = s.organisation_number AND e.deleted_at IS NULL").
		Where("s.deleted_at IS NULL").
		Where("s.stakeholder_type = ? AND s.status = ?", stakeholderModel.StakeholderTypeOwner, "ACTIVE").
		Where("s.organisation_number NOT IN ?", []string{orphanOrg, liveOrg}).
		Count(&externalOwners).Error; err != nil {
		t.Fatalf("checking external owners: %v", err)
	}
	if externalOwners > 0 {
		t.Skip("test database already contains live owners from another fixture")
	}

	user := authModel.User{
		Number:     setupUser,
		FirstName:  "Setup",
		SecondName: "User",
		Gender:     "M",
		Type:       "Admin",
		Status:     "ACTIVE",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seeding setup user: %v", err)
	}
	owner := stakeholderModel.Stakeholder{
		Number:             orphanOwner,
		OrganisationNumber: orphanOrg,
		UserNumber:         setupUser,
		StakeholderType:    stakeholderModel.StakeholderTypeOwner,
		Status:             "ACTIVE",
		AddedByNumber:      setupUser,
	}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("seeding orphan owner: %v", err)
	}

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	get := func(path string) (int, []byte) {
		t.Helper()
		resp, err := client.Get(server.URL + path)
		if err != nil {
			t.Fatalf("GET %s failed: %v", path, err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	statusState := func() string {
		t.Helper()
		status, body := get("/setup/status")
		if status != http.StatusOK {
			t.Fatalf("status endpoint = %d (%s)", status, string(body))
		}
		var parsed struct {
			State string `json:"state"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding status: %v", err)
		}
		return parsed.State
	}

	progress := func() map[string]any {
		t.Helper()
		status, body := get("/setup/progress")
		if status != http.StatusOK {
			t.Fatalf("progress endpoint = %d (%s)", status, string(body))
		}
		var parsed struct {
			Steps map[string]any `json:"steps"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding progress: %v", err)
		}
		return parsed.Steps
	}

	if state := statusState(); state != "uninitialized" {
		t.Fatalf("orphan owner established installation state %q, want uninitialized", state)
	}
	if steps := progress(); steps["owner_exists"] != false {
		t.Fatalf("orphan owner reported owner_exists=true: %v", steps)
	}

	liveEntity := eosModel.Entity{Number: liveOrg, Name: "Live Setup Organisation", AddedBy: setupUser}
	if err := db.Create(&liveEntity).Error; err != nil {
		t.Fatalf("seeding live entity: %v", err)
	}
	live := stakeholderModel.Stakeholder{
		Number:             liveOwner,
		OrganisationNumber: liveOrg,
		UserNumber:         setupUser,
		StakeholderType:    stakeholderModel.StakeholderTypeOwner,
		Status:             "ACTIVE",
		AddedByNumber:      setupUser,
	}
	if err := db.Create(&live).Error; err != nil {
		t.Fatalf("seeding live owner: %v", err)
	}

	if state := statusState(); state != "initialized" {
		t.Fatalf("live owner state = %q, want initialized", state)
	}
	if steps := progress(); steps["owner_exists"] != true || steps["organisation_exists"] != true {
		t.Fatalf("live owner progress = %v", steps)
	}

	// Partially completed: the live owner is removed, the orphan remains.
	if err := db.Unscoped().Where("number = ?", liveOwner).Delete(&stakeholderModel.Stakeholder{}).Error; err != nil {
		t.Fatalf("removing live owner: %v", err)
	}
	if state := statusState(); state != "uninitialized" {
		t.Fatalf("partial setup state = %q, want uninitialized", state)
	}
	steps := progress()
	if steps["users_exist"] != true || steps["admin_active"] != true || steps["owner_exists"] != false {
		t.Fatalf("partial setup progress = %v", steps)
	}
}

func TestOrganisationIntegrityIndexesStage0C2(t *testing.T) {
	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running AutoMigrate on the dedicated test database: %v", err)
	}

	for _, spec := range []struct{ table, index string }{
		{"ems_organisation_units", "idx_ems_organisation_units_org_status"},
		{"ems_stakeholders", "idx_ems_stakeholders_org_type_status"},
		{"ems_user_assignments", "idx_ems_user_assignments_org_status"},
		{"ems_decision_authorities", "idx_ems_decision_authorities_org_status"},
	} {
		var count int64
		if err := db.Raw(
			"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			spec.table, spec.index,
		).Scan(&count).Error; err != nil {
			t.Fatalf("checking index %s: %v", spec.index, err)
		}
		if count == 0 {
			t.Fatalf("index %s on %s was not created", spec.index, spec.table)
		}
	}
}
