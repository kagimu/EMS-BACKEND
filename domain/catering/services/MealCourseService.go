package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
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
		SchoolNumber: mealCourseCountRequest.SchoolNumber,
	}

	if oldMealCourseCount.Number != "" {
		newMealCourseCount.Number = oldMealCourseCount.Number
	}

	if _, err = repositories.SaveMealCourseCount(newMealCourseCount); err != nil {
		return err
	}

	for _, mealCourse := range mealCourseCountRequest.Courses {

		newMealCourse := models.MealCourse{
			Number:                uuid.NewV4().String(),
			Name:                  mealCourse.Name,
			MealCourseCountNumber: newMealCourseCount.Number,
			Status:                constants.STATUS_ACTIVE,
			UserNumber:            currentUserNumber,
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

	courses, err := repositories.GetMealCourses(mealCourseCount.Number, constants.STATUS_ACTIVE)
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
