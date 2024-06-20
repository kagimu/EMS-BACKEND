package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	authResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	campusRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	cateringRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	schoolRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

func AddMealCoursesToFeedingGroup(courseToFeedingGroupRequest requests.MealCourseToFeedingGroupRequest, currentUser string) error {

	feedingGroup, err := cateringRepositories.GetFeedingGroupByNumber(courseToFeedingGroupRequest.FeedingGroupNumber)

	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	if feedingGroup.Status != constants.StatusActive {
		return errors.New("selected group is not active")
	}

	trn := config.DB.Begin()

	for _, day := range courseToFeedingGroupRequest.Days {

		if checkErr := commons.CheckDay(day); checkErr != "" {
			return errors.New(checkErr)
		}

		groupMealCourses, err := cateringRepositories.GetAllFeedingGroupMealCoursesForDay(feedingGroup.Number, day, trn)
		if err != nil {
			trn.Rollback()
			return errors.New("failed to get all feeding group meal courses: " + err.Error())
		}

		for _, groupMealCourse := range groupMealCourses {
			err = cateringRepositories.DeleteFeedingGroupMealCourse(groupMealCourse, trn)
			if err != nil {
				trn.Rollback()
				return errors.New("failed to clear feeding group meal course: " + err.Error())
			}
		}

		if err = saveMealCoursesForDay(courseToFeedingGroupRequest.MealCourses, day, feedingGroup, currentUser, trn); err != nil {
			trn.Rollback()
			return err
		}

	}

	if err = trn.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func saveMealCoursesForDay(mealCourses []string, day string, feedingGroup models.FeedingGroup, currentUser string, trn *gorm.DB) error {

	for _, mealCourse := range mealCourses {
		checkMeal, err := cateringRepositories.GetMealCourseByNumber(mealCourse, trn)
		if err != nil {
			return err
		}
		if checkMeal.Number == "" {
			return errors.New(fmt.Sprintf("invalid '%s' meal ID", day))
		}

		if checkMeal.Status != constants.StatusActive {
			return errors.New(fmt.Sprintf("meal course '%s' is not active", checkMeal.Name))
		}

		groupMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(checkMeal.Number, day, feedingGroup.Number, trn)
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

		if err := cateringRepositories.SaveFeedingGroupMealCourse(newGroupMealCourse, trn); err != nil {
			return err
		}

	}

	return nil
}

func GetFeedingGroupMealCourseCalendar(schoolNumber string) (responses.FeedingGroupMealCourseCalendar, error) {

	var response responses.FeedingGroupMealCourseCalendar

	school, err := schoolRepositories.GetSchoolByNumber(schoolNumber)

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

	groups, err := cateringRepositories.GetAllFeedingGroups(school.Number, constants.StatusActive)

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

		mealCourses, err := cateringRepositories.GetAllFeedingGroupMealCoursesForDay(groupNumber, day)

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

func GetUserMeals(mealCourseNumber string, campusNumber string, feedingGroupNumber string, weekDay ...string) ([]responses.UserMealResponse, error) {

	mealsResponse := make([]responses.UserMealResponse, 0)

	currentTime := time.Now()
	day := currentTime.Format("Mon")

	if len(weekDay) > 0 {
		day = weekDay[0]
	} else {
		checkDay := commons.CheckDay(day)
		if checkDay != "" {
			return mealsResponse, errors.New(checkDay)
		}
	}

	campus, err := campusRepositories.GetCampusByNumber(campusNumber)

	if err != nil {
		return mealsResponse, err
	}

	if campus.Number == "" {
		return mealsResponse, errors.New("invalid campus number")
	}

	mealCourse, err := cateringRepositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourse.Status != constants.StatusActive {
		return mealsResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}

	checkSchool, err := cateringRepositories.GetMealCourseByNumberAndSchool(mealCourse.Number, campus.SchoolNumber)
	if err != nil {
		return mealsResponse, err
	}

	if checkSchool.Number == "" {
		return mealsResponse, errors.New("meal course can not be used on the selected campus")
	}

	mealCourseFeedingSession, err := cateringRepositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourseFeedingSession.Number == "" {
		return mealsResponse, errors.New("menu not set up for this meal course on '" + day + "'")
	}

	feedingSession := mealCourseFeedingSession.FeedingSession

	feedItems, err := cateringRepositories.GetFeedingSessionItems(feedingSession.Number)
	if err != nil {
		return mealsResponse, err
	}

	if len(feedItems) == 0 {
		return mealsResponse, errors.New("no items set to '" + feedingSession.Number + "' session")
	}

	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {

		checkGroup, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkGroup.Number == "" || checkGroup.SchoolNumber != campus.SchoolNumber {
			return mealsResponse, errors.New("invalid group number")
		}

		checkAllowedToHaveMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkAllowedToHaveMeal.Number == "" {
			return mealsResponse, errors.New("this group should not be served today")
		}

		feedingGroups = append(feedingGroups, checkGroup)

	} else {

		allFeedingGroupsForMealAndDay, err := cateringRepositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return mealsResponse, err
		}

		for _, feedingGroup := range allFeedingGroupsForMealAndDay {
			feedingGroups = append(feedingGroups, feedingGroup.FeedingGroup)
		}

	}

	for _, feedingGroup := range feedingGroups {

		allUsers, err := cateringRepositories.GetFeedingGroupUserByGroupNumber("", "", feedingGroup.Number)
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

		menuItem, err := cateringRepositories.GetMenuItemByNumber(feedingSessionItem.MenuItemNumber)
		if err != nil {
			fmt.Println("Error getting menu item from feeding session")
			return items
		}

		unitOfMeasure, err := cateringRepositories.GetUnitOfMeasureByNumber(menuItem.UnitOfMeasureNumber)
		if err != nil {
			fmt.Println("Error getting unit of measure" + err.Error())
			return items
		}

		unitOfMeasureName := "None"
		if unitOfMeasure.Number != "" {
			unitOfMeasureName = unitOfMeasure.Name
		}

		var served float64

		servedStock, err := cateringRepositories.GetMealStockingByMealNumberAndDateUserNumber(mealCourse, date_, userNumber, menuItem.Number)

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
			UnitOfMeasure:  unitOfMeasureName,
		}

		items = append(items, singleUserMenu)
	}

	return items
}

func GetFeedingGroupsBySchoolNumber(schoolNumber string, status string) ([]responses.FeedingGroupResponse, error) {

	var response []responses.FeedingGroupResponse

	feedingGroups, err := cateringRepositories.GetAllFeedingGroups(schoolNumber, status)
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

		school, err := schoolRepositories.GetSchoolByNumber(feedingGroup.SchoolNumber)
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

	currentStock, err := cateringRepositories.GetMenuItemCurrentStock(mealRequest.MenuItemNumber)
	if err != nil {
		return err
	}

	if mealRequest.Quantity <= 0 {
		return errors.New("quantities served can not be zero")
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

	newMeal, err := cateringRepositories.GetMealByMealNumberAndDateAndUserNumber(mealRequest.MealCourseNumber, formatDate, mealRequest.UserNumber)
	if err != nil {
		return errors.New("Checking if item served error: " + err.Error())
	}

	groupUser, err := cateringRepositories.GetFeedingGroupUserByUserNumber(mealRequest.UserNumber)
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
		newMeal.FeedingGroupNumber = groupUser.FeedingGroupNumber
	} else {
		newMeal.StaffNumber = currentUser
		newMeal.CampusNumber = groupUser.CampusNumber
		newMeal.FeedingSessionNumber = mealRequest.FeedingSessionNumber
		newMeal.FeedingGroupNumber = groupUser.FeedingGroupNumber
	}

	err = cateringRepositories.SaveMeal(tx, newMeal)
	if err != nil {
		return errors.New("error creating meal: " + err.Error())
	}

	stocking, err := cateringRepositories.GetMealStockByMealNumberAndMenuItemNumber(tx, newMeal.Number, mealRequest.MenuItemNumber)
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

	err = cateringRepositories.SaveStocking(tx, stocking)
	if err != nil {
		return errors.New("error reducing stock: " + err.Error())
	}

	mealStocking, err := cateringRepositories.GetMealStockingByMealNumberAndStockingNumber(tx, newMeal.Number, stocking.Number)
	if err != nil {
		return err
	}

	if mealStocking.MealNumber == "" {
		mealStocking.MealNumber = newMeal.Number
		mealStocking.StockingNumber = stocking.Number
	}

	err = cateringRepositories.SaveMealStocking(tx, mealStocking)
	if err != nil {
		return errors.New("error mapping stock with meal taken: " + err.Error())
	}

	if err = tx.Commit().Error; err != nil {
		return errors.New("Commit error: " + err.Error())
	}

	return nil
}

func UpdateFeedingGroup(payload requests.UpdateFeedingGroupRequest, currentUserNumber string) error {

	feedingGroup, err := cateringRepositories.GetFeedingGroupByNumber(payload.FeedingGroupNumber)
	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	checkFeedingGroupDup, err := cateringRepositories.GetFeedingGroupByName(payload.Name, feedingGroup.SchoolNumber)
	if err != nil {
		return err
	}

	if checkFeedingGroupDup.Number != "" && checkFeedingGroupDup.Number != feedingGroup.Number {
		return errors.New("this name is already defined")
	}

	feedingGroup.Status = payload.Status
	feedingGroup.Name = payload.Name
	feedingGroup.UserNumber = currentUserNumber

	err = cateringRepositories.SaveFeedingGroup(feedingGroup)
	if err != nil {
		return err
	}

	return nil
}

func CreateFeedingGroup(payload requests.FeedingGroupRequest, currentUserNumber string) error {

	checkName, err := cateringRepositories.GetFeedingGroupByName(payload.Name, payload.SchoolNumber)
	if err != nil {
		return err
	}
	if checkName.Number != "" {
		return errors.New("group with same name already defined")
	}

	newFeedingGroup := models.FeedingGroup{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		SchoolNumber: payload.SchoolNumber,
		Status:       payload.Status,
		UserNumber:   currentUserNumber,
	}

	err = cateringRepositories.SaveFeedingGroup(newFeedingGroup)
	if err != nil {
		return err
	}

	return nil
}

func SaveFeedingSession(payload requests.FeedingSessionRequest, currentUserNumber string) error {

	checkFeedingSession, err := cateringRepositories.GetFeedingSessionByNameAndSchool(payload.Name, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if checkFeedingSession.Number != "" {
		return errors.New("this name is already defined")
	}

	checkFeedingSession, err = cateringRepositories.GetFeedingSessionByCodeAndSchool(payload.Code, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if checkFeedingSession.Number != "" {
		return errors.New("the provided code is already used")
	}

	newFeedingSession := models.FeedingSession{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		Status:       payload.Status,
		SchoolNumber: payload.SchoolNumber,
		Code:         payload.Code,
		UserNumber:   currentUserNumber,
	}

	err = cateringRepositories.SaveFeedingSession(newFeedingSession)
	if err != nil {
		return err
	}

	return nil
}

func GetFeedingSession(schoolNumber string) ([]responses.FeedingSessionResponse, error) {

	var feedingSessions []responses.FeedingSessionResponse

	sessions, err := cateringRepositories.GetFeedingSessionsBySchoolNumber(schoolNumber)

	if err != nil {
		return feedingSessions, err
	}

	for _, session := range sessions {
		feedingSessions = append(feedingSessions, BuildFeedingSessionResponse(session))
	}

	return feedingSessions, nil
}

func BuildFeedingSessionResponse(session models.FeedingSession) responses.FeedingSessionResponse {

	response := responses.FeedingSessionResponse{
		Name:   session.Name,
		Code:   session.Code,
		Number: session.Number,
	}

	user, err := repository.GetUserByNumber(session.UserNumber)
	if err != nil {
		fmt.Println("Getting session user error", err.Error())
	} else if user.Number == "" {
		response.AddedBy = user.FirstName
	}

	return response
}

func MealUsers(schoolNumber string, campusNumber string, groupNumber string) ([]responses.FeedingUserGroupResponse, error) {

	var response = make([]responses.FeedingUserGroupResponse, 0)

	users, err := cateringRepositories.GetFeedingGroupUserByGroupNumber(schoolNumber, campusNumber, groupNumber)
	if err != nil {
		return nil, errors.New("Getting group users error: " + err.Error())
	}

	for _, user := range users {

		response = append(response, createFeedingUserResponse(user))

	}

	return response, nil

}

func createFeedingUserResponse(user models.FeedingGroupUser) responses.FeedingUserGroupResponse {

	_user := user.User

	userNames := _user.FirstName + " " + _user.LastName

	return responses.FeedingUserGroupResponse{
		Number:             user.Number,
		UserNumber:         user.UserNumber,
		FeedingGroupName:   user.FeedingGroup.Name,
		FeedingGroupNumber: user.FeedingGroupNumber,
		CampusName:         user.Campus.Name,
		Gender:             user.User.Gender,
		Name:               userNames,
	}
}

func AvailableUser(schoolNumber string, userType string) ([]authResponse.UserResponse, error) {

	var response = make([]authResponse.UserResponse, 0)

	users, err := cateringRepositories.AvailableUsers(schoolNumber, constants.StatusActive, userType)
	if err != nil {
		return nil, errors.New("Getting available users error: " + err.Error())
	}

	for _, user := range users {
		response = append(response, services.CreateUserResponse(user))
	}

	return response, nil
}

func DeleteFeedingGroup(feedingGroupNumber string) error {

	group, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
	if err != nil {
		return errors.New("Checking group number error: " + err.Error())
	}

	if group.Number == "" {
		return errors.New("invalid group number")
	}

	err = cateringRepositories.RemoveFeedingGroup(group)
	if err != nil {
		fmt.Println("Removing feed group error: " + err.Error())
		return errors.New("Can not remove '" + group.Name + "' group, try removing all its usages first.")
	}

	return nil

}

func UpdateFeedingSession(updateFeedingSessionRequest requests.UpdateFeedingSessionRequest, currentUserNumber string) error {
	feedingSession, err := cateringRepositories.GetFeedingSessionByNumberAndSchoolNumber(updateFeedingSessionRequest.FeedingSessionNumber, updateFeedingSessionRequest.SchoolNumber)

	if err != nil {
		return errors.New("Checking feed session number error: " + err.Error())
	}

	if feedingSession.Number == "" {
		return errors.New("invalid feed session number")
	}

	checkName, err := cateringRepositories.GetFeedingSessionByNameAndSchool(updateFeedingSessionRequest.Name, feedingSession.SchoolNumber)
	if err != nil {
		return errors.New("Checking feed session name error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != updateFeedingSessionRequest.FeedingSessionNumber {
		return errors.New("name already used, choose a unique one")
	}

	checkCode, err := cateringRepositories.GetFeedingSessionByCodeAndSchool(updateFeedingSessionRequest.Code, updateFeedingSessionRequest.SchoolNumber)
	if err != nil {
		return errors.New("Checking feeding session code error: " + err.Error())
	}

	if checkCode.Number != "" && checkCode.Number != updateFeedingSessionRequest.FeedingSessionNumber {
		return errors.New("feeding session code already used, choose a unique one")
	}

	feedingSession.Code = updateFeedingSessionRequest.Code
	feedingSession.UserNumber = currentUserNumber
	feedingSession.Name = updateFeedingSessionRequest.Name
	feedingSession.Status = updateFeedingSessionRequest.Status

	err = cateringRepositories.SaveFeedingSession(feedingSession)
	if err != nil {
		return errors.New("updating feed session error: " + err.Error())
	}

	return nil
}

func DeleteFeedingSession(feedingSessionNumber string) error {

	feedingSession, err := cateringRepositories.GetFeedingSessionByNumberAndSchoolNumber(feedingSessionNumber, "")
	if err != nil {
		return errors.New("Checking feed session number error: " + err.Error())
	}

	if feedingSession.Number == "" {
		return errors.New("invalid feed session number")
	}

	err = cateringRepositories.DeleteFeedingSession(feedingSession)
	if err != nil {
		fmt.Println("Deleting feed session error: " + err.Error())
		return errors.New("this feeding session could not be removed, try removing all its usages first")
	}

	return nil
}
