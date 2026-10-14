package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
)

// apiClient is a minimal authenticated HTTP client for the regression flows.
type apiClient struct {
	t      *testing.T
	base   string
	token  string
	userNo string
	http   *http.Client
}

// login boots the real router against the test database and authenticates via
// the real /auth/login endpoint.
func login(t *testing.T) *apiClient {
	t.Helper()

	testutil.RequireTestDB(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)

	c := &apiClient{t: t, base: server.URL, http: server.Client()}

	email, password := testutil.LoginCredentials(t)
	payload, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		t.Fatalf("encoding login request: %v", err)
	}

	resp, err := c.http.Post(c.base+"/auth/login", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("login returned %d, want 200: %s", resp.StatusCode, string(raw))
	}

	var body struct {
		Token string `json:"token"`
		Data  struct {
			Number      string   `json:"number"`
			Permissions []string `json:"permissions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decoding login response: %v", err)
	}
	if body.Token == "" {
		t.Fatal("login response did not contain a token")
	}

	c.token = body.Token
	c.userNo = body.Data.Number

	return c
}

func (c *apiClient) get(path string) *http.Response {
	c.t.Helper()

	req, err := http.NewRequest(http.MethodGet, c.base+path, nil)
	if err != nil {
		c.t.Fatalf("building GET %s: %v", path, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("user-number", c.userNo)

	resp, err := c.http.Do(req)
	if err != nil {
		c.t.Fatalf("GET %s failed: %v", path, err)
	}

	return resp
}

// expectReachable asserts that an authenticated request is accepted (not 401
// or 403) and does not produce a server error (not 5xx).
//
// Release 1 deliberately asserts reachability rather than an exact body or
// status, because several Education endpoints depend on query parameters and
// seeded data that this foundation release does not fabricate. Tightening these
// assertions to exact payloads requires purpose-built fixtures and is deferred.
func (c *apiClient) expectReachable(path string) {
	c.t.Helper()

	resp := c.get(path)
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		c.t.Errorf("GET %s = 401; authentication was rejected", path)
	case resp.StatusCode == http.StatusForbidden:
		c.t.Errorf("GET %s = 403; authorization was rejected", path)
	case resp.StatusCode >= 500:
		raw, _ := io.ReadAll(resp.Body)
		c.t.Errorf("GET %s = %d; server error: %s", path, resp.StatusCode, string(raw))
	}
}

// TestEducationRegressionFlows is the Release 1 regression foundation for the
// eight critical Education flows.
//
// Every subtest requires TEST_APP_DB_HOST and TEST_LOGIN_EMAIL /
// TEST_LOGIN_PASSWORD to be configured, and is skipped otherwise. The tests
// exercise the existing HTTP contracts and require no production changes.
func TestEducationRegressionFlows(t *testing.T) {
	t.Run("1_auth_contract_login", func(t *testing.T) {
		c := login(t)
		if c.userNo == "" {
			t.Error("login response data.number was empty")
		}
	})

	t.Run("2_scoping_isolation_campus_list", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/campus/all")
	})

	t.Run("3_portal_profile_reachable", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/profile")
	})

	t.Run("4_finance_invoices", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/invoices")
	})

	t.Run("5_student_lifecycle_applications", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/applications/all")
	})

	t.Run("6_assessment_areas", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/assessment-areas")
	})

	t.Run("7_catering_meal_calendar", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/catering/meal-calendar")
	})

	t.Run("8_curriculum_contract", func(t *testing.T) {
		c := login(t)
		c.expectReachable("/main/curriculum/all")
	})
}
