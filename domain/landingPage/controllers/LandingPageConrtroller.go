package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var landingPageService = service.NewLandingPageService()

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

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	pageTitle := c.PostForm("title")
	if pageTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide the page title"})
		return
	}

	err = landingPageService.LandingPageDetails(c, pageTitle, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Landing page details successfully created"})
}
