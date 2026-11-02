package middleware

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	"github.com/gin-gonic/gin"
)

const supportRoleKey = "support_role"
const supportUserKey = "support_user_number"

// RequireSupportStaff allows the deployment software owner and active support
// staff. It never grants ordinary client Administrators support access.
func RequireSupportStaff() gin.HandlerFunc {
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

		if commons.IsSoftwareOwner(user.Number) {
			c.Set(supportRoleKey, "OWNER")
			c.Set(supportUserKey, user.Number)
			c.Next()
			return
		}

		var staff supportModels.Staff
		if err := config.DB.Where("user_number = ? AND status = ?", user.Number, supportModels.SupportStatusActive).First(&staff).Error; err != nil {
			c.JSON(http.StatusForbidden, commons.Feedback("You do not have support workspace access"))
			c.Abort()
			return
		}
		c.Set(supportRoleKey, staff.Role)
		c.Set(supportUserKey, user.Number)
		c.Next()
	}
}

func SupportRole(c *gin.Context) string {
	value, _ := c.Get(supportRoleKey)
	role, _ := value.(string)
	return role
}

func SupportUser(c *gin.Context) string {
	value, _ := c.Get(supportUserKey)
	user, _ := value.(string)
	return user
}
