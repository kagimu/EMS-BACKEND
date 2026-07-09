package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/services"
	"github.com/gin-gonic/gin"
)

func SaveStaffScope(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.StaffScopeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input - " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := services.SaveStaffScope(&payload, user.Number); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save staff scope: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff scope has been successfully updated"})
}

func GetStaffScope(c *gin.Context) {
	staffNumber := c.Param("staff_number")
	if staffNumber == "" {
		staffNumber = c.Query("staff_number")
	}
	if staffNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "staff_number is required"})
		return
	}

	schoolNumber := c.Query("school_number")

	scopeTree, err := services.GetStaffScopeTree(staffNumber, schoolNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve staff scope: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scopeTree})
}

func UpdateStaffScope(c *gin.Context) {
	SaveStaffScope(c)
}

func DeleteStaffScope(c *gin.Context) {
	staffNumber := c.Param("staff_number")
	if staffNumber == "" {
		staffNumber = c.Query("staff_number")
	}
	if staffNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "staff_number is required"})
		return
	}

	schoolNumber := c.Query("school_number")

	if err := services.ClearStaffScope(staffNumber, schoolNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to clear staff scope: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff scope has been successfully cleared"})
}

