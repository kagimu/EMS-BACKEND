package middleware

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	eosServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/services"
	"github.com/gin-gonic/gin"
)

// installationLevelPermissions are the catalogue permissions whose resource is
// installation-wide rather than owned by one organisation: the module catalogue
// and the installation package catalogue. They are resolved with the
// installation-wide lookup and do not require an active organisation. Every
// other permission is organisation-scoped.
var installationLevelPermissions = map[string]struct{}{
	"view_ems_modules":    {},
	"manage_ems_modules":  {},
	"manage_ems_packages": {},
}

// ValidatePermission enforces an organisation-scoped business permission.
//
// The organisation context is always resolved server-side from the
// authenticated user's active_user_entities row; a client-supplied organisation
// identifier is never consulted. If no active organisation exists the request
// fails with 409 instead of falling back to broader permissions.
//
// Installation-level Admin authority is bounded to that active organisation.
// The users.type == "Admin" flag preserves the existing installation
// administration workflow inside an organisation the user has selected, but it
// never grants business permissions in another organisation, and it grants
// nothing at all when no active organisation is set. Genuine
// cross-organisation control-plane operations use their explicit authorities
// (manage_licence_registry, manage_ems_audit), not this middleware.
func ValidatePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if permission == "" {
			c.Next()
			return
		}

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

		if _, installationLevel := installationLevelPermissions[permission]; installationLevel {
			if user.Type == "Admin" {
				c.Next()
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
			return
		}

		entityService := eosServices.EntityService{}
		activeEntity, err := entityService.GetCurrentEntity(user.Number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to resolve the active organisation"))
			c.Abort()
			return
		}
		if activeEntity.Number == "" {
			c.JSON(http.StatusConflict, commons.Feedback("No active organisation is selected"))
			c.Abort()
			return
		}
		c.Set("entity_number", activeEntity.Number)

		if user.Type == "Admin" {
			c.Next()
			return
		}

		permissions, err := repository.FetchPermissionsByUserAndEntity(user.Number, activeEntity.Number)
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
