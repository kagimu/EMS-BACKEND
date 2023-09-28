package middleware

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization not set",
			})
			return
		}

		if len(c.Request.Header["User-Number"]) < 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Not Authenticated!",
			})
			return
		}

		tokenString := authHeader[len(BearerSchema):]

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Barer not set",
			})
			return
		}

		token, err := services.JWTAuthService().ValidateToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Sprintf("%s", err),
			})
			return
		}

		if token.Valid {
			claims := token.Claims.(*services.AuthCustomClaims)
			c.Request.Header.Add("User-Number", claims.UserNumber)
			c.Next()

		} else {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
		}

	}
}

func ValidatePermission(auth error, permission string) gin.HandlerFunc {

	return func(c *gin.Context) {
		if auth != nil && permission != "" {
			var err error
			var allowed = true
			//	allowed, err := auth.CheckPermission(c.GetHeader("User-Number"), permission)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": fmt.Sprintf("%s", err),
				})
			}

			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "You don't have permission to access this resource",
				})
				return
			} else {
				c.Next()
			}
		}

		/*c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})*/

	}

}
