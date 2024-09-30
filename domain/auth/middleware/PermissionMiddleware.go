package middleware

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidatePermission(permission string) gin.HandlerFunc {

	return func(c *gin.Context) {
		var user models.User
		user = c.MustGet("user").(models.User)
		if user.Type == "Admin" {
			c.Next()
			return
		}

		if permission != "" {

			var allowed = false
			permissions, err := repository.FetchAllPermissionsByUser(c.GetHeader("User-Number"))

			if err != nil {
				c.JSON(http.StatusForbidden, commons.Feedback(err.Error()))
				c.Abort()
				return
			}

			// Check if permission is in the list of permissions
			for _, p := range permissions {
				if p == permission {
					allowed = true
					c.Next()
					break
				}
			}

			if !allowed {
				c.JSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
				c.Abort()
				return
			} else {
				c.Next()
			}

		}

	}

}
