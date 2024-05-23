package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveMealCourse(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MealCourseCountRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = services.SaveMealCourse(payload, user.Number); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ERROR: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal course successfully saved"})
}

func UpdateSaveMealCourse(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateMealCourseCountRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = services.UpdateSaveMealCourse(payload, user.Number); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Updating ERROR: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal course successfully saved"})
}

func GetMealCourses(ctx *gin.Context) {

	_, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := ctx.Query("school_number")
	if schoolNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "provide school number"})
		return
	}

	response, err := services.GetMealCourses(schoolNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Getting meal courses error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func MealCourseCalendar(ctx *gin.Context) {

	_, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := ctx.Query("campus_number")
	if campusNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "provide campus number"})
		return
	}

	response, err := services.MealCourseCalendar(campusNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Getting Meal course calendar error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func FeedingSessionsToMealCourse(ctx *gin.Context) {

	_, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MealCourseFeedingSessionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.AddFeedingSessionsToMealCourse(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Adding feeding sessions error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Feeding sessions successfully added"))
}

func AllMealCourseFeedingSessions(ctx *gin.Context) {

	_, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	mealCourseNumber := ctx.Query("meal_course_number")
	if mealCourseNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "provide the meal course number"})
		return
	}

	response, err := services.MealCourseFeedingSessions(mealCourseNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Getting Meal course calendar error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func ActivateMealCourse(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ActivateMealCourseRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.ActivateMealCourse(payload, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Adding feeding sessions error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Feeding sessions successfully added"))

}

func FetUsersMeals(c *gin.Context) {

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := c.Query("campus_number")
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose campus"})
		return
	}

	mealCourseNumber := c.Query("meal_course_number")
	if mealCourseNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a meal course"})
		return
	}

	feedingGroupNumber := c.Query("feeding_group_number")

	usersMeal, err := services.GetUserMeals(mealCourseNumber, campusNumber, feedingGroupNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": usersMeal})
}

func AddMeal(c *gin.Context) {
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

	err = services.AddMeal(payload, loggedInUser.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meal created successfully"})
}

func UpdateMealCourseItem(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = services.UpdateMealCourseName(payload, user.Number); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ERROR: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal course name updated successfully"})
}
