package middleware

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidatePermission(permission string) gin.HandlerFunc {

	return func(c *gin.Context) {
		if permission != "" {

			var allowed = false
			permissions, err := repository.FetchAllPermissionsByUser(c.GetHeader("User-Number"))

			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, commons.Feedback(fmt.Sprintf("%s", err)))
			}

			// Check if permission is in the list of permissions
			for _, p := range permissions {
				if p == permission {
					allowed = true
					break
				}
			}

			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, commons.Feedback("You don't have permission to access this resource"))
				return
			}

			c.Next()

		}

	}

}
