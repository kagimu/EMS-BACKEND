package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
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
		return err
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	var oldMealCourseCount models.MealCourseCount

	oldMealCourseCount, err = repositories.GetMealCourseCountBySchoolNumber(school.Number)

	if err != nil {
		return err
	}

	newMealCourseCount := models.MealCourseCount{
		Count:        mealCourseCountRequest.Count,
		UserNumber:   currentUserNumber,
		Number:       uuid.NewV4().String(),
		SchoolNumber: school.Number,
	}

	if oldMealCourseCount.Number != "" {
		newMealCourseCount.Number = oldMealCourseCount.Number
		newMealCourseCount.ID = oldMealCourseCount.ID
	}

	if _, err = repositories.SaveMealCourseCount(newMealCourseCount); err != nil {
		return err
	}

	allOtherMealCourses, err := repositories.GetAllMealCoursesByMealCourseCount(newMealCourseCount.Number, constants.STATUS_ACTIVE)
	if err != nil {
		return err
	}

	for _, otherMealCourse := range allOtherMealCourses {

		otherMealCourse.Status = constants.STATUS_INACTIVE

		if _, err = repositories.SaveMealCourse(otherMealCourse); err != nil {
			return err
		}
	}

	for _, mealCourse := range mealCourseCountRequest.Courses {

		newMealCourse := models.MealCourse{
			Number:                uuid.NewV4().String(),
			Name:                  mealCourse.Name,
			MealCourseCountNumber: newMealCourseCount.Number,
			Status:                constants.STATUS_ACTIVE,
			UserNumber:            currentUserNumber,
			Code:                  mealCourse.Code,
		}

		var oldMealCourse models.MealCourse

		oldMealCourse, err = repositories.GetMealCourseByName(mealCourse.Name, newMealCourseCount.Number)
		if err != nil {
			return err
		}

		if oldMealCourse.Number != "" {
			newMealCourse.Number = oldMealCourse.Number
		}

		if _, err = repositories.SaveMealCourse(newMealCourse); err != nil {
			return err
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

func MealCourseCalendar(schoolNumber string) (responses.MealCourseCalendarResponse, error) {

	var mealCourseCalendarResponse responses.MealCourseCalendarResponse

	meal, err := repositories.GetMealCourseCountBySchoolNumber(schoolNumber)
	if err != nil {
		return mealCourseCalendarResponse, err
	}

	if meal.Number == "" {
		return mealCourseCalendarResponse, errors.New("no meal courses configured for the selected school")
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

			item := responses.DayItemsResponse{
				Day:   day,
				Items: []string{"Beans", "Posho", "Rice"},
			}

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
			Name: _session.Name,
			Code: _session.Code,
		}

		feedingSessionsResponse = append(feedingSessionsResponse, session)
	}

	return feedingSessionsResponse, nil
}
