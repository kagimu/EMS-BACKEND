package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	status := ctx.Query("status")

	response, err := services.GetMealCourses(schoolNumber, status)
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

func UpdateFeedingSessionsToMealCourse(ctx *gin.Context) {

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

	// AddFeedingSessionsToMealCourse already handles replacing existing ones
	err = services.AddFeedingSessionsToMealCourse(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Updating feeding sessions error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Feeding sessions successfully updated"))
}

func DeleteMealCourseFeedingSession(ctx *gin.Context) {

	_, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.DeleteMealCourseFeedingSessionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Check request body and try again"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.DeleteMealCourseFeedingSession(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Deleting feeding session error: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Feeding session successfully deleted"))
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

// FetchLearnersMeals mirrors FetUsersMeals but only returns enrolled learners (students
// with an active enrollment).  Each record includes full enrollment details just like
// the enrollment GET endpoint, and the list can be narrowed down with the optional filters:
//
//	campus_number          (required)
//	meal_course_number     (required)
//	feeding_group_number   (optional)
//	study_period_number    (optional)
//	academic_year_number   (optional)
//	campus_class_number    (optional)
//	stream_number          (optional)
func FetchLearnersMeals(c *gin.Context) {

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
	studyPeriodNumber := c.Query("study_period_number")
	academicYearNumber := c.Query("academic_year_number")
	campusClassNumber := c.Query("campus_class_number")
	streamNumber := c.Query("stream_number")
	classNumber := c.Query("class_number")

	learnersMeal, err := services.GetLearnerMeals(mealCourseNumber, campusNumber, feedingGroupNumber, studyPeriodNumber, academicYearNumber, campusClassNumber, streamNumber, classNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": learnersMeal})
}

// FetchLearnersMealsHistory returns historical meal delivery data for enrolled learners
// for a specific date, along with the same filtering capabilities as FetchLearnersMeals.
func FetchLearnersMealsHistory(c *gin.Context) {

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

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Date is required"})
		return
	}

	feedingGroupNumber := c.Query("feeding_group_number")
	studyPeriodNumber := c.Query("study_period_number")
	academicYearNumber := c.Query("academic_year_number")
	campusClassNumber := c.Query("campus_class_number")
	streamNumber := c.Query("stream_number")
	classNumber := c.Query("class_number")

	learnersMealHistory, err := services.GetLearnerMealsHistory(mealCourseNumber, campusNumber, feedingGroupNumber, studyPeriodNumber, academicYearNumber, campusClassNumber, streamNumber, classNumber, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": learnersMealHistory})
}

// FetchStaffMeals returns meal eligibility and served status for staff members
// that are assigned to feeding groups for the given school and meal course.
//
// Query params:
//
//	school_number        (required)
//	meal_course_number   (required)
//	feeding_group_number (optional)
func FetchStaffMeals(c *gin.Context) {

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a school"})
		return
	}

	mealCourseNumber := c.Query("meal_course_number")
	if mealCourseNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a meal course"})
		return
	}

	feedingGroupNumber := c.Query("feeding_group_number")

	staffMeals, err := services.GetStaffMeals(mealCourseNumber, schoolNumber, feedingGroupNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffMeals})
}

// FetchStaffMealsHistory returns historical meal delivery data for staff members
func FetchStaffMealsHistory(c *gin.Context) {

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a school"})
		return
	}

	mealCourseNumber := c.Query("meal_course_number")
	if mealCourseNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a meal course"})
		return
	}

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Date is required"})
		return
	}

	feedingGroupNumber := c.Query("feeding_group_number")

	staffMealsHistory, err := services.GetStaffMealsHistory(mealCourseNumber, schoolNumber, feedingGroupNumber, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffMealsHistory})
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

func AddStaffMeal(c *gin.Context) {
	loggedInUser, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StaffMealRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.AddStaffMeal(payload, loggedInUser.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff meal recorded successfully"})
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

func MealsProjections(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := context.Query("campus_number")
	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a campus"})
		return
	}

	_startDate := context.Query("start_date")
	if _startDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a start date"})
		return
	}

	_endDate := context.Query("end_date")
	if _endDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the end date"})
		return
	}

	var startDate time.Time
	if startDate, err = time.Parse(time.DateOnly, _startDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
		return
	}

	var endDate time.Time
	if endDate, err = time.Parse(time.DateOnly, _endDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
		return
	}

	if endDate.Before(startDate) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "End date: " + _endDate + " should be after the start date: " + _startDate})
		return
	}

	mealCourseNumber := context.Query("meal_course_number")

	response, err := services.MenuProjectedCalender(startDate, campusNumber, endDate, mealCourseNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Getting projections error: " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": response})
}
