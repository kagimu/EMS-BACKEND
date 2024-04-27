package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
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

	var meals []models.Meal
	if err := config.DB.Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch meals"})
		return
	}

	var response []responses.MealResponse
	for _, meal := range meals {
		var userMeals []responses.UserMeals

		response = append(response, responses.MealResponse{
			Number:               meal.Number,
			FeedingSessionNumber: meal.FeedingSessionNumber,
			FeedingSessionName:   meal.FeedingSession.Name,
			FeedingSessionCode:   meal.FeedingSession.Code,
			FeedingSessionStatus: meal.FeedingSession.Status,
			StockingNumber:       meal.StockingNumber,
			Date:                 meal.Date,
			Time:                 meal.Time,
			MenuItemNumber:       meal.MenuItemNumber,
			StaffIDNumber:        meal.StaffNumber,
			UserMeals:            userMeals,
			JobTitle:             meal.Staff.JobTitle.Name,
		})

		for _, meal := range meals {
			userMeals = append(userMeals, responses.UserMeals{
				UserName:     meal.User.FirstName + " " + meal.User.SecondName + " " + meal.User.LastName,
				MenuItemName: meal.MenuItem.Name,
				MenuItemType: meal.MenuItem.Type,
				MenuItemCode: meal.MenuItem.Code,
			})
		}

	}

	c.JSON(http.StatusOK, gin.H{"meals": response})
}
