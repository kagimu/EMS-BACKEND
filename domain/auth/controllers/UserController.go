package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/hassPassTokenGen"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	uuid "github.com/satori/go.uuid"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func FetchUsers(c *gin.Context) {

	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func FetchUser(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateUser(c *gin.Context) {

	var user requests.StaffAccountRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := user.ValidateStaffAccountRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	tx := config.DB.Begin()
	var userObject = models.User{
		Number:     uuid.NewV4().String(),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Gender:     user.Gender,
		Type:       "Staff",
		Status:     "Active",
	}

	if err := tx.Create(&userObject).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	//after you need to save user email
	saveEmailContact := models.Contact{
		Contact:     user.Email,
		ContactType: "Email",
		UserNumber:  userObject.Number,
	}
	if err := tx.Create(&saveEmailContact).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	savePhoneContact := models.Contact{
		Contact:     user.Phone,
		ContactType: "Phone",
		UserNumber:  userObject.Number,
	}
	if err := tx.Create(&savePhoneContact).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// append roles

	var userRole permissionModels.UserRole

	for _, roleNumber := range user.RoleNumbers {
		if err := tx.Where("number = ?", roleNumber).First(&userRole.Role).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch role"})
			return
		}

		userRole.UserNumber = userObject.Number
		userRole.RoleNumber = roleNumber
		userRole.AddedBy = c.GetHeader("User-Number")

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to append role to user"})
			return

		}
	}

	// set password
	hashPassword, hashErr := hassPassTokenGen.Hash(user.Password)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to encrypt password"})
		return
	}

	var userPassword = models.UserPassword{
		UserNumber:   userObject.Number,
		UserPassword: hashPassword,
	}

	if err := tx.Create(&userPassword).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, commons.Feedback("User created successfully"))
}
