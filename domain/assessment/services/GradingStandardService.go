package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
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
	scale, err := repositories.GetScaleByNameAndSchoolNumber(aggregateScaleRequest.ScaleName, aggregateScaleRequest.SchoolNumber)
	if err != nil {
		return errors.New("Checking scale name error: " + err.Error())
	}

	if scale.Number != "" {
		return errors.New("Name already used: " + aggregateScaleRequest.ScaleName)
	}

	grades := aggregateScaleRequest.Grades

	if len(grades) != aggregateScaleRequest.Count {
		return errors.New("you should define the grades up to set count")
	}

	err = checkDuplicates(grades)
	if err != nil {
		return err
	}

	trn := config.DB.Begin()

	scale, err = repositories.SaveScale(models2.Scale{
		Number:       uuid.NewV4().String(),
		Name:         aggregateScaleRequest.ScaleName,
		SchoolNumber: aggregateScaleRequest.SchoolNumber,
		Count:        aggregateScaleRequest.Count,
		UserNumber:   currentUser,
	}, trn)

	if err != nil {
		return errors.New("creating base scale error: " + err.Error())
	}

	for _, grade := range grades {

		err = createGrade(trn, grade, scale.Number)
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

	err := CheckOverlaps(grades)
	if err != nil {
		return err
	}

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

func Overlap(aggregateScaleItem1, aggregateScaleItem2 requests.AggregateScaleItem) bool {
	return aggregateScaleItem1.MinValue <= aggregateScaleItem2.MaxValue && aggregateScaleItem2.MinValue <= aggregateScaleItem1.MaxValue
}

func CheckOverlaps(aggregateScaleItems []requests.AggregateScaleItem) error {

	scaleLength := len(aggregateScaleItems)

	for i := 0; i < scaleLength; i++ {

		for j := i + 1; j < scaleLength; j++ {

			if Overlap(aggregateScaleItems[i], aggregateScaleItems[j]) {
				return errors.New("the scale has overlapping entries")
			}

		}

	}

	return nil
}

func createGrade(trn *gorm.DB, grade requests.AggregateScaleItem, scaleNumber string) error {

	checkDuplicate, err := repositories.GetAggregateScaleByScaleNumberAndCode(scaleNumber, grade.Code, trn)
	if err != nil {
		return errors.New("checking code duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used this code")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberAndAssessmentGradeNumber(scaleNumber, grade.GradeNumber, trn)
	if err != nil {
		return errors.New("checking grade duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used grade in your scale")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberAndWeight(scaleNumber, grade.Weight, trn)
	if err != nil {
		return errors.New("checking weight duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used weight on this scale")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberAndRemarks(scaleNumber, grade.Remark, trn)
	if err != nil {
		return errors.New("checking remark duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used remark on this scale")
	}

	newAggregateScale := models.AggregateScale{
		Weight:                grade.Weight,
		AssessmentGradeNumber: grade.GradeNumber,
		ScaleNumber:           scaleNumber,
		Code:                  grade.Code,
		Number:                uuid.NewV4().String(),
		MinValue:              grade.MinValue,
		MaxValue:              grade.MaxValue,
	}

	err = repositories.SaveAggregateScale(newAggregateScale, trn)
	if err != nil {
		return errors.New("Saving grading scale error: " + err.Error())
	}

	return nil
}

func FetchAggregateScales(campusClassNumber, campusNumber, gradingStandardNumber, assessmentAreaNumber string) ([]responses.AggregateScaleResponse, error) {

	var response []responses.AggregateScaleResponse

	classNumber := campusClassNumber

	if campusClassNumber != "" {
		class, err := repositories2.GetClassByCampusClassNumber(campusClassNumber)
		if err != nil {
			return response, errors.New("checking class error: " + err.Error())
		}

		if class.Number == "" {
			return response, errors.New("invalid class campus number")
		}

		classNumber = class.Number
	}

	aggregateScales, err := repositories.FindAggregateScaleByClassNumberAndGradingStandardNumber(classNumber, gradingStandardNumber)
	if err != nil {
		return response, errors.New("getting aggregate scale error: " + err.Error())
	}

	for _, aggregateScale := range aggregateScales {
		assessmentGrade, err := repositories.GetAssessmentGradeByNumber(aggregateScale.AssessmentGradeNumber)
		if err != nil {
			return response, errors.New("getting assessment grade error: " + err.Error())
		}

		scale, err := repositories.GetScaleByNumber(aggregateScale.ScaleNumber)
		if err != nil {
			return response, errors.New("getting scale error: " + err.Error())
		}

		response = append(response, createAggregateScaleResponse(aggregateScale, assessmentGrade, scale))
	}

	return response, nil
}

func createAggregateScaleResponse(aggregateScale models.AggregateScale, assessmentGrade models2.AssessmentGrade, scale models2.Scale) responses.AggregateScaleResponse {
	//baseURL + aggregateScale.Illustration
	return responses.AggregateScaleResponse{
		Code:         aggregateScale.Code,
		Number:       aggregateScale.Number,
		MinValue:     aggregateScale.MinValue,
		MaxValue:     aggregateScale.MaxValue,
		Weight:       aggregateScale.Weight,
		Illustration: "",
		Name:         assessmentGrade.Name,
		ScaleName:    scale.Name,
	}
}
