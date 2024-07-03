package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func SaveGradingStandard(gradingStandardRequest requests.GradingStandardRequest, currentUser string) error {

	byCode, err := repositories.FindGradingStandardByCodeAndClassNumber(gradingStandardRequest.Code, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking code error: " + err.Error())
	}

	if byCode.Number != "" {
		return errors.New("code already used: " + byCode.Number)
	}

	byName, err := repositories.FindGradingStandardByNameAndClassNumber(gradingStandardRequest.Name, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking by name error: " + err.Error())
	}

	if byName.Number != "" {
		return errors.New("name already used: " + byName.Name)
	}

	gs := models.GradingStandard{
		Code:              gradingStandardRequest.Code,
		Number:            uuid.NewV4().String(),
		Name:              gradingStandardRequest.Name,
		Description:       gradingStandardRequest.Description,
		Type:              gradingStandardRequest.Type,
		CampusClassNumber: gradingStandardRequest.CampusClassNumber,
		UserNumber:        currentUser,
		CampusNumber:      gradingStandardRequest.CampusNumber,
	}

	err = repositories.SaveGradingStandard(gs)
	if err != nil {
		return errors.New("Saving grading standard error: " + err.Error())
	}

	return nil
}

func UpdateGradingStandard(gradingStandardRequest requests.UpdateGradingStandardRequest, currentUser string) error {

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardRequest.GradingStandardNumber)
	if err != nil {
		return errors.New("Checking grading standard error: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return errors.New("grading standard number is empty: " + gradingStandardRequest.GradingStandardNumber)
	}

	byCode, err := repositories.FindGradingStandardByCodeAndClassNumber(gradingStandardRequest.Code, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking code error: " + err.Error())
	}

	if byCode.Number != "" && byCode.Number != gradingStandard.Number {
		return errors.New("code already used: " + byCode.Number)
	}

	byName, err := repositories.FindGradingStandardByNameAndClassNumber(gradingStandardRequest.Name, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking by name error: " + err.Error())
	}

	if byName.Number != "" && byCode.Number != gradingStandard.Number {
		return errors.New("name already used: " + byName.Name)
	}

	gradingStandard.Code = gradingStandardRequest.Code
	gradingStandard.Name = gradingStandardRequest.Name
	gradingStandard.Description = gradingStandardRequest.Description
	gradingStandard.Type = gradingStandardRequest.Type
	gradingStandard.CampusClassNumber = gradingStandardRequest.CampusClassNumber
	gradingStandard.CampusNumber = gradingStandardRequest.CampusNumber
	gradingStandard.UserNumber = currentUser

	err = repositories.SaveGradingStandard(gradingStandard)
	if err != nil {
		return errors.New("Saving grading standard error: " + err.Error())
	}

	return nil
}

func SaveAggregateScale(aggregateScaleRequest requests.AggregateScaleRequest, currentUser string) error {
	aggregateStandard, err := repositories.GetGradingStandardByNumber(aggregateScaleRequest.GradingStandardNumber)
	if err != nil {
		return errors.New("Checking grading standard error: " + err.Error())
	}

	if aggregateStandard.Number == "" {
		return errors.New("invalid grading standard number: " + aggregateScaleRequest.GradingStandardNumber)
	}

	oldScales, err := repositories.FindAggregateScaleByGradingStandardNumber(aggregateStandard.Number)
	if err != nil {
		return errors.New("checking previous scale error: " + err.Error())
	}

	grades := aggregateScaleRequest.Grades

	err = checkDuplicates(grades)
	if err != nil {
		return err
	}

	trn := config.DB.Begin()

	for _, oldScale := range oldScales {
		err = repositories.RemoveAggregateScale(oldScale, trn)
		if err != nil {
			return errors.New("failed clearing old scale, try again.")
		}
	}

	for _, grade := range grades {

		checkDuplicate, err := repositories.GetAggregateScaleByGradingStandardNumberAndGradeNumber(aggregateStandard.Number, grade.GradeNumber, trn)
		if err != nil {
			return errors.New("checking grade duplicate error: " + err.Error())
		}

		if checkDuplicate.Number != "" {
			return errors.New("you have a used grade in your scale")
		}

		err = createGrade(trn, grade, currentUser, aggregateStandard)
		if err != nil {
			trn.Rollback()
			return err
		}
	}

	if err := trn.Commit().Error; err != nil {
		return errors.New("commit grading scales error: " + err.Error())
	}

	return nil
}

func checkDuplicates(grades []requests.AggregateScaleItem) error {

	cleanGrades := make([]requests.AggregateScaleItem, len(grades))

	for i, grade := range grades {

		shouldReject := ""

		for _, cleanGrade := range cleanGrades {

			if cleanGrade.GradeNumber == grade.GradeNumber {
				shouldReject = fmt.Sprintf("grade %s already exists", grade.GradeNumber)
			}

		}

		if shouldReject == "" {
			cleanGrades[i] = grade
		} else {
			return errors.New(shouldReject)
		}

	}

	return nil
}

func createGrade(trn *gorm.DB, aggregateScaleItem requests.AggregateScaleItem, user string, gradingStandard models.GradingStandard) error {

	newAggregateScale := models.AggregateScale{
		Code:                  aggregateScaleItem.Code,
		Number:                uuid.NewV4().String(),
		MinValue:              aggregateScaleItem.MinValue,
		MaxValue:              aggregateScaleItem.MaxValue,
		CampusClassNumber:     gradingStandard.CampusClassNumber,
		CampusNumber:          gradingStandard.CampusNumber,
		UserNumber:            user,
		GradeNumber:           aggregateScaleItem.GradeNumber,
		GradingStandardNumber: gradingStandard.Number,
	}

	err := repositories.SaveAggregateScale(newAggregateScale, trn)
	if err != nil {
		return errors.New("Saving grading scale error: " + err.Error())
	}

	return nil
}
