package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateParent(c *gin.Context) {
	var parent requests.ParentRequest

	if err := c.ShouldBindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server errror"})
		return
	}

	parent.Prepare()

	if err := parent.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input validation failed"})
		return
	}

	if err := checkmail.ValidateFormat(parent.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is invalid"})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	var user authModel.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	tx := config.DB.Begin()
	saveApplication := applicationModel.Application{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
	}

	if err := tx.Create(&saveApplication).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to create application",
			"details": err.Error(),
		})
		return
	}

	saveParent := models.Parent{
		Number:                     uuid.NewV4().String(),
		Surname:                    parent.Surname,
		SecondName:                 parent.SecondName,
		OtherName:                  parent.OtherName,
		MainTelephoneNumber:        parent.MainTelephoneNumber,
		AlternativeTelephoneNumber: parent.AlternativeTelephoneNumber,
		Email:                      parent.Email,
		AlternativeEmail:           parent.AlternativeEmail,
		UserNumber:                 userNumber,
		ApplicationNumber:          saveApplication.Number,
	}

	if parent.AlternativeEmail == parent.Email {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and alternative email should not be the same "})
		return
	}

	if parent.AlternativeTelephoneNumber == parent.MainTelephoneNumber {
		c.JSON(http.StatusBadRequest, gin.H{"message": "telephone and alternative telephone should not be the same "})
		return
	}

	if err := tx.Create(&saveParent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to save parent",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to save parent"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "parent information created successfully",
		"data":    saveParent,
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
