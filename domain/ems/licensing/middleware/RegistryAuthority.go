package middleware

import (
	"net/http"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

// ManageLicenceRegistryPermission is the explicit control-plane authority for
// the licensing registry (Foundation F, §15).
const ManageLicenceRegistryPermission = "manage_licence_registry"

// RequireLicenceRegistryAuthority enforces the explicit licensing-management
// permission for the licensing control plane.
//
// Unlike middleware.ValidatePermission it deliberately does NOT grant the
// Admin user-type pass-through: licensing registry access is a narrow control
// plane authority and must be granted explicitly. An Admin without the
// permission is rejected. No second identity system is introduced — the
// existing user-role-permission infrastructure is used.
func RequireLicenceRegistryAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		permissions, err := authRepositories.FetchAllPermissionsByUser(user.Number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to check permissions"))
			c.Abort()
			return
		}
		for _, permission := range permissions {
			if permission == ManageLicenceRegistryPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
		c.Abort()
	}
}
