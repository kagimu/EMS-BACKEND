package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateSubjectGroupHandler handles POST /subject-groups
func CreateSubjectGroupHandler(c *gin.Context) {
	var req requests.CreateSubjectGroupRequest

	// Bind incoming JSON to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate struct
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err})
		return
	}

	// Call service to handle saving logic
	if err := services.CreateSubjectGroup(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subject group", "details": err})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{"message": "Subject group created successfully"})
}
