package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	authMiddleware "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/middleware"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	"github.com/gin-gonic/gin"
)

type sec1View struct {
	Number             string `json:"number"`
	OrganisationNumber string `json:"organisation_number"`
	RequestedByNumber  string `json:"requested_by_number"`
}

// TestSEC1IdentityCannotBeSpoofed is the SEC-1 regression: a valid JWT for User A
// combined with a client-supplied User-Number header for User B must still
// resolve to User A, both for the application identity and for the raw
// User-Number header that legacy/Education controllers read.
func TestSEC1IdentityCannotBeSpoofed(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
	g4CleanNotifications(t, fx.db, organisations)
	t.Cleanup(func() {
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		g4CleanNotifications(t, fx.db, organisations)
	})

	h := newG4Harness(t)

	do := func(token, method, path string, payload any, headers map[string][]string) (int, []byte) {
		t.Helper()
		var reader io.Reader
		if payload != nil {
			raw, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("encoding payload: %v", err)
			}
			reader = bytes.NewReader(raw)
		}
		req, err := http.NewRequest(method, h.server.URL+path, reader)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		resp, err := h.client.Do(req)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	createPayload := map[string]any{
		"module": "ems", "resource_type": "sec1", "resource_number": "1",
		"decision_type": "purchase_request", "amount": 100,
	}
	mustCreateAsA := func(t *testing.T, headers map[string][]string) sec1View {
		t.Helper()
		status, body := do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, createPayload, headers)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		var envelope struct {
			Data sec1View `json:"data"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding create: %v", err)
		}
		return envelope.Data
	}
	assertResolvedToA := func(t *testing.T, view sec1View) {
		t.Helper()
		if view.OrganisationNumber != eOrgA {
			t.Fatalf("organisation = %q, want %q (spoof must not select org B)", view.OrganisationNumber, eOrgA)
		}
		if view.RequestedByNumber != eAdminA {
			t.Fatalf("requester = %q, want %q (spoof must not select user B)", view.RequestedByNumber, eAdminA)
		}
		var asA int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND entity_id = ? AND actor_user_number = ? AND organisation_number = ?",
				"approval.create", view.Number, eAdminA, eOrgA).Count(&asA)
		if asA != 1 {
			t.Fatalf("audit approval.create actor must be the JWT user %q (found %d)", eAdminA, asA)
		}
		var asB int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("entity_id = ? AND actor_user_number = ?", view.Number, eAdminB).Count(&asB)
		if asB != 0 {
			t.Fatalf("audit must not attribute the action to the spoofed user %q", eAdminB)
		}
		var underB int64
		fx.db.Table("ems_approval_requests").
			Where("number = ? AND organisation_number = ?", view.Number, eOrgB).Count(&underB)
		if underB != 0 {
			t.Fatal("spoofed identity must not create a record under organisation B")
		}
	}

	t.Run("valid_jwt_without_spoof_resolves_jwt_user", func(t *testing.T) {
		assertResolvedToA(t, mustCreateAsA(t, nil))
	})

	t.Run("jwt_a_with_user_number_b_still_resolves_a", func(t *testing.T) {
		assertResolvedToA(t, mustCreateAsA(t, map[string][]string{"User-Number": {eAdminB}}))
	})

	t.Run("multiple_client_user_numbers_cannot_override", func(t *testing.T) {
		assertResolvedToA(t, mustCreateAsA(t, map[string][]string{"User-Number": {eAdminB, eAdminB}}))
	})

	t.Run("missing_jwt_rejected", func(t *testing.T) {
		if status, _ := do("", http.MethodPost, eApprovalsBase, createPayload, nil); status != http.StatusUnauthorized {
			t.Fatalf("missing JWT = %d, want 401", status)
		}
	})

	t.Run("invalid_jwt_rejected", func(t *testing.T) {
		if status, _ := do("not-a-valid-token", http.MethodPost, eApprovalsBase, createPayload, nil); status != http.StatusUnauthorized {
			t.Fatalf("invalid JWT = %d, want 401", status)
		}
	})

	// The raw User-Number header (read directly by legacy/Education controllers)
	// must also carry the JWT-derived value, never the client-supplied one.
	t.Run("raw_user_number_header_is_jwt_derived", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		engine := gin.New()
		engine.GET("/whoami", authMiddleware.AuthorizeJWT(), func(c *gin.Context) {
			user := c.MustGet("user").(authModel.User)
			c.JSON(http.StatusOK, gin.H{
				"user":   user.Number,
				"header": c.GetHeader("User-Number"),
				"all":    c.Request.Header.Values("User-Number"),
			})
		})
		server := httptest.NewServer(engine)
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL+"/whoami", nil)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+fx.tokenAdminA)
		req.Header.Add("User-Number", eAdminB)
		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var got struct {
			User   string   `json:"user"`
			Header string   `json:"header"`
			All    []string `json:"all"`
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("decoding whoami: %v", err)
		}
		if got.User != eAdminA {
			t.Fatalf("context identity = %q, want %q", got.User, eAdminA)
		}
		if got.Header != eAdminA {
			t.Fatalf("User-Number header = %q, want %q", got.Header, eAdminA)
		}
		if len(got.All) != 1 || got.All[0] != eAdminA {
			t.Fatalf("User-Number header values = %v, want exactly [%q]", got.All, eAdminA)
		}
	})
}
