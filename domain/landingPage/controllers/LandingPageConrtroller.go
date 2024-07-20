package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLandingPage(c *gin.Context) {
	landingPageService := service.NewLandingPageService()
	page, err := landingPageService.GetLandingPage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}

func CreateLandingPage(c *gin.Context) {
	var page models.LandingPage
	var pageRequest requests.LandingPageRequest
	landingPageService := service.NewLandingPageService()

	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := pageRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	err := landingPageService.CreateLandingPage(&page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, page)
}
