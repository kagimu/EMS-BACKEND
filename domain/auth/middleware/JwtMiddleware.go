package middleware

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	"github.com/gin-gonic/gin"
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
				"message": "Bearer not set",
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

			var userData models.User
			//check if user is a valid user
			if err := config.DB.Where("number = ?", claims.UserNumber).First(&userData).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
				return
			}
			c.Next()

		} else {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
		}

	}
}
