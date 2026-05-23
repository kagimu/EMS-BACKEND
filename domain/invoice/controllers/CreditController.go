package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"github.com/gin-gonic/gin"
)

func CreateCredit(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.CreateCreditRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide a valid credit request: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	data, err := services.CreateCredit(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Failed to record credit: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit recorded successfully", "data": data})
}

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
