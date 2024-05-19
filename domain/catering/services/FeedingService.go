package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
)

func AddMealCoursesToFeedingGroup(courseToFeedingGroupRequest requests.MealCourseToFeedingGroupRequest, currentUser string) error {

	feedingGroup, err := repositories.GetFeedingGroupByNumber(courseToFeedingGroupRequest.FeedingGroupNumber)

	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	if feedingGroup.Status != constants.STATUS_ACTIVE {
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

		if checkMeal.Status != constants.STATUS_ACTIVE {
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

	groups, err := repositories.GetAllFeedingGroups(school.Number, constants.STATUS_ACTIVE)

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

		mealCourses, err := repositories.GetAllFeedingGroupMealCoursesForDay(groupNumber, day)

		var meals []string

		if err != nil {
			item.Error = err.Error()
		} else {

			for _, _mealCourse := range mealCourses {
				mealCourse := _mealCourse.MealCourse
				if mealCourse.Status == constants.STATUS_ACTIVE {
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
		weeklyItems = append(weeklyItems, item)

	}

	return weeklyItems
}
