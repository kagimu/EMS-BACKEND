package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/gin-gonic/gin"
)

func CreateApplication(c *gin.Context) {
	var application models.Application
	var err error

	err = c.ShouldBindJSON(&application)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server errror"})
		return
	}

	err = config.DB.Create(&application).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "application created successfully",
		"data":    application,
	})
}
