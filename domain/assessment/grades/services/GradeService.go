package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/gin-gonic/gin"
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

	checkName, err := repository.FindGradeByName(request.Name)
	if err != nil {
		return errors.New("error checking grade name: " + err.Error())
	}

	if checkName.Number != "" {
		return errors.New("name already used")
	}

	checkNameAndShortName, err := repository.FindGradeByNameAndShortName(request.Name, request.ShortName, request.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade name: " + err.Error())
	}

	if checkNameAndShortName.Number != "" {
		return errors.New("grade name with short name already used")
	}

	var grade = models.Grade{
		Number:           uuid.NewV4().String(),
		Name:             request.Name,
		Code:             request.Code,
		UserNumber:       userNumber,
		ShortName:        request.ShortName,
		SchoolNumber:     request.SchoolNumber,
		GradeGroupNumber: request.GradeGroupNumber,
	}

	if err := repository.SaveGrade(grade); err != nil {
		return errors.New("creating grade error: " + err.Error())
	}

	return nil
}

func FetchGrades(schoolNumber string, ctx *gin.Context) ([]responses.NameOfAggregates, error) {
	var grades []models.Grade
	grades, err := repository.FindAllGrades(schoolNumber, ctx.Query("grade_group_number"))
	if err != nil {
		return nil, err
	}

	var gradeNames []responses.NameOfAggregates
	for _, grade := range grades {

		gradeNames = append(gradeNames, responses.NameOfAggregates{
			Name:        grade.Name,
			Number:      grade.Number,
			ShortName:   grade.ShortName,
			Code:        grade.Code,
			GradeGroup:  grade.GradeGroup.Name,
			Status:      grade.Status,
			GroupNumber: grade.GradeGroup.Number,
		})
	}

	return gradeNames, nil
}

func UpdateGrade(userNumber string, editGradeDataRequest requests.EditGradeDataRequest) error {
	grade, err := repository.FindGradeByNumber(editGradeDataRequest.GradeNumber)
	if err != nil {
		return errors.New("error checking grade code: " + err.Error())
	}

	if grade.Number == "" {
		return errors.New("invalid grade number")
	}

	checkGrade, err := repository.FindGradeByCode(editGradeDataRequest.Code, editGradeDataRequest.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade code: " + err.Error())
	}

	if checkGrade.Number != "" && grade.Number != checkGrade.Number {
		return errors.New("grade code already used")
	}

	checkName, err := repository.FindGradeByNameAndShortName(editGradeDataRequest.Name, editGradeDataRequest.ShortName, editGradeDataRequest.SchoolNumber)
	if err != nil {
		return errors.New("error checking grade name: " + err.Error())
	}

	if checkName.Number != "" && grade.Number != checkName.Number {
		return errors.New("grade name with short name already used")
	}

	grade.Name = editGradeDataRequest.Name
	grade.SchoolNumber = editGradeDataRequest.SchoolNumber
	grade.Code = editGradeDataRequest.Code
	grade.ShortName = editGradeDataRequest.ShortName
	grade.UserNumber = userNumber
	grade.GradeGroupNumber = editGradeDataRequest.GradeGroupNumber

	if err := repository.UpdateGrade(grade); err != nil {
		return errors.New("error updating grade: " + err.Error())
	}

	return nil
}

func RemoveGrade(gradeNumber string) error {
	grade, err := repository.FindGradeByNumber(gradeNumber)
	if err != nil {
		return errors.New("checking grade number error: " + err.Error())
	}

	if grade.Number == "" {
		return errors.New("invalid grade number")
	}

	err = repository.RemoveGrade(grade)
	if err != nil {
		fmt.Println("removing grade error: " + err.Error())

		if grade.Status == constants.StatusActive {

			grade.Status = constants.StatusInactive
			err = repository.UpdateGrade(grade)
			if err != nil {
				return errors.New("grade could not be deleted or updated to: " + grade.Status)
			}

		}

	}

	return err
}
