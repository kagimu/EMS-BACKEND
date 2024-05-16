package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	campusRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
)

func SaveMealCourse(mealCourseCountRequest requests.MealCourseCountRequest, currentUserNumber string) error {

	var school schoolModels.School

	school, err := repositories2.GetSchoolByNumber(mealCourseCountRequest.SchoolNumber)

	if err != nil {
		return errors.New("checking for valid school number error : " + err.Error())
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	var oldMealCourseCount models.MealCourseCount

	oldMealCourseCount, err = repositories.GetMealCourseCountBySchoolNumber(school.Number)

	if err != nil {
		return errors.New("getting course count by school number: " + err.Error())
	}

	newMealCourseCount := models.MealCourseCount{
		Count:        mealCourseCountRequest.Count,
		UserNumber:   currentUserNumber,
		Number:       uuid.NewV4().String(),
		SchoolNumber: school.Number,
	}

	if oldMealCourseCount.Number == "" {
		oldMealCourseCount = newMealCourseCount
	} else {
		oldMealCourseCount.UserNumber = currentUserNumber
		oldMealCourseCount.Count = mealCourseCountRequest.Count
	}

	if _, err = repositories.SaveMealCourseCount(oldMealCourseCount); err != nil {
		return errors.New("saving or updating meal course count error: " + err.Error())
	}

	allOtherMealCourses, err := repositories.GetAllMealCoursesByMealCourseCount(oldMealCourseCount.Number, constants.STATUS_ACTIVE)
	if err != nil {
		return errors.New("getting all active meal course by course count number error : " + err.Error())
	}

	for _, otherMealCourse := range allOtherMealCourses {

		otherMealCourse.Status = constants.STATUS_INACTIVE

		if _, err = repositories.SaveMealCourse(otherMealCourse); err != nil {
			return errors.New("updating current meal course names to inactive : " + err.Error())
		}
	}

	for _, mealCourse := range mealCourseCountRequest.Courses {

		newMealCourse := models.MealCourse{
			Number:                uuid.NewV4().String(),
			Name:                  mealCourse.Name,
			MealCourseCountNumber: oldMealCourseCount.Number,
			Status:                constants.STATUS_ACTIVE,
			UserNumber:            currentUserNumber,
			Code:                  mealCourse.Code,
		}

		var checkMeal models.MealCourse

		checkMeal, err = repositories.GetMealCourseByName(mealCourse.Name, oldMealCourseCount.Number)
		if err != nil {
			return errors.New("checking meal course by name error : " + err.Error())
		}

		if checkMeal.Number != "" {
			checkMeal.Name = newMealCourse.Name
			checkMeal.Status = newMealCourse.Status
			checkMeal.UserNumber = newMealCourse.UserNumber
			checkMeal.Code = newMealCourse.Code
		} else {
			checkMeal = newMealCourse
		}

		if _, err = repositories.SaveMealCourse(checkMeal); err != nil {
			return errors.New("save or update meal course error : " + err.Error())
		}

	}

	return nil

}

func GetMealCourses(schoolNumber string) (responses.MealCourseResponse, error) {

	school, err := repositories2.GetSchoolByNumber(schoolNumber)

	if err != nil {
		return responses.MealCourseResponse{}, err
	}

	mealCourseCount, err := repositories.GetMealCourseCountBySchoolNumber(school.Number)
	if err != nil {
		return responses.MealCourseResponse{}, err
	}

	courses, err := repositories.GetAllMealCoursesByMealCourseCount(mealCourseCount.Number, constants.STATUS_ACTIVE)
	if err != nil {
		return responses.MealCourseResponse{}, err
	}

	var mealCourseItemResponses []responses.MealCourseItemResponse

	for _, course := range courses {

		user, _ := repository.GetUserByNumber(course.UserNumber)

		courseResponse := responses.MealCourseItemResponse{
			Name:    course.Name,
			Code:    course.Code,
			AddedBy: user.FirstName,
			Number:  course.Number,
		}

		mealCourseItemResponses = append(mealCourseItemResponses, courseResponse)

	}

	courseResponse := responses.MealCourseResponse{
		Count:        mealCourseCount.Count,
		SchoolNumber: schoolNumber,
		Courses:      mealCourseItemResponses,
	}

	return courseResponse, nil

}

func MealCourseCalendar(campusNumber string) (responses.MealCourseCalendarResponse, error) {

	var mealCourseCalendarResponse responses.MealCourseCalendarResponse

	campus, err := campusRepositories.GetCampusByNumber(campusNumber)

	if err != nil {
		return mealCourseCalendarResponse, err
	}

	if campus.Number == "" {
		return mealCourseCalendarResponse, errors.New("invalid campus number")
	}

	schoolNumber := campus.SchoolNumber
	meal, err := repositories.GetMealCourseCountBySchoolNumber(schoolNumber)
	if err != nil {
		return mealCourseCalendarResponse, err
	}

	if meal.Number == "" {
		return mealCourseCalendarResponse, errors.New("no meal courses configured for the school")
	}

	courses, err := repositories.GetAllMealCoursesByMealCourseCount(meal.Number, constants.STATUS_ACTIVE)
	if err != nil {
		return mealCourseCalendarResponse, err
	}

	_days := commons.DaysOfTheWeek()

	headers := []string{" "}

	for _, rowDay := range _days {
		headers = append(headers, rowDay)
	}

	mealCourseCalendarResponse.Headers = headers

	var calendarRawResponse []responses.CalendarRowResponse

	for _, course := range courses {

		var daysItems []responses.DayItemsResponse

		for _, day := range _days {

			item := processItemsOnDayForMealCourseAtCampus(day, course.Number, campusNumber)

			daysItems = append(daysItems, item)
		}

		calendarRaw := responses.CalendarRowResponse{
			CourseName: course.Name,
			DayItems:   daysItems,
		}

		calendarRawResponse = append(calendarRawResponse, calendarRaw)
	}

	mealCourseCalendarResponse.CalendarRawResponse = calendarRawResponse

	return mealCourseCalendarResponse, nil

}

func processItemsOnDayForMealCourseAtCampus(day string, mealCourseNumber string, campusNumber string) responses.DayItemsResponse {

	item := responses.DayItemsResponse{
		Day:    day,
		Active: false,
	}

	mealFeedingSession, err := repositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)

	if err != nil {
		item.Error = err.Error()
	} else {

		if mealFeedingSession.Number == "" {
			item.Error = "no feeding session set for day on this course"
		} else {

			feedingSession := mealFeedingSession.FeedingSession
			item.Name = feedingSession.Name

			itemsList, err := processFeedingSessionItemNames(feedingSession.Number)
			if err != nil {
				item.Error = err.Error()
			} else {
				item.Active = true
				item.Items = itemsList
			}
		}

	}

	return item
}

func processFeedingSessionItemNames(feedingSessionNumber string) ([]string, error) {

	var items []string

	feedingSessionItems, err := repositories.GetFeedingSessionItems(feedingSessionNumber)
	if err != nil {
		return items, err
	}

	if len(feedingSessionItems) == 0 {
		return items, errors.New("no menu items set")
	}

	for _, item := range feedingSessionItems {

		items = append(items, item.MenuItem.Name)

	}

	return items, nil
}

func AddFeedingSessionsToMealCourse(mealCourseFeedingSessionRequest requests.MealCourseFeedingSessionRequest) error {

	mealCourse, err := repositories.GetMealCourseByNumber(mealCourseFeedingSessionRequest.MealCourseNumber)

	if err != nil {
		return err
	}

	if mealCourse.Number == "" {
		return errors.New("invalid meal course selected")
	}

	for _, feedingSession := range mealCourseFeedingSessionRequest.FeedingSessions {

		mealCourseFeedingSession, err := repositories.GetMealCourseFeedingSession(mealCourse.Name, feedingSession)
		if err != nil {
			return err
		}

		if mealCourseFeedingSession.Number != "" {
			return errors.New("feeding session already added to the meal course")
		}

		newMealCourseFeedingSession := models.MealCourseFeedingSession{
			MealCourseNumber:     mealCourse.Number,
			FeedingSessionNumber: feedingSession,
			Number:               uuid.NewV4().String(),
		}

		err = repositories.SaveMealCourseFeedingSession(newMealCourseFeedingSession)
		if err != nil {
			return err
		}

	}

	return nil
}

func MealCourseFeedingSessions(mealCourseNumber string) ([]responses.FeedingSessionResponse, error) {

	var feedingSessionsResponse []responses.FeedingSessionResponse

	mealCourse, err := repositories.GetMealCourseByNumber(mealCourseNumber)

	if err != nil {
		return feedingSessionsResponse, err
	}

	if mealCourse.Number == "" {
		return feedingSessionsResponse, errors.New("invalid meal course selected")
	}

	feedingSessions, err := repositories.GetMealCourseFeedingSessionsByMealCourseNumber(mealCourse.Number)
	if err != nil {
		return feedingSessionsResponse, err
	}

	for _, feedingSession := range feedingSessions {

		_session := feedingSession.FeedingSession

		session := responses.FeedingSessionResponse{
			Name:   _session.Name,
			Code:   _session.Code,
			Number: _session.Number,
		}

		feedingSessionsResponse = append(feedingSessionsResponse, session)
	}

	return feedingSessionsResponse, nil
}

func ActivateMealCourse(activateMealCourseRequest requests.ActivateMealCourseRequest, currentUserNumber string) error {

	mealCourse, err := repositories.GetMealCourseByNumber(activateMealCourseRequest.MealCourseNumber)
	if err != nil {
		return err
	}

	if mealCourse.Status != constants.STATUS_ACTIVE {
		return errors.New("selected meal course is not active")
	}

	feedingSession, err := repositories.GetFeedingSessionByNumber(activateMealCourseRequest.FeedingSessionNumber)

	if err != nil {
		return err
	}

	if feedingSession.Status != constants.STATUS_ACTIVE {
		return errors.New("selected feed session is not active")
	}

	mealCourseFeedingSession, err := repositories.GetMealCourseFeedingSession(mealCourse.Number, feedingSession.Number)

	if err != nil {
		return err
	}

	if mealCourseFeedingSession.Number == "" {
		return errors.New("feeding session not added to this meal course")
	}

	for _, day := range activateMealCourseRequest.Days {
		if err = activateForDay(day, activateMealCourseRequest.Campuses, currentUserNumber, mealCourseFeedingSession.Number); err != nil {
			return err
		}
	}

	return nil
}

func activateForDay(day string, campuses []string, currentUserNumber string, mealCourseFeedingSessionNumber string) error {

	for _, campusNumber := range campuses {

		campus, err := campusRepositories.GetCampusByNumber(campusNumber)

		campusMealCourse, err := repositories.GetCampusMealCourseByCampusAndDay(day, mealCourseFeedingSessionNumber, campusNumber)
		if err != nil {
			return err
		}

		if campusMealCourse.Number != "" {
			return errors.New("this campus already has this meal course activated")
		}

		campusMealCourse = models.CampusMealCourse{
			Number:                         uuid.NewV4().String(),
			UserNumber:                     currentUserNumber,
			MealCourseFeedingSessionNumber: mealCourseFeedingSessionNumber,
			Day:                            day,
			CampusNumber:                   campus.Number,
		}

		if err = repositories.SaveCampusMealCourse(campusMealCourse); err != nil {
			return err
		}

	}

	return nil
}
