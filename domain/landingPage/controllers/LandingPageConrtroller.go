package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
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
	landingPageService := service.NewLandingPageService()

	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := landingPageService.CreateLandingPage(&page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, page)
}
