package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/services"
	"github.com/gin-gonic/gin"
)

func FetchSummary(c *gin.Context) {
	var filters requests.ReportsFilterRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	summary, err := services.FetchSummary(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch summary reports"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func FetchSchoolStructureReport(c *gin.Context) {
	var filters requests.ReportsFilterRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	report, err := services.FetchSchoolStructureReport(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch school structure report", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
