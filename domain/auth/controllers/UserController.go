package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/requests"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	uuid "github.com/satori/go.uuid"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func FetchUsers(c *gin.Context) {

	entityNumber, _ := commons.GetUserEntityNumber(c)

	var paginatedResponse = pagination.PaginatedResponse{}
	var users []models.User
	if err := config.DB.Scopes(
		filters.ByStatus(c.Query("status")),
		pagination.Paginate(c, &paginatedResponse),
		filters.ByAnyField("entity_users.entity_number", entityNumber),
	).Joins("JOIN entity_users ON entity_users.user_number = users.number").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	type roleObject struct {
		RoleName   string `json:"role_name"`
		RoleNumber string `json:"role_number"`
		ExpiryDate string `json:"expiry_date"`
	}

	type CLeanUser struct {
		models.User
		Roles []roleObject `json:"roles"`
	}

	var cleanUsers = make([]CLeanUser, 0)

	for _, user := range users {
		var userRoles []permissionModels.UserRole

		var roles []roleObject
		if err := config.DB.Preload("Role").Where("user_number = ?", user.Number).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user roles"})
			return
		}

		for _, userRole := range userRoles {
			roles = append(roles, roleObject{RoleName: userRole.Role.Name,
				RoleNumber: userRole.Role.Number, ExpiryDate: userRole.ExpiryDate.Format("2006-01-02")})
		}

		cleanUsers = append(cleanUsers, CLeanUser{User: user, Roles: roles})

	}

	paginatedResponse.Content = cleanUsers

	c.JSON(http.StatusOK, paginatedResponse)
}
func FetchAllUsers(c *gin.Context) {

	//var paginatedResponse = pagination.PaginatedResponse{}
	var users []models.User
	if err := config.DB.Scopes(
		filters.ByStatus(c.Query("status")),
	).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	type CLeanUser struct {
		models.User
		Roles []string `json:"roles"`
	}

	var cleanUsers = make([]CLeanUser, 0)

	for _, user := range users {
		var userRoles []permissionModels.UserRole
		var roles []string
		if err := config.DB.Preload("Role").Where("user_number = ?", user.Number).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user roles"})
			return
		}

		for _, userRole := range userRoles {
			roles = append(roles, userRole.Role.Name)
		}

		cleanUsers = append(cleanUsers, CLeanUser{User: user, Roles: roles})

	}

	c.JSON(http.StatusOK, cleanUsers)
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
		Type:       constants.TypeStaff,
		Status:     "Active",
	}

	if err := tx.Create(&userObject).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// check if email or phone already exists

	if err := tx.Where("contact = ?", user.Email).First(&models.Contact{}).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "This Email was already used to create an account"})
		return
	}

	if err := tx.Where("contact = ?", user.Phone).First(&models.Contact{}).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "This Phone was already used to create an account"})
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

	for _, role := range user.Roles {
		var userRole permissionModels.UserRole
		var roleModel permissionModels.Role
		if err := tx.Where("number = ?", role.RoleNumber).First(&roleModel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch role"})
			return
		}

		userRole.UserNumber = userObject.Number
		userRole.RoleNumber = role.RoleNumber
		if role.ExpiryDate == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Expiry date is required"})
			return

		}
		date, err := commons.ConvertStringToDate(role.ExpiryDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to append role to user" + err.Error()})
			return
		}
		userRole.ExpiryDate = date
		userRole.AddedBy = c.GetHeader("User-Number")

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create role to user"})
			return

		}
	}

	// set password
	//hashPassword, hashErr := hassPassTokenGen.Hash(user.Password)
	//if hashErr != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to encrypt password"})
	//	return
	//}
	//
	//var userPassword = models.UserPassword{
	//	UserNumber:   userObject.Number,
	//	UserPassword: hashPassword,
	//}
	//
	//if err := tx.Create(&userPassword).Error; err != nil {
	//	tx.Rollback()
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user" + err.Error()})
	//	return
	//}

	// attach user to this entity
	if user.EntityNumber != "" {
		var entityUser = eosModels.EntityUser{
			EntityNumber:  user.EntityNumber,
			UserNumber:    userObject.Number,
			AddedByNumber: c.GetHeader("User-Number"),
		}

		if err := tx.Create(&entityUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user" + err.Error()})
			return
		}

	}

	//create / generate a token for the user to reset their password
	//userToken := tokenGenerator()
	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userObject.Number,
		Token:      userToken,
		Status:     "Inactive",
	}

	if err := tx.Create(&saveUserPasswordToken).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to create token"})
		return
	}
	//email sender
	Token := userToken
	d := emailSender.Info{Email: user.Email, Token: Token, MailType: "WholesomeEd: User Account Confirmation", UserName: userObject.SecondName, CurrentYear: time.Now().Year(), Type: "registration"}
	err := d.SendEmail()
	if err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//fmt.Println(fmt.Printf("Sending mail error %s\n", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to send email\n" + err.Error()})
		return
		//fmt.Printf("Sending mail error %s\n", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Failed :"+err.Error()))
	}

	c.JSON(http.StatusCreated, commons.Feedback("User created successfully"))
}

func CreateUserAcc(c *gin.Context) {
	userRepo := repository.NewSqlUserRepository(config.DB)
	fmt.Println(userRepo)
	//userService := services.NewUserService(userRepo)
}
