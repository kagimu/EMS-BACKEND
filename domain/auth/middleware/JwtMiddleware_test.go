package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRecorderContext(method, target string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, target, nil)

	return c, recorder
}

// TestAuthorizeJWTRejectsMissingAuthorization covers the deterministic,
// database-independent contract of the middleware: a request without an
// Authorization header and without a query token is rejected with 401.
func TestAuthorizeJWTRejectsMissingAuthorization(t *testing.T) {
	c, recorder := newRecorderContext(http.MethodGet, "/main/profile")

	AuthorizeJWT()(c)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
