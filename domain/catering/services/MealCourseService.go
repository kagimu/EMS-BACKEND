package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
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

	var oldMealCourse models.MealCourseCount

	oldMealCourse, err = repositories.GetMealCourseCountBySchoolNumber(school.Number)

	if err != nil {
		return err
	}

	newMealCourse := models.MealCourseCount{
		Count:        mealCourseCountRequest.Count,
		UserNumber:   currentUserNumber,
		Number:       uuid.NewV4().String(),
		SchoolNumber: mealCourseCountRequest.SchoolNumber,
	}

	if oldMealCourse.Number != "" {
		newMealCourse.Number = oldMealCourse.Number
	}

	if _, err = repositories.SaveMealCourseCount(newMealCourse); err != nil {
		return err
	}

	return nil

}
