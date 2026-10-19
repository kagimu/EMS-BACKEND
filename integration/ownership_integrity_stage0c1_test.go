package integration

import (
	"bytes"
	"context"
	"database/sql"
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
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	stakeholdersRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/repositories"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"gorm.io/gorm"
)

const s0c1StakePath = "/main/ems/stakeholders"

type stage0c1Fixture struct {
	db *gorm.DB
}

type stage0c1Organisation struct {
	org   string
	token string
}

func setupStage0C1Fixture(t *testing.T) *stage0c1Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	t.Setenv("JWT_SECRET", "stage-0c1-test-secret")

	return &stage0c1Fixture{db: db}
}

// newOrganisation creates an isolated organisation, an Admin who is its active
// member, and the requested owner users. Every subtest uses its own organisation
// so the effective-owner count is never polluted by another subtest.
func (fx *stage0c1Fixture) newOrganisation(t *testing.T, suffix string, ownerUsers ...string) *stage0c1Organisation {
	t.Helper()

	org := "stage0c1-org-" + suffix
	admin := "stage0c1-admin-" + suffix
	users := append([]string{admin}, ownerUsers...)

	cleanup := func() {
		fx.db.Unscoped().Where("organisation_number = ?", org).Delete(&stakeholderModel.Stakeholder{})
		fx.db.Unscoped().Where("organisation_number = ?", org).Delete(&auditModel.AuditLog{})
		fx.db.Unscoped().Where("organisation_number = ?", org).Delete(&notificationModel.Notification{})
		fx.db.Unscoped().Where("entity_number = ?", org).Delete(&eosModel.EntityUser{})
		fx.db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		fx.db.Unscoped().Where("number = ?", org).Delete(&eosModel.Entity{})
		fx.db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, number := range users {
		userType := "User"
		if number == admin {
			userType = "Admin"
		}
		user := authModel.User{
			Number:     number,
			FirstName:  "Stage0C1",
			SecondName: number,
			Gender:     "M",
			Type:       userType,
			Status:     "ACTIVE",
		}
		if err := fx.db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	entity := eosModel.Entity{Number: org, Name: "Stage 0C1 " + suffix, AddedBy: admin}
	if err := fx.db.Create(&entity).Error; err != nil {
		t.Fatalf("seeding entity %s: %v", org, err)
	}

	membership := eosModel.EntityUser{
		EntityNumber:  org,
		UserNumber:    admin,
		AddedByNumber: admin,
	}
	if err := fx.db.Omit("Entity", "User", "AddedBy").Create(&membership).Error; err != nil {
		t.Fatalf("seeding membership: %v", err)
	}

	active := eosModel.ActiveUserEntity{UserNumber: admin, EntityNumber: org}
	if err := fx.db.Omit("Entity", "User").Create(&active).Error; err != nil {
		t.Fatalf("setting active organisation: %v", err)
	}

	token, err := authServices.JWTAuthService().GenerateToken(admin)
	if err != nil {
		t.Fatalf("generating token: %v", err)
	}

	return &stage0c1Organisation{org: org, token: token}
}

func TestOwnershipIntegrityStage0C1(t *testing.T) {
	fx := setupStage0C1Fixture(t)

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

	ownerPayload := func(userNumber string, validFrom, validTo *time.Time) map[string]any {
		payload := map[string]any{
			"user_number":      userNumber,
			"stakeholder_type": "owner",
		}
		if validFrom != nil {
			payload["valid_from"] = validFrom.UTC().Format(time.RFC3339)
		}
		if validTo != nil {
			payload["valid_to"] = validTo.UTC().Format(time.RFC3339)
		}
		return payload
	}

	createOwner := func(token, userNumber string, validFrom, validTo *time.Time) string {
		t.Helper()
		status, body := must(http.MethodPost, s0c1StakePath, token, ownerPayload(userNumber, validFrom, validTo))
		if status != http.StatusCreated {
			t.Fatalf("create owner %s = %d (%s)", userNumber, status, string(body))
		}

		var parsed struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding owner response: %v", err)
		}
		if parsed.Data.Number == "" {
			t.Fatalf("owner number missing: %s", string(body))
		}
		return parsed.Data.Number
	}

	setStatus := func(token, number, status string) (int, []byte) {
		return must(http.MethodPatch, s0c1StakePath+"/"+number+"/status", token, map[string]string{"status": status})
	}

	updateValidity := func(token, number string, validTo time.Time) (int, []byte) {
		return must(http.MethodPut, s0c1StakePath+"/"+number, token, map[string]any{
			"valid_to": validTo.UTC().Format(time.RFC3339),
		})
	}

	effectiveOwners := func(org string) int64 {
		count, err := stakeholdersRepositories.CountEffectiveOwners(org, time.Now())
		if err != nil {
			t.Fatalf("counting effective owners: %v", err)
		}
		return count
	}

	t.Run("A_single_effective_owner_cannot_be_deactivated", func(t *testing.T) {
		o := fx.newOrganisation(t, "a", "stage0c1-owner-a")
		owner := createOwner(o.token, "stage0c1-owner-a", nil, nil)
		if status, body := setStatus(o.token, owner, "INACTIVE"); status != http.StatusConflict {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("B_single_effective_owner_cannot_be_expired_by_update", func(t *testing.T) {
		o := fx.newOrganisation(t, "b", "stage0c1-owner-b")
		owner := createOwner(o.token, "stage0c1-owner-b", nil, nil)
		if status, body := updateValidity(o.token, owner, time.Now().Add(-time.Hour)); status != http.StatusConflict {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("C_two_effective_owners_can_deactivate_one", func(t *testing.T) {
		o := fx.newOrganisation(t, "c", "stage0c1-owner-c1", "stage0c1-owner-c2")
		first := createOwner(o.token, "stage0c1-owner-c1", nil, nil)
		second := createOwner(o.token, "stage0c1-owner-c2", nil, nil)
		if status, body := setStatus(o.token, first, "INACTIVE"); status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if status, body := must(http.MethodGet, s0c1StakePath+"/"+second, o.token, nil); status != http.StatusOK {
			t.Fatalf("second owner = %d (%s)", status, string(body))
		}
		if remaining := effectiveOwners(o.org); remaining != 1 {
			t.Fatalf("effective owners = %d, want 1", remaining)
		}
	})

	t.Run("D_two_effective_owners_can_shorten_one_validity", func(t *testing.T) {
		o := fx.newOrganisation(t, "d", "stage0c1-owner-d1", "stage0c1-owner-d2")
		first := createOwner(o.token, "stage0c1-owner-d1", nil, nil)
		createOwner(o.token, "stage0c1-owner-d2", nil, nil)
		if status, body := updateValidity(o.token, first, time.Now().Add(-time.Hour)); status != http.StatusOK {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if remaining := effectiveOwners(o.org); remaining != 1 {
			t.Fatalf("effective owners = %d, want 1", remaining)
		}
	})

	t.Run("E_expired_owner_is_not_counted", func(t *testing.T) {
		o := fx.newOrganisation(t, "e", "stage0c1-owner-e1", "stage0c1-owner-e2")
		effective := createOwner(o.token, "stage0c1-owner-e1", nil, nil)
		expired := createOwner(o.token, "stage0c1-owner-e2", nil, timePointer(time.Now().Add(-time.Hour)))

		if status, body := setStatus(o.token, effective, "INACTIVE"); status != http.StatusConflict {
			t.Fatalf("deactivating the only effective owner = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
		// The expired owner is not effective, so ending it is allowed.
		if status, body := setStatus(o.token, expired, "INACTIVE"); status != http.StatusOK {
			t.Fatalf("deactivating expired owner = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("F_future_owner_is_not_counted", func(t *testing.T) {
		o := fx.newOrganisation(t, "f", "stage0c1-owner-f1", "stage0c1-owner-f2")
		effective := createOwner(o.token, "stage0c1-owner-f1", nil, nil)
		createOwner(o.token, "stage0c1-owner-f2", timePointer(time.Now().Add(24*time.Hour)), nil)

		if status, body := setStatus(o.token, effective, "INACTIVE"); status != http.StatusConflict {
			t.Fatalf("status = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}
	})

	t.Run("G_concurrent_final_deactivations_keep_one_owner", func(t *testing.T) {
		o := fx.newOrganisation(t, "g", "stage0c1-owner-g1", "stage0c1-owner-g2")
		first := createOwner(o.token, "stage0c1-owner-g1", nil, nil)
		second := createOwner(o.token, "stage0c1-owner-g2", nil, nil)

		statuses := make(chan int, 2)
		var wg sync.WaitGroup
		for _, number := range []string{first, second} {
			wg.Add(1)
			go func(number string) {
				defer wg.Done()
				status, _, err := doRequest(http.MethodPatch, s0c1StakePath+"/"+number+"/status", o.token, map[string]string{"status": "INACTIVE"})
				if err != nil {
					status = 0
				}
				statuses <- status
			}(number)
		}
		wg.Wait()
		close(statuses)

		ok, conflict := 0, 0
		for status := range statuses {
			switch status {
			case http.StatusOK:
				ok++
			case http.StatusConflict:
				conflict++
			}
		}
		if ok != 1 || conflict != 1 {
			t.Fatalf("concurrent deactivations: %d ok, %d conflict; want 1 and 1", ok, conflict)
		}
		if remaining := effectiveOwners(o.org); remaining < 1 {
			t.Fatalf("effective owners after concurrent deactivation = %d, want >= 1", remaining)
		}
	})

	t.Run("H_concurrent_duplicate_owner_creation_is_serialised", func(t *testing.T) {
		o := fx.newOrganisation(t, "h", "stage0c1-owner-h")
		statuses := make(chan int, 2)
		var wg sync.WaitGroup
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				status, _, err := doRequest(http.MethodPost, s0c1StakePath, o.token, ownerPayload("stage0c1-owner-h", nil, nil))
				if err != nil {
					status = 0
				}
				statuses <- status
			}()
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
			t.Fatalf("concurrent owner creation: %d created, %d conflict; want 1 and 1", created, conflict)
		}

		var active int64
		if err := fx.db.Model(&stakeholderModel.Stakeholder{}).
			Where("organisation_number = ? AND user_number = ? AND stakeholder_type = ? AND status = ?",
				o.org, "stage0c1-owner-h", stakeholderModel.StakeholderTypeOwner, "ACTIVE").
			Count(&active).Error; err != nil {
			t.Fatalf("counting duplicate owner rows: %v", err)
		}
		if active != 1 {
			t.Fatalf("ACTIVE owner rows = %d, want 1", active)
		}
	})

	t.Run("I_lock_timeout_leaves_ownership_unchanged", func(t *testing.T) {
		o := fx.newOrganisation(t, "i", "stage0c1-owner-i1", "stage0c1-owner-i2")
		first := createOwner(o.token, "stage0c1-owner-i1", nil, nil)
		createOwner(o.token, "stage0c1-owner-i2", nil, nil)

		sqlDB, err := fx.db.DB()
		if err != nil {
			t.Fatalf("getting sql.DB: %v", err)
		}
		conn, err := sqlDB.Conn(context.Background())
		if err != nil {
			t.Fatalf("pinning connection: %v", err)
		}
		lockName := stakeholdersRepositories.OwnershipLockName(o.org)
		var got sql.NullInt64
		if err := conn.QueryRowContext(context.Background(), "SELECT GET_LOCK(?, 0)", lockName).Scan(&got); err != nil || !got.Valid || got.Int64 != 1 {
			conn.Close()
			t.Fatalf("could not take the ownership lock for the test: got=%v err=%v", got, err)
		}
		defer func() {
			_, _ = conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK(?)", lockName)
			conn.Close()
		}()

		if status, body := setStatus(o.token, first, "INACTIVE"); status != http.StatusConflict {
			t.Fatalf("status while lock held = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}

		var stored stakeholderModel.Stakeholder
		if err := fx.db.Where("number = ?", first).First(&stored).Error; err != nil {
			t.Fatalf("loading owner: %v", err)
		}
		if stored.Status != "ACTIVE" {
			t.Fatalf("owner status after failed lock = %q, want ACTIVE", stored.Status)
		}
	})

	t.Run("J_legitimate_create_update_and_deactivate_still_work", func(t *testing.T) {
		o := fx.newOrganisation(t, "j", "stage0c1-owner-j1", "stage0c1-owner-j2")
		first := createOwner(o.token, "stage0c1-owner-j1", nil, nil)
		second := createOwner(o.token, "stage0c1-owner-j2", nil, nil)

		if status, body := must(http.MethodPut, s0c1StakePath+"/"+first, o.token, map[string]any{
			"title": "Chairperson",
		}); status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}
		if status, body := setStatus(o.token, second, "INACTIVE"); status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
		if status, body := setStatus(o.token, second, "ACTIVE"); status != http.StatusOK {
			t.Fatalf("reactivate = %d (%s)", status, string(body))
		}
	})
}

func timePointer(value time.Time) *time.Time {
	return &value
}
