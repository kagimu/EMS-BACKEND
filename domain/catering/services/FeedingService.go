package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

func AddMealCoursesToFeedingGroup(courseToFeedingGroupRequest requests.MealCourseToFeedingGroupRequest, currentUser string) error {

	feedingGroup, err := repositories.GetFeedingGroupByNumber(courseToFeedingGroupRequest.FeedingGroupNumber)

	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	if feedingGroup.Status != constants.StatusActive {
		return errors.New("selected group is not active")
	}

	for _, day := range courseToFeedingGroupRequest.Days {

		if checkErr := commons.CheckDay(day); checkErr != "" {
			return errors.New(checkErr)
		}

		if err = saveMealCoursesForDay(courseToFeedingGroupRequest.MealCourses, day, feedingGroup, currentUser); err != nil {
			return err
		}

	}

	return nil
}

func saveMealCoursesForDay(mealCourses []string, day string, feedingGroup models.FeedingGroup, currentUser string) error {

	for _, mealCourse := range mealCourses {
		checkMeal, err := repositories.GetMealCourseByNumber(mealCourse)
		if err != nil {
			return err
		}
		if checkMeal.Number == "" {
			return errors.New(fmt.Sprintf("invalid '%s' meal ID", day))
		}

		if checkMeal.Status != constants.StatusActive {
			return errors.New(fmt.Sprintf("meal course '%s' is not active", checkMeal.Name))
		}

		groupMeal, err := repositories.GetFeedingGroupMealCourseByMealAndGroup(checkMeal.Number, day, feedingGroup.Number)
		if err != nil {
			return err
		}

		if groupMeal.Number != "" {
			return errors.New(fmt.Sprintf("meal course '%s' already added to group '%s' for '%s'", checkMeal.Name, feedingGroup.Number, day))
		}

		newGroupMealCourse := models.FeedingGroupMealCourse{
			Number:             uuid.NewV4().String(),
			FeedingGroupNumber: feedingGroup.Number,
			MealCourseNumber:   mealCourse,
			Day:                day,
			UserNumber:         currentUser,
		}

		if err := repositories.SaveFeedingGroupMealCourse(newGroupMealCourse); err != nil {
			return err
		}

	}

	return nil
}

func GetFeedingGroupMealCourseCalendar(schoolNumber string) (responses.FeedingGroupMealCourseCalendar, error) {

	var response responses.FeedingGroupMealCourseCalendar

	school, err := repositories2.GetSchoolByNumber(schoolNumber)

	if err != nil {
		return response, err
	}

	if school.Number == "" {
		return response, errors.New("invalid school number")
	}

	headers := []string{" "}

	for _, day := range commons.DaysOfTheWeek() {
		headers = append(headers, day)
	}

	groups, err := repositories.GetAllFeedingGroups(school.Number, constants.StatusActive)

	if err != nil {
		return response, err
	}

	var items []responses.GroupCalendarRawResponse

	for _, group := range groups {

		dailyMenu := getGroupDailyMenu(group.Number)

		item := responses.GroupCalendarRawResponse{
			GroupName:   group.Name,
			MealCourses: dailyMenu,
		}

		items = append(items, item)

	}

	response.FeedingGroup = items

	response.Headers = headers

	return response, nil
}

func getGroupDailyMenu(groupNumber string) []responses.GroupDayMealCourses {

	var weeklyItems []responses.GroupDayMealCourses

	for _, day := range commons.DaysOfTheWeek() {

		var item responses.GroupDayMealCourses
		item.Day = day
		item.Active = false

		mealCourses, err := repositories.GetAllFeedingGroupMealCoursesForDay(groupNumber, day)

		var meals []string

		if err != nil {
			item.Error = err.Error()
		} else {

			for _, _mealCourse := range mealCourses {
				mealCourse := _mealCourse.MealCourse
				if mealCourse.Status == constants.StatusActive {
					meals = append(meals, mealCourse.Name)
				}
			}

			if len(mealCourses) > 0 && len(meals) == 0 {
				item.Error = "Inactive meal courses"
			} else {
				item.Error = "No meal course(s) set"
			}

		}

		item.MealCourses = meals

		if len(meals) > 0 {
			item.Active = true
		}

		weeklyItems = append(weeklyItems, item)

	}

	return weeklyItems
}

func GetUserMeals(mealCourseNumber string, campusNumber string, feedingGroupNumber string) ([]responses.UserMealResponse, error) {

	var mealsResponse []responses.UserMealResponse

	currentTime := time.Now()
	day := currentTime.Format("Mon")
	checkDay := commons.CheckDay(day)
	if checkDay != "" {
		return mealsResponse, errors.New(checkDay)
	}

	campus, err := repositories3.GetCampusByNumber(campusNumber)

	if err != nil {
		return mealsResponse, err
	}

	if campus.Number == "" {
		return mealsResponse, errors.New("invalid campus number")
	}

	mealCourse, err := repositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourse.Status != constants.StatusActive {
		return mealsResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}

	checkSchool, err := repositories.GetMealCourseByNumberAndSchool(mealCourse.Number, campus.SchoolNumber)
	if err != nil {
		return mealsResponse, err
	}

	if checkSchool.Number == "" {
		return mealsResponse, errors.New("meal course can not be used on the selected campus")
	}

	//GetFeedingSessionForCampusAndMealCourseAndDay

	mealCourseFeedingSession, err := repositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourseFeedingSession.Number == "" {
		return mealsResponse, errors.New("menu not set up for this meal course on '" + day + "'")
	}

	feedingSession := mealCourseFeedingSession.FeedingSession

	feedItems, err := repositories.GetFeedingSessionItems(feedingSession.Number)
	if err != nil {
		return mealsResponse, err
	}

	if len(feedItems) == 0 {
		return mealsResponse, errors.New("no items set to '" + feedingSession.Number + "' session")
	}

	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {

		checkGroup, err := repositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkGroup.Number == "" || checkGroup.SchoolNumber != campus.SchoolNumber {
			return mealsResponse, errors.New("invalid group number")
		}

		checkAllowedToHaveMeal, err := repositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkAllowedToHaveMeal.Number == "" {
			return mealsResponse, errors.New("this group should not be served today")
		}

		feedingGroups = append(feedingGroups, checkGroup)

	} else {

		allFeedingGroupsForMealAndDay, err := repositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return mealsResponse, err
		}

		for _, feedingGroup := range allFeedingGroupsForMealAndDay {
			feedingGroups = append(feedingGroups, feedingGroup.FeedingGroup)
		}

	}

	for _, feedingGroup := range feedingGroups {

		allUsers, err := repositories.GetFeedingGroupUserByGroupNumber(feedingGroup.Number)
		if err != nil {
			return mealsResponse, err
		}

		for _, _user := range allUsers {

			user := _user.User

			userMeal := responses.UserMealResponse{
				UserNumber:           user.Number,
				Name:                 user.FirstName,
				FeedingGroupNumber:   feedingGroup.Number,
				FeedingGroupName:     feedingGroup.Name,
				FeedingSessionNumber: feedingSession.Number,
			}

			userMeal.Meals = findUserMeal(user.Number, feedItems, mealCourse.Number)

			mealsResponse = append(mealsResponse, userMeal)
		}

	}

	return mealsResponse, nil
}

func findUserMeal(userNumber string, feedingSessionItems []models.FeedingSessionItem, mealCourse string) []responses.SingleUserMealItem {

	var items []responses.SingleUserMealItem

	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	date_, err := time.Parse(time.DateOnly, currentDate)
	if err != nil {
		fmt.Println("Error formatting today's date: " + currentDate)
	}

	for _, feedingSessionItem := range feedingSessionItems {

		menuItem, err := repositories.GetMenuItemByNumber(feedingSessionItem.MenuItemNumber)
		if err != nil {
			fmt.Println("Error getting menu item from feeding session")
			return items
		}

		var served float64

		servedStock, err := repositories.GetMealStockingByMealNumberAndDateUserNumber(mealCourse, date_, userNumber, menuItem.Number)

		if err != nil {
			fmt.Println("Error getting menu item served quantity: " + err.Error())
			return items
		}

		if servedStock.Number != "" {
			served = -1 * servedStock.Quantity
		} else {
			served = 0
		}

		singleUserMenu := responses.SingleUserMealItem{
			Quantity:       feedingSessionItem.Quantity,
			ServedQuantity: served,
			MenuItemNumber: menuItem.Number,
			MenuItemName:   menuItem.Name,
		}

		items = append(items, singleUserMenu)
	}

	return items
}

func GetFeedingGroupsBySchoolNumber(schoolNumber string, status string) ([]responses.FeedingGroupResponse, error) {

	var response []responses.FeedingGroupResponse

	feedingGroups, err := repositories.GetAllFeedingGroups(schoolNumber, status)
	if err != nil {
		return response, err
	}

	for _, feedingGroup := range feedingGroups {

		user, err := repository.GetUserByNumber(feedingGroup.UserNumber)
		if err != nil {
			return response, err
		}

		if user.Number == "" {
			return response, errors.New("user not found")
		}

		school, err := repositories2.GetSchoolByNumber(feedingGroup.SchoolNumber)
		if err != nil {
			return response, err
		}

		if school.Number == "" {
			return response, errors.New("school not found")
		}

		response = append(response, responses.FeedingGroupResponse{
			Number:       feedingGroup.Number,
			Name:         feedingGroup.Name,
			UserName:     user.FirstName + " " + user.LastName,
			Status:       feedingGroup.Status,
			SchoolNumber: school.Number,
			SchoolName:   school.Name,
		})

	}

	return response, nil
}

func AddMeal(mealRequest requests.MealRequest, currentUser string) error {

	currentStock, err := repositories.GetMenuItemCurrentStock(mealRequest.MenuItemNumber)
	if err != nil {
		return err
	}

	if currentStock < mealRequest.Quantity {
		return errors.New("not enough stock")
	}

	tx := config.DB.Begin()

	formatDate, err := time.Parse(time.DateOnly, mealRequest.Date)
	if err != nil {
		return err
	}

	day := formatDate.Format("Mon")
	checkDay := commons.CheckDay(day)
	if checkDay != "" {
		return errors.New(checkDay)
	}

	newMeal, err := repositories.GetMealByMealNumberAndDateAndUserNumber(mealRequest.MealCourseNumber, formatDate, mealRequest.UserNumber)
	if err != nil {
		return errors.New("Checking if item served error: " + err.Error())
	}

	groupUser, err := repositories.GetFeedingGroupUserByUserNumber(mealRequest.UserNumber)
	if err != nil {
		return err
	}

	if groupUser.Number == "" {
		return errors.New("user not added to any group")
	}

	if newMeal.Number == "" {
		newMeal.Number = uuid.NewV4().String()
		newMeal.MealCourseNumber = mealRequest.MealCourseNumber
		newMeal.FeedingSessionNumber = mealRequest.FeedingSessionNumber
		newMeal.UserNumber = mealRequest.UserNumber
		newMeal.Date = formatDate
		newMeal.Time = mealRequest.Time
		newMeal.AddedByNumber = currentUser
		newMeal.StaffNumber = currentUser
		newMeal.Day = day
		newMeal.CampusNumber = groupUser.CampusNumber
	} else {
		newMeal.StaffNumber = currentUser
		newMeal.CampusNumber = groupUser.CampusNumber
		newMeal.FeedingSessionNumber = mealRequest.FeedingSessionNumber
		newMeal.FeedingGroupNumber = groupUser.FeedingGroupNumber
	}

	err = repositories.SaveMeal(tx, newMeal)
	if err != nil {
		return errors.New("error creating meal: " + err.Error())
	}

	stocking, err := repositories.GetMealStockByMealNumberAndMenuItemNumber(tx, newMeal.Number, mealRequest.MenuItemNumber)
	if err != nil {
		return err
	}

	if stocking.Number == "" {
		stocking.Number = uuid.NewV4().String()
		stocking.MenuItemNumber = mealRequest.MenuItemNumber
		stocking.Quantity = mealRequest.Quantity * -1
		stocking.CampusNumber = groupUser.CampusNumber
		stocking.UserNumber = currentUser
	} else {
		stocking.Quantity = mealRequest.Quantity * -1
		stocking.CampusNumber = groupUser.CampusNumber
		stocking.UserNumber = currentUser
	}

	err = repositories.SaveStocking(tx, stocking)
	if err != nil {
		return errors.New("error reducing stock: " + err.Error())
	}

	mealStocking, err := repositories.GetMealStockingByMealNumberAndStockingNumber(tx, newMeal.Number, stocking.Number)
	if err != nil {
		return err
	}

	if mealStocking.MealNumber == "" {
		mealStocking.MealNumber = newMeal.Number
		mealStocking.StockingNumber = stocking.Number
	}

	err = repositories.SaveMealStocking(tx, mealStocking)
	if err != nil {
		return errors.New("error mapping stock with meal taken: " + err.Error())
	}

	if err = tx.Commit().Error; err != nil {
		return errors.New("Commit error: " + err.Error())
	}

	return nil
}
