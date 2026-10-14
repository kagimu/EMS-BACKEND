package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	ntViewerA     = "ems-nt-viewer-a"
	ntViewerB     = "ems-nt-viewer-b"
	ntAdminNoPerm = "ems-nt-admin-noperm"
	ntNoPerm      = "ems-nt-noperm"
	ntNoCtx       = "ems-nt-nocontext"

	ntOrgA = "ems-nt-org-a"
	ntOrgB = "ems-nt-org-b"

	ntRoleView = "ems-nt-role-view"

	ntNotificationsURL = "/main/ems/notifications"
)

type ntFixture struct {
	db          *gorm.DB
	tokenA      string
	tokenB      string
	tokenAdmin  string
	tokenNoPerm string
	tokenNoCtx  string

	notifA1    uint
	notifA2    uint
	notifAdmin uint
	notifB1    uint
}

func setupNotificationFixture(t *testing.T) *ntFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-g-notification-test-secret")

	users := []string{ntViewerA, ntViewerB, ntAdminNoPerm, ntNoPerm, ntNoCtx}
	organisations := []string{ntOrgA, ntOrgB}
	roles := []string{ntRoleView}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&notificationModel.Notification{})
		db.Unscoped().Where("role_number IN ?", roles).Delete(&permissionsModel.RolePermission{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("number IN ?", roles).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct{ number, userType string }{
		{ntViewerA, "User"},
		{ntViewerB, "User"},
		{ntAdminNoPerm, "Admin"},
		{ntNoPerm, "User"},
		{ntNoCtx, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationG",
			SecondName: "Notification",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{ntOrgA, "Foundation G Notification Org A"},
		{ntOrgB, "Foundation G Notification Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: ntViewerA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{ntViewerA, ntOrgA},
		{ntAdminNoPerm, ntOrgA},
		{ntNoPerm, ntOrgA},
		{ntViewerB, ntOrgB},
		// ntNoCtx intentionally has no active entity.
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	role := permissionsModel.Role{
		Number: ntRoleView, Name: ntRoleView, EntityNumber: ntOrgA, Status: "ACTIVE", UserNumber: ntViewerA,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding role: %v", err)
	}
	rolePermission := permissionsModel.RolePermission{
		RoleNumber: ntRoleView, Permission: "view_ems_notifications", AddedByNumber: ntViewerA,
	}
	if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
		t.Fatalf("seeding permission: %v", err)
	}
	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, userNumber := range []string{ntViewerA, ntViewerB, ntNoCtx} {
		userRole := permissionsModel.UserRole{
			UserNumber: userNumber, RoleNumber: ntRoleView, ExpiryDate: expiry, AddedBy: ntViewerA,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding user role for %s: %v", userNumber, err)
		}
	}

	now := time.Now()
	seed := func(org, recipient, notificationType, title string, read bool) uint {
		notification := notificationModel.Notification{
			OrganisationNumber:  org,
			RecipientUserNumber: recipient,
			Type:                notificationType,
			Title:               title,
			Body:                "body of " + title,
		}
		if read {
			notification.ReadAt = &now
		}
		if err := db.Create(&notification).Error; err != nil {
			t.Fatalf("seeding notification for %s: %v", recipient, err)
		}
		return notification.ID
	}

	notifA1 := seed(ntOrgA, ntViewerA, "approval.assigned", "Action required", false)
	notifA2 := seed(ntOrgA, ntViewerA, "approval.completed", "Approved", true)
	notifAdmin := seed(ntOrgA, ntAdminNoPerm, "approval.assigned", "Admin action", false)
	notifB1 := seed(ntOrgB, ntViewerB, "licence.released", "Licence released", false)

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &ntFixture{
		db:          db,
		tokenA:      token(ntViewerA),
		tokenB:      token(ntViewerB),
		tokenAdmin:  token(ntAdminNoPerm),
		tokenNoPerm: token(ntNoPerm),
		tokenNoCtx:  token(ntNoCtx),
		notifA1:     notifA1,
		notifA2:     notifA2,
		notifAdmin:  notifAdmin,
		notifB1:     notifB1,
	}
}

type ntNotificationData struct {
	Number              uint            `json:"number"`
	OrganisationNumber  string          `json:"organisation_number"`
	RecipientUserNumber string          `json:"recipient_user_number"`
	Type                string          `json:"type"`
	Title               string          `json:"title"`
	Body                string          `json:"body"`
	Payload             json.RawMessage `json:"payload"`
	ReadAt              *time.Time      `json:"read_at"`
	CreatedAt           time.Time       `json:"created_at"`
}

type ntEnvelope struct {
	Message string             `json:"message"`
	Data    ntNotificationData `json:"data"`
}

type ntListEnvelope struct {
	Message string               `json:"message"`
	Data    []ntNotificationData `json:"data"`
	Total   int64                `json:"total"`
}

// TestFoundationGNotification covers Foundation G Increment 3: the
// recipient-scoped in-app notification store, recipient and organisation
// isolation, read/unread behaviour, mark-read (which changes only read_at) and
// permission enforcement. No external delivery, scheduler or worker is involved.
func TestFoundationGNotification(t *testing.T) {
	fx := setupNotificationFixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	do := func(token, method, path string, payload any) (int, []byte) {
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
		req.Header.Set("Authorization", "Bearer "+token)
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

	list := func(token, query string) ntListEnvelope {
		t.Helper()
		status, body := do(token, http.MethodGet, ntNotificationsURL+query, nil)
		if status != http.StatusOK {
			t.Fatalf("list %q = %d (%s)", query, status, string(body))
		}
		var envelope ntListEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding list: %v (%s)", err, string(body))
		}
		return envelope
	}

	t.Run("list_is_recipient_and_org_scoped", func(t *testing.T) {
		envelope := list(fx.tokenA, "")
		if envelope.Total != 2 {
			t.Fatalf("viewer A total = %d, want 2", envelope.Total)
		}
		for _, item := range envelope.Data {
			if item.RecipientUserNumber != ntViewerA {
				t.Fatalf("viewer A received %s's notification", item.RecipientUserNumber)
			}
			if item.OrganisationNumber != ntOrgA {
				t.Fatalf("viewer A received org %s notification", item.OrganisationNumber)
			}
		}
	})

	t.Run("unread_filter", func(t *testing.T) {
		envelope := list(fx.tokenA, "?unread=true")
		if envelope.Total != 1 {
			t.Fatalf("unread total = %d, want 1", envelope.Total)
		}
		if envelope.Data[0].ReadAt != nil {
			t.Fatal("unread filter returned a read notification")
		}
	})

	t.Run("get_own_notification", func(t *testing.T) {
		status, body := do(fx.tokenA, http.MethodGet, ntNotificationsURL+"/"+itoa(fx.notifA1), nil)
		if status != http.StatusOK {
			t.Fatalf("get own = %d (%s)", status, string(body))
		}
		var envelope ntEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if envelope.Data.Number != fx.notifA1 || envelope.Data.RecipientUserNumber != ntViewerA {
			t.Fatalf("unexpected notification: %+v", envelope.Data)
		}
	})

	t.Run("cannot_get_another_recipient", func(t *testing.T) {
		// Same organisation, different recipient.
		if status, _ := do(fx.tokenA, http.MethodGet, ntNotificationsURL+"/"+itoa(fx.notifAdmin), nil); status != http.StatusNotFound {
			t.Fatalf("cross-recipient get = %d, want 404", status)
		}
	})

	t.Run("cannot_get_another_organisation", func(t *testing.T) {
		if status, _ := do(fx.tokenA, http.MethodGet, ntNotificationsURL+"/"+itoa(fx.notifB1), nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want 404", status)
		}
	})

	t.Run("admin_sees_only_own_notifications", func(t *testing.T) {
		envelope := list(fx.tokenAdmin, "")
		if envelope.Total != 1 {
			t.Fatalf("admin total = %d, want 1", envelope.Total)
		}
		if envelope.Data[0].Number != fx.notifAdmin {
			t.Fatalf("admin saw notification %d, want %d", envelope.Data[0].Number, fx.notifAdmin)
		}
	})

	t.Run("mark_read_changes_only_read_at", func(t *testing.T) {
		var before notificationModel.Notification
		if err := fx.db.Where("id = ?", fx.notifA1).First(&before).Error; err != nil {
			t.Fatalf("loading before: %v", err)
		}

		status, body := do(fx.tokenA, http.MethodPost, ntNotificationsURL+"/"+itoa(fx.notifA1)+"/read", map[string]any{})
		if status != http.StatusOK {
			t.Fatalf("mark read = %d (%s)", status, string(body))
		}
		var envelope ntEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if envelope.Data.ReadAt == nil {
			t.Fatal("read_at was not set")
		}
		firstRead := *envelope.Data.ReadAt

		var after notificationModel.Notification
		if err := fx.db.Where("id = ?", fx.notifA1).First(&after).Error; err != nil {
			t.Fatalf("loading after: %v", err)
		}
		if after.Title != before.Title || after.Body != before.Body ||
			after.Type != before.Type || after.RecipientUserNumber != before.RecipientUserNumber {
			t.Fatal("mark-read mutated notification content")
		}
		if !after.CreatedAt.Equal(before.CreatedAt) {
			t.Fatal("mark-read mutated created_at")
		}

		// Idempotent: a second mark-read does not change read_at.
		status, body = do(fx.tokenA, http.MethodPost, ntNotificationsURL+"/"+itoa(fx.notifA1)+"/read", nil)
		if status != http.StatusOK {
			t.Fatalf("second mark read = %d (%s)", status, string(body))
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding second: %v", err)
		}
		if envelope.Data.ReadAt == nil || !envelope.Data.ReadAt.Equal(firstRead) {
			t.Fatal("mark-read is not idempotent")
		}
	})

	t.Run("cannot_mark_read_another_recipient", func(t *testing.T) {
		if status, _ := do(fx.tokenA, http.MethodPost, ntNotificationsURL+"/"+itoa(fx.notifAdmin)+"/read", nil); status != http.StatusNotFound {
			t.Fatalf("cross-recipient mark-read = %d, want 404", status)
		}
	})

	t.Run("permission_required", func(t *testing.T) {
		if status, _ := do(fx.tokenNoPerm, http.MethodGet, ntNotificationsURL, nil); status != http.StatusForbidden {
			t.Fatalf("no-permission list = %d, want 403", status)
		}
	})

	t.Run("no_active_organisation", func(t *testing.T) {
		if status, _ := do(fx.tokenNoCtx, http.MethodGet, ntNotificationsURL, nil); status != http.StatusConflict {
			t.Fatalf("no-context list = %d, want 409", status)
		}
	})
}

// itoa renders a notification identifier for URL construction.
func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
