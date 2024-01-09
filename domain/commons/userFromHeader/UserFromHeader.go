package userfromheader

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromHeader(c *gin.Context) (*models.User, error) {
	userNumber := c.Request.Header["User-Number"][0]
	var user models.User

	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
