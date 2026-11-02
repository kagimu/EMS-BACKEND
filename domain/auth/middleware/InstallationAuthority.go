package middleware

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

// RequireInstallationAuthority protects installation-wide authority. A
// deployment-configured software owner is always allowed. Otherwise the actor
// must hold the requested permission explicitly through a role; the ordinary
// Admin user-type pass-through is deliberately NOT applied here.
//
// This middleware never requires an active organisation: owner and
// installation-wide operations must work independently of a selected client.
func RequireInstallationAuthority(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}
		user, ok := value.(models.User)
		if !ok || user.Number == "" {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}

		if commons.IsSoftwareOwner(user.Number) {
			c.Next()
			return
		}

		if permission == "" {
			c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
			c.Abort()
			return
		}

		permissions, err := repository.FetchAllPermissionsByUser(user.Number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to check permissions"))
			c.Abort()
			return
		}
		for _, granted := range permissions {
			if granted == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
		c.Abort()
	}
}

// RequireSoftwareOwner restricts a route to the deployment-configured software
// owner. It deliberately does not accept a role permission or the Admin
// user-type bypass: owner-console access is an explicit deployment grant.
func RequireSoftwareOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}
		user, ok := value.(models.User)
		if !ok || user.Number == "" {
			c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
			c.Abort()
			return
		}
		if !commons.IsSoftwareOwner(user.Number) {
			c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access the owner console"))
			c.Abort()
			return
		}
		c.Next()
	}
}
