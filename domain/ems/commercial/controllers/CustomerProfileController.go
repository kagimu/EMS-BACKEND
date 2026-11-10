package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	commercialRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/requests"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/services"
	"github.com/gin-gonic/gin"
)

func writeProfileError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, commercialServices.ErrProfileNotFound), errors.Is(err, commercialServices.ErrProfileOrganisation):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, commercialServices.ErrProfileExists):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, commercialServices.ErrProfileValidation), errors.Is(err, commercialServices.ErrProfileManager):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		if stringsContainsInvalid(err) {
			c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, commons.Feedback("The customer profile request could not be completed"))
	}
}

func stringsContainsInvalid(err error) bool {
	return err != nil && (containsText(err.Error(), "invalid commercial status") || containsText(err.Error(), "invalid customer category"))
}

func containsText(value, target string) bool {
	return len(value) >= len(target) && (value == target || len(value) > len(target) && (indexOf(value, target) >= 0))
}

func indexOf(value, target string) int {
	for index := 0; index+len(target) <= len(value); index++ {
		if value[index:index+len(target)] == target {
			return index
		}
	}
	return -1
}

func actorNumber(c *gin.Context) (string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil || user.Number == "" {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	return user.Number, true
}

func GetCustomerProfile(c *gin.Context) {
	profile, err := commercialServices.GetProfile(c.Param("number"))
	if err != nil {
		writeProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer profile fetched", "data": profile})
}

func CreateCustomerProfile(c *gin.Context) {
	actor, ok := actorNumber(c)
	if !ok {
		return
	}
	var request commercialRequests.CustomerProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	profile, err := commercialServices.CreateProfile(actor, c.Param("number"), request)
	if err != nil {
		writeProfileError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Customer profile created", "data": profile})
}

func UpdateCustomerProfile(c *gin.Context) {
	actor, ok := actorNumber(c)
	if !ok {
		return
	}
	var request commercialRequests.CustomerProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	profile, err := commercialServices.UpdateProfile(actor, c.Param("number"), request)
	if err != nil {
		writeProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer profile updated", "data": profile})
}
