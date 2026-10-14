package integration

import (
	"strings"
	"testing"
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	g4LicOwner   = "ems-g4-lic-owner"
	g4LicControl = "ems-g4-lic-control"
	g4LicDual    = "ems-g4-lic-dual"
	g4LicAdmin   = "ems-g4-lic-admin"
	g4LicOwnerB  = "ems-g4-lic-owner-b"

	g4LicOrg    = "ems-g4-lic-org"
	g4LicOrphan = "ems-g4-lic-orphan"
	g4LicOrgB   = "ems-g4-lic-org-b"

	g4LicRole = "ems-g4-lic-role"

	g4LicenceKey     = "SECRET-KEY-DO-NOT-LEAK"
	g4LicenceNumber  = "g4-licence-1"
	g4LicenceNumber2 = "g4-licence-2"
)

func setupG4LicensingFixture(t *testing.T) *gorm.DB {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}

	users := []string{g4LicOwner, g4LicControl, g4LicDual, g4LicAdmin, g4LicOwnerB}
	organisations := []string{g4LicOrg, g4LicOrphan, g4LicOrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&notificationModel.Notification{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.LicenceEvent{})
		db.Unscoped().Where("role_number = ?", g4LicRole).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("role_number = ?", g4LicRole).Delete(&permissionsModel.RolePermission{})
		db.Unscoped().Where("number = ?", g4LicRole).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, number := range users {
		userType := "User"
		if number == g4LicAdmin {
			// A generic Admin must NOT be a licensing-management recipient.
			userType = "Admin"
		}
		user := authModel.User{Number: number, FirstName: "G4", SecondName: "Lic", Gender: "M", Type: userType, Status: "ACTIVE"}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, number := range organisations {
		entity := eosModel.Entity{Number: number, Name: number, AddedBy: g4LicOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", number, err)
		}
	}

	// Active owner stakeholders: g4LicOwner (owner only) and g4LicDual (owner
	// AND registry authority) for the dedup case. g4LicOrphan has no owner so
	// the fail-closed path can be observed.
	for _, spec := range []struct{ org, user string }{
		{g4LicOrg, g4LicOwner},
		{g4LicOrg, g4LicDual},
		{g4LicOrgB, g4LicOwnerB},
	} {
		stakeholder := stakeholderModel.Stakeholder{
			Number: uuid.NewV4().String(), OrganisationNumber: spec.org, UserNumber: spec.user,
			StakeholderType: stakeholderModel.StakeholderTypeOwner, Status: "ACTIVE", AddedByNumber: g4LicOwner,
		}
		if err := db.Create(&stakeholder).Error; err != nil {
			t.Fatalf("seeding owner stakeholder for %s: %v", spec.user, err)
		}
	}

	// The explicit control-plane authority: manage_licence_registry.
	role := permissionsModel.Role{Number: g4LicRole, Name: g4LicRole, EntityNumber: g4LicOrg, Status: "ACTIVE", UserNumber: g4LicOwner}
	if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding licensing role: %v", err)
	}
	rolePermission := permissionsModel.RolePermission{RoleNumber: g4LicRole, Permission: "manage_licence_registry", AddedByNumber: g4LicOwner}
	if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
		t.Fatalf("seeding licensing permission: %v", err)
	}
	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, user := range []string{g4LicControl, g4LicDual} {
		userRole := permissionsModel.UserRole{UserNumber: user, RoleNumber: g4LicRole, ExpiryDate: expiry, AddedBy: g4LicOwner}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding licensing user role for %s: %v", user, err)
		}
	}

	return db
}

func g4Licence(number, org string) *licensingModel.Licence {
	return &licensingModel.Licence{
		Number: number, LicenceKey: g4LicenceKey, OrganisationNumber: org,
		ValidFrom: time.Now().Add(-time.Hour), ValidTo: time.Now().Add(24 * time.Hour),
		MaxInstallations: 5, LicensedAddons: "education", Status: "ACTIVE",
	}
}

func g4CountFor(db *gorm.DB, user, eventType string) int64 {
	var count int64
	db.Model(&notificationModel.Notification{}).
		Where("recipient_user_number = ? AND type = ?", user, eventType).Count(&count)
	return count
}

// TestFoundationG4LicensingNotifications covers the ratified licensing map
// (D-GQ1-4/5/6/16): client-facing events go to the active owner, control-plane
// events to holders of manage_licence_registry, suspension/revocation to both,
// with recipient dedup (D-GQ1-18), fail-closed on no owner, organisation
// isolation and no secret leakage.
func TestFoundationG4LicensingNotifications(t *testing.T) {
	db := setupG4LicensingFixture(t)

	t.Run("client_facing_event_notifies_owner_only", func(t *testing.T) {
		if err := licensingServices.RecordLicenceIssued(g4Licence(g4LicenceNumber, g4LicOrg), g4LicControl); err != nil {
			t.Fatalf("record issued: %v", err)
		}
		if !g4HasNotification(db, g4LicOwner, "LICENCE_ISSUED", false) {
			t.Fatal("owner was not notified of the issued licence")
		}
		// g4LicDual is also an owner, so is notified; g4LicControl (registry
		// only) and the generic Admin must not be.
		if g4CountFor(db, g4LicControl, "LICENCE_ISSUED") != 0 {
			t.Fatal("licensing-management recipient received a client-facing event")
		}
		if g4CountFor(db, g4LicAdmin, "LICENCE_ISSUED") != 0 {
			t.Fatal("a generic Admin received a client-facing event")
		}
		if g4CountFor(db, g4LicOwnerB, "LICENCE_ISSUED") != 0 {
			t.Fatal("another organisation's owner received the notification")
		}
	})

	t.Run("client_initiated_event_notifies_registry", func(t *testing.T) {
		if err := licensingServices.RecordClientRenewalRequested(g4Licence(g4LicenceNumber, g4LicOrg), g4LicOwner); err != nil {
			t.Fatalf("record renewal request: %v", err)
		}
		if !g4HasNotification(db, g4LicControl, "CLIENT_RENEWAL_REQUESTED", true) {
			t.Fatal("licensing-management recipient was not notified of the renewal request")
		}
		if !g4HasNotification(db, g4LicDual, "CLIENT_RENEWAL_REQUESTED", true) {
			t.Fatal("dual-role licensing-management recipient was not notified")
		}
		if g4CountFor(db, g4LicAdmin, "CLIENT_RENEWAL_REQUESTED") != 0 {
			t.Fatal("a generic Admin received a control-plane event")
		}
	})

	t.Run("suspension_notifies_owner_and_registry_once_each", func(t *testing.T) {
		if err := licensingServices.RecordLicenceTransition(licensingModel.LicenceEventSuspended, g4Licence(g4LicenceNumber, g4LicOrg), g4LicControl, "ACTIVE"); err != nil {
			t.Fatalf("record suspended: %v", err)
		}
		if !g4HasNotification(db, g4LicOwner, "LICENCE_SUSPENDED", true) {
			t.Fatal("owner was not told to act on suspension")
		}
		if !g4HasNotification(db, g4LicControl, "LICENCE_SUSPENDED", true) {
			t.Fatal("licensing-management recipient was not notified of suspension")
		}
		// D-GQ1-18: g4LicDual qualifies as owner AND registry authority and
		// must receive exactly one notification.
		if count := g4CountFor(db, g4LicDual, "LICENCE_SUSPENDED"); count != 1 {
			t.Fatalf("dual-role recipient got %d suspension notifications, want 1", count)
		}
	})

	t.Run("installation_event_notifies_registry", func(t *testing.T) {
		installation := &licensingModel.Installation{
			Number: "g4-installation-1", InstallationIdentity: "g4-iid-1",
			OrganisationNumber: g4LicOrg, LicenceNumber: g4LicenceNumber, Status: "ACTIVE",
		}
		if err := licensingServices.RecordInstallationActivated(installation, g4LicOwner); err != nil {
			t.Fatalf("record installation activated: %v", err)
		}
		if !g4HasNotification(db, g4LicControl, "INSTALLATION_ACTIVATED", false) {
			t.Fatal("licensing-management recipient was not notified of activation")
		}
		if g4CountFor(db, g4LicOwner, "INSTALLATION_ACTIVATED") != 0 {
			t.Fatal("client owner received a control-plane installation event")
		}
	})

	t.Run("fail_closed_without_owner", func(t *testing.T) {
		before := g4CountOrg(db, g4LicOrphan)
		if err := licensingServices.RecordLicenceIssued(g4Licence(g4LicenceNumber2, g4LicOrphan), g4LicControl); err != nil {
			t.Fatalf("record orphan issued: %v", err)
		}
		// No owner exists: no notification, and no fallback recipient invented.
		if g4CountOrg(db, g4LicOrphan) != before {
			t.Fatal("a notification was created for an organisation with no owner")
		}
		var events int64
		db.Model(&licensingModel.LicenceEvent{}).Where("organisation_number = ?", g4LicOrphan).Count(&events)
		if events == 0 {
			t.Fatal("the authoritative licensing fact was not recorded")
		}
	})

	t.Run("no_secret_leaks_into_notifications", func(t *testing.T) {
		var rows []notificationModel.Notification
		db.Where("organisation_number = ?", g4LicOrg).Find(&rows)
		if len(rows) == 0 {
			t.Fatal("expected notifications to inspect")
		}
		for _, row := range rows {
			if strings.Contains(row.Title, g4LicenceKey) || strings.Contains(row.Body, g4LicenceKey) || strings.Contains(row.Payload, g4LicenceKey) {
				t.Fatalf("licence key leaked into notification %d", row.ID)
			}
		}
	})
}

func g4CountOrg(db *gorm.DB, organisationNumber string) int64 {
	var count int64
	db.Model(&notificationModel.Notification{}).
		Where("organisation_number = ?", organisationNumber).Count(&count)
	return count
}
