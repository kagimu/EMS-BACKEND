package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
)

func CreateParent(c *gin.Context) {
	var parent models.Parent
	var err error

	err = c.ShouldBindJSON(&parent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server errror"})
		return
	}

	parent.Prepare()

	err = parent.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input validation failed"})
		return
	}

	err = checkmail.ValidateFormat(parent.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is invalid"})
		return
	}

	err = config.DB.Create(&parent).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "parent created successfully",
		"data":    parent,
	})
}

func FetchParentInfo(c *gin.Context) {
	parent := models.Parent{}
	parentID := c.Param("id")

	err := config.DB.First(&parent, parentID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no parent found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": parent})
}

func FetchParents(c *gin.Context) {
	var parents []models.Parent

	err := config.DB.Find(&parents).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": parents})
}

func UpdateParent(c *gin.Context) {
	parent := models.Parent{}
	parentID := c.Param("id")

	err := config.DB.First(&parent, parentID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no parent found",
		})
		return
	}

	if err := c.BindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid request",
		})
		return
	}

	if err := config.DB.Save(&parent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to update parent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": parent})
}

func DeleteParent(c *gin.Context) {
	parent := models.Parent{}
	parentID := c.Param("id")

	err := config.DB.First(&parent, parentID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no parent found",
		})
		return
	}

	config.DB.Delete(&parent)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "parent deleted successfully",
	})
}
