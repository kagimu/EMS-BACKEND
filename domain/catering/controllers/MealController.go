package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateMeal(c *gin.Context) {
	loggedInUser, err := commons.GetUserFromHeader(c)
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

	//var staff staffModels.Staff
	//if err := config.DB.Where("number = ?", payload.StaffNumber).First(&staff).Error; err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "Staff member not found"})
	//	return
	//}

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

	tx := config.DB.Begin()

	// Checks if the menu items, user, date and feeding sessions are the same. Updates the meal details if true.
	var existingMeal models.Meal
	if err := config.DB.Where("feeding_session_number = ? AND menu_item_number = ? AND user_number = ? AND date = ?",
		payload.FeedingSessionNumber, payload.MenuItemNumber, payload.UserNumber, payload.Date).First(&existingMeal).Error; err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if existingMeal.Number != "" {
		// update the meal

		// set values to existing meal and save
		existingMeal.Quantity = payload.Quantity
		existingMeal.Time = payload.Time
		existingMeal.StaffNumber = loggedInUser.Number
		existingMeal.UserNumber = payload.UserNumber
		existingMeal.AddedByNumber = loggedInUser.Number
		existingMeal.StockingNumber = stocking.Number

		if err := tx.Save(&existingMeal).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update meal"})
			return
		}

		// update the stocking

		var existingStocking models.Stocking
		if err := tx.Where("number = ?", existingMeal.StockingNumber).First(&existingStocking).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stocking"})
			return
		}

		existingStocking.Quantity = 0 - payload.Quantity
		if err := tx.Save(&stocking).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update stocking"})
			return

		}

	} else {

		updateStock := models.Stocking{
			Number:         uuid.NewV4().String(),
			MenuItemNumber: payload.MenuItemNumber,
			Quantity:       0 - payload.Quantity,
			UserNumber:     loggedInUser.Number,
			CampusNumber:   feedingSession.CampusNumber,
			UnitCost:       0,
		}

		if err := tx.Create(&updateStock).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update menu item stock"})
			return
		}

		formatDate, err := time.Parse(time.RFC3339, payload.Date)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse date"})
			return
		}

		// Use the date part of formatDate
		formatDate = formatDate.UTC().Truncate(24 * time.Hour)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse time"})
			return
		}

		saveMeal := models.Meal{
			Number:               uuid.NewV4().String(),
			FeedingSessionNumber: payload.FeedingSessionNumber,
			MenuItemNumber:       payload.MenuItemNumber,
			StaffNumber:          loggedInUser.Number,
			Quantity:             payload.Quantity,
			UserNumber:           payload.StaffNumber,
			Date:                 formatDate,
			Time:                 payload.Time,
			CampusNumber:         feedingSession.CampusNumber,
			AddedByNumber:        loggedInUser.Number,
			StockingNumber:       updateStock.Number,
		}

		if err := tx.Create(&saveMeal).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create meal"})
			return
		}

	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create meal"})
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

	// Feeding group number, feeding session  number as params
	// Query from feeding group where feeding group number = returns users
	//  Iterate
	// Query from feeding session items where feeding session number is given  = returns a list of menu items
	// Meal - attach menu items (result is the meals a person is supposed to eat)
	// For each menu item, confirm if the item has been served on this user
	// if the menu item exists in the meal table, pick the qunaity and supply it, else prefix with a zero
	// return list of users
	type singleMeal struct {
		MenuItemNumber string  `json:"menu_item_number"`
		MenuItemName   string  `json:"menu_item_name"`
		Quantity       float64 `json:"quantity"`
	}
	type UserMeals struct {
		UserNumber         string       `json:"user_number"`
		Name               string       `json:"name"`
		FeedingGroupName   string       `json:"feeding_group_name"`
		FeedingGroupNumber string       `json:"feeding_group_number"`
		Meals              []singleMeal `json:"meals"`
	}

	var userMeals []UserMeals

	var feedingGroupUsers []models.FeedingGroupUser
	if err := config.DB.Preload("User").Preload("FeedingGroup").Where("feeding_group_number = ?", feedingGroupNumber).Find(&feedingGroupUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group does not have any users"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding group users"})
		return
	}

	for groupUserIndex, v := range feedingGroupUsers {
		var feedingSessionItems []models.FeedingSessionItem
		if err := config.DB.Where("feeding_session_number = ?", feedingSessionNumber).Preload("MenuItem").Find(&feedingSessionItems).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding session items"})
			return
		}

		userMeals = append(userMeals, UserMeals{
			UserNumber:         v.UserNumber,
			Name:               v.User.FirstName + " " + v.User.LastName,
			FeedingGroupName:   v.FeedingGroup.Name,
			FeedingGroupNumber: v.FeedingGroup.Number,
		})

		for _, item := range feedingSessionItems {

			var meal models.Meal
			if err := config.DB.Where("feeding_session_number = ? AND user_number = ? ", feedingSessionNumber, v.UserNumber).First(&meal).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					userMeals[groupUserIndex].Meals = append(userMeals[groupUserIndex].Meals, singleMeal{
						MenuItemNumber: item.MenuItemNumber,
						MenuItemName:   item.MenuItem.Name,
						Quantity:       0,
					})

				}
			} else {
				userMeals[groupUserIndex].Meals = append(userMeals[groupUserIndex].Meals, singleMeal{
					MenuItemNumber: item.MenuItemNumber,
					MenuItemName:   item.MenuItem.Name,
					Quantity:       meal.Quantity,
				})
			}

		}

	}

	c.JSON(http.StatusOK, gin.H{"data": userMeals})
}

func SaveMealCourse(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MealCourseCountRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = services.SaveMealCourse(payload, user.Number); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Meal course error: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal course successfully saved"})
}
