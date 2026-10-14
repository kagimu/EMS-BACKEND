package middleware

import (
	"errors"
	"net/http"
	"strings"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

// licenceExemptPrefixes lists paths that must remain reachable even when the
// licence is expired past grace, suspended, revoked or not entitled (F-Q2,
// §10/§11). EMS Core routes are not licensed addons and are not gated.
var licenceExemptPrefixes = []string{
	"/main/ems/", // EMS Core + licensing status/activation/renewal/registry
	"/main/profile",
}

// bootstrapExemptEndpoints lists the minimal method+path pairs required to
// bootstrap the FIRST organisation of a brand-new installation. These routes
// must be reachable before any active organisation (and therefore before any
// licence) can exist; otherwise the licence gate deadlocks first-time setup.
// The exemption is exact (method AND path) and skips ONLY the licence gate —
// JWT authentication, membership and handler validation still apply.
var bootstrapExemptEndpoints = map[string]struct{}{
	http.MethodPost + " /main/eos/entity/add":     {},
	http.MethodGet + " /main/eos/my-entities":     {},
	http.MethodPost + " /main/eos/entity/current": {},
}

// devToolExemptEndpoints lists exact developer-only endpoints that must remain
// reachable regardless of licence state so that the disposable acceptance
// database can be reset to a clean slate (which by definition has no active
// organisation or licence). These endpoints are only registered when the
// developer feature is explicitly enabled, and they still require JWT
// authentication plus their own authorisation checks — the licence gate is the
// only check skipped here.
var devToolExemptEndpoints = map[string]struct{}{
	http.MethodPost + " /main/dev/reset-onboarding": {},
}

// RequireLicenceEntitlement is the additive Foundation F enforcement boundary.
// It answers only whether the authenticated organisation's current licence
// permits use of the given licensed module; it never grants permissions, roles
// or authority, and it never trusts a client-supplied organisation.
//
// It is deliberately not an authorization mechanism: existing authentication
// and permission middleware continue to run after it (or are skipped by the
// exemption list for licensing functionality).
func RequireLicenceEntitlement(moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, prefix := range licenceExemptPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

		// Exact method+path exemptions for first-time organisation bootstrap.
		if _, exempt := bootstrapExemptEndpoints[c.Request.Method+" "+path]; exempt {
			c.Next()
			return
		}

		// Exact method+path exemptions for developer-only tooling.
		if _, exempt := devToolExemptEndpoints[c.Request.Method+" "+path]; exempt {
			c.Next()
			return
		}

		value, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}
		user, ok := value.(authModels.User)
		if !ok || user.Number == "" {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}

		// Organisation context comes exclusively from the authenticated user
		// (via the existing eos ActiveUserEntity mechanism). A request-supplied
		// organisation is never consulted.
		organisationNumber, err := licensingServices.OrganisationNumberForUser(user.Number)
		if err != nil {
			if errors.Is(err, licensingServices.ErrNoActiveOrganisation) {
				c.JSON(http.StatusConflict, commons.Feedback("Your EMS licence does not permit access to this module"))
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to verify the EMS licence"))
			c.Abort()
			return
		}

		// Reuse the Increment 1 predicate: stored ACTIVE, within
		// [valid_from, valid_to + grace), and the module in licensed_addons.
		licensed, err := licensingServices.IsAddonLicensed(organisationNumber, moduleCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to verify the EMS licence"))
			c.Abort()
			return
		}
		if !licensed {
			c.JSON(http.StatusConflict, commons.Feedback("Your EMS licence does not permit access to this module"))
			c.Abort()
			return
		}

		c.Next()
	}
}
