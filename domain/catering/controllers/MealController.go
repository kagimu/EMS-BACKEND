package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	staffModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateMeal(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MealRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var feedingSession models.FeedingSession
	if err := config.DB.Where("number = ?", payload.FeedingSessionNumber).First(&feedingSession).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session not found"})
		return
	}

	var menuItem models.MenuItem
	if err := config.DB.Where("number = ?", payload.MenuItemNumber).First(&menuItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Menu item not found"})
		return
	}

	var staff staffModels.Staff
	if err := config.DB.Where("number = ?", payload.StaffNumber).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Staff member not found"})
		return
	}

	var user authModels.User
	if err := config.DB.Where("number = ?", payload.UserNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User number not found"})
		return
	}

	var stocking models.Stocking
	if err := config.DB.Where("menu_item_number = ?", payload.MenuItemNumber).First(&stocking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stocking not found"})
		return
	}

	updateStock := models.Stocking{
		MenuItemNumber: payload.MenuItemNumber,
		Quantity:       stocking.Quantity + -payload.Quantity,
	}

	if err := config.DB.Model(&stocking).Where("number = ?", stocking.Number).Updates(updateStock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update menu item stock"})
		return
	}

	saveMeal := models.Meal{
		Number:               uuid.NewV4().String(),
		FeedingSessionNumber: payload.FeedingSessionNumber,
		MenuItemNumber:       payload.MenuItemNumber,
		StaffNumber:          payload.StaffNumber,
		UserNumber:           payload.UserNumber,
		CampusNumber:         feedingSession.CampusNumber,
	}

	if err := config.DB.Create(&saveMeal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create meal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meal created successfully"})
}

func FetchMeals(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	feedingGroupNumber := c.Query("feeding_group_number")
	if feedingGroupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group number is required"})
		return
	}

	feedingSessionNumber := c.Query("feeding_session_number")
	if feedingSessionNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session number is required"})
		return
	}

	var feedingGroupUsers []models.FeedingGroup
	if err := config.DB.Where("number = ?", feedingGroupNumber).First(&feedingGroupUsers).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding group users"})
		return
	}

	var feedingSessionItems []models.FeedingSessionItem
	if err := config.DB.Where("feeding_session_number = ?", feedingSessionNumber).Preload("MenuItem").First(&feedingSessionItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding session items"})
		return
	}

	var feedingGroupItems []interface{}
	for _, v := range feedingGroupUsers {
		feedingGroupItems = append(feedingGroupItems, v)
	}

	for _, v := range feedingSessionItems {
		feedingGroupItems = append(feedingGroupItems, v)
	}

	var meals []models.Meal
	if err := config.DB.Where("feeding_session_number = ?", feedingSessionNumber).Find(&meals).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch meals"})
		return
	}

	var servedItems []models.FeedingSessionItem
	for _, item := range feedingSessionItems {
		served := false

		for _, meal := range meals {
			if meal.MenuItemNumber == item.MenuItemNumber {
				served = true
				break
			}
		}

		if served {
			servedItems = append(servedItems, item)
		}
	}

	// Feeding group number, feeding session  number as params
	// Query from feeding group where feeding group number = returns users
	//  Iterate
	// Query from feeding session items where feeding session number is given  = returns a list of menu items
	// Meal - attach menu items (result is the meals a person is supposed to eat)
	// For each menu item, confirm if the item has been served on this user
	// if the menu item exists in the meal table, pick the qunaity and supply it, else prefix with a zero
	// return list of users

	c.JSON(http.StatusOK, gin.H{"meals": servedItems, "menuitems": feedingGroupItems})
}
