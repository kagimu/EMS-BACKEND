package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	uuid "github.com/satori/go.uuid"
)

func CreateGrade(userNumber string, request requests.GradeDataRequest) error {

	checkGrade, err := repository.FindGradeByCode(request.Code, request.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade code: " + err.Error())
	}

	if checkGrade.Number != "" {
		return errors.New("grade code already used")
	}

	checkName, err := repository.FindGradeByNameAndShortName(request.Name, request.ShortName, request.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade name: " + err.Error())
	}

	if checkName.Number != "" {
		return errors.New("grade name with short name already used")
	}

	var grade = models.Grade{
		Number:       uuid.NewV4().String(),
		Name:         request.Name,
		Code:         request.Code,
		UserNumber:   userNumber,
		ShortName:    request.ShortName,
		SchoolNumber: request.SchoolNumber,
	}

	if err := repository.SaveGrade(grade); err != nil {
		return errors.New("creating grade error: " + err.Error())
	}

	return nil
}

func FetchGrades(schoolNumber string) ([]models.Grade, error) {
	var grades []models.Grade
	grades, err := repository.FindAllGrades(schoolNumber)
	if err != nil {
		return nil, err
	}
	return grades, nil
}

func UpdateGrade(userNumber string, request requests.EditGradeDataRequest) error {
	grade, err := repository.FindGradeByNumber(request.GradeNumber)
	if err != nil {
		return errors.New("error checking grade code: " + err.Error())
	}

	if grade.Number == "" {
		return errors.New("invalid grade number")
	}

	checkGrade, err := repository.FindGradeByCode(request.Code, request.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade code: " + err.Error())
	}

	if checkGrade.Number != "" && grade.Number != checkGrade.Number {
		return errors.New("grade code already used")
	}

	checkName, err := repository.FindGradeByNameAndShortName(request.Name, request.ShortName, request.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade name: " + err.Error())
	}

	if checkName.Number != "" && grade.Number != checkName.Number {
		return errors.New("grade name with short name already used")
	}

	grade.Name = request.Name
	grade.SchoolNumber = request.SchoolNumber
	grade.Code = request.Code
	grade.ShortName = request.ShortName
	grade.UserNumber = userNumber

	if err := repository.UpdateGrade(grade); err != nil {
		return errors.New("error updating grade: " + err.Error())
	}

	return nil
}
