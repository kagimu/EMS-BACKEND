package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"github.com/gin-gonic/gin"
)

func GetCreditsByUser(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := c.Query("user_number")
	if userNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_number is required"})
		return
	}

	data, err := services.GetCreditsByUser(userNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Failed to get credits: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
