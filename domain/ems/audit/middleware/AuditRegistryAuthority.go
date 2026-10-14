package middleware

import (
	"net/http"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

// ManageAuditRegistryPermission is the explicit control-plane authority for
// cross-organisation audit access (Foundation G §8/§23). It is deliberately
// distinct from the client-plane view_ems_audit permission.
const ManageAuditRegistryPermission = "manage_ems_audit"

// RequireAuditRegistryAuthority enforces the explicit cross-organisation audit
// authority. Unlike middleware.ValidatePermission it deliberately does NOT
// grant the Admin user-type pass-through: cross-organisation audit access
// requires an explicit authority, never a generic Admin bypass (Foundation G
// §8/§20/§9). No second identity system is introduced — the existing
// user-role-permission infrastructure is used.
func RequireAuditRegistryAuthority() gin.HandlerFunc {
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
			if permission == ManageAuditRegistryPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
		c.Abort()
	}
}
