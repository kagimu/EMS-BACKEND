package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
)

func SaveGradingStandard(gradingStandardRequest requests.GradingStandardRequest, currentUser string) error {

	trn := config.DB.Begin()

	gradingStandard := models.GradingStandard{
		Code:         gradingStandardRequest.Code,
		Number:       uuid.NewV4().String(),
		Name:         gradingStandardRequest.Name,
		Description:  gradingStandardRequest.Description,
		UserNumber:   currentUser,
		SchoolNumber: gradingStandardRequest.SchoolNumber,
	}

	if gradingStandardRequest.Status != "" {
		gradingStandard.Status = gradingStandardRequest.Status
	}

	err := repositories.SaveGradingStandard(gradingStandard, trn)
	if err != nil {
		trn.Rollback()
		return errors.New("Saving grading standard error: " + err.Error())
	}

	for _, classNumber := range gradingStandardRequest.Classes {
		err = processGradingStandardClass(classNumber, gradingStandard.Number, trn)
		if err != nil {
			trn.Rollback()
			return err
		}
	}

	if err := trn.Commit().Error; err != nil {
		return errors.New("commiting all changes error: " + err.Error())
	}

	return nil
}

func processGradingStandardClass(gradingStandardClass requests.GradingStandardClasses, gradingStandardNumber string, trn *gorm.DB) error {
	class, err := repositories2.GetClassByCampusClassNumber(gradingStandardClass.CampusClassNumber)
	if err != nil {
		trn.Rollback()
		return errors.New("getting class from campus class number error: " + err.Error())
	}

	if class.Number == "" {
		return errors.New("campus class number has no class")
	}

	duplicateClass, err := repositories.GetGradingStandardClassScaleByClassNumberAndScaleNumber(class.Number, gradingStandardClass.ScaleNumber, trn)
	if err != nil {
		return errors.New("checking class duplicate error: " + err.Error())
	}

	if duplicateClass.ClassNumber != "" {
		return errors.New("you already selected class: " + class.Name)
	}

	classGrade, err := repositories.GetGradingStandardClassScaleByClassNumberAndStandard(class.Number, gradingStandardNumber)
	if err != nil {
		return errors.New("checking class duplicate error: " + err.Error())
	}

	if classGrade.ID > 0 {
		err = repositories.RemoveGradingStandardClassScale(classGrade, trn)
		if err != nil {
			return errors.New("clearing class duplicate error: " + err.Error())
		}
	}

	err = repositories.SaveGradingStandardClassScale(models2.GradingStandardClassScale{
		ClassNumber:           class.Number,
		ScaleNumber:           gradingStandardClass.ScaleNumber,
		GradingStandardNumber: gradingStandardNumber,
	}, trn)

	if err != nil {
		trn.Rollback()
		return errors.New("Saving grading standard classes error: " + err.Error())
	}

	return nil
}

func UpdateGradingStandard(gradingStandardRequest requests.UpdateGradingStandardRequest, currentUser string) error {

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardRequest.GradingStandardNumber)
	if err != nil {
		return errors.New("Checking grading standard error: " + err.Error())
	}

	oldMappings, err := repositories.FindGradingStandardClassScaleByGradingStandardNumber(gradingStandard.Number)
	if err != nil {
		return errors.New("Checking grading standard old mapping error: " + err.Error())
	}

	gradingStandard.Code = gradingStandardRequest.Code
	gradingStandard.Name = gradingStandardRequest.Name
	gradingStandard.Description = gradingStandardRequest.Description
	gradingStandard.UserNumber = currentUser

	if gradingStandardRequest.Status != "" {
		gradingStandard.Status = gradingStandardRequest.Status
	}

	trn := config.DB.Begin()

	for _, oldMapping := range oldMappings {
		err = repositories.RemoveGradingStandardClassScale(oldMapping)
		if err != nil {
			return errors.New("clearing old mapping error: " + err.Error())
		}
	}

	err = repositories.SaveGradingStandard(gradingStandard, trn)
	if err != nil {
		return errors.New("Saving grading standard error: " + err.Error())
	}

	for _, class := range gradingStandardRequest.Classes {
		err = processGradingStandardClass(class, gradingStandard.Number, trn)
		if err != nil {
			return err
		}
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

	if int64(len(grades)) != aggregateScaleRequest.Count {
		return errors.New("you should define the grades up to set count")
	}

	err = checkDuplicates(grades)
	if err != nil {
		return err
	}

	trn := config.DB.Begin()

	scale, err = repositories.SaveScale(models.Scale{
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

		err = createAggregateScales(trn, grade, scale.Number)
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

func createAggregateScales(trn *gorm.DB, aggregateScaleItem requests.AggregateScaleItem, scaleNumber string) error {

	checkDuplicate, err := repositories.GetAggregateScaleByScaleNumberAndCode(scaleNumber, aggregateScaleItem.Code, trn)
	if err != nil {
		return errors.New("checking code duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used this code")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberGradeNumber(scaleNumber, aggregateScaleItem.GradeNumber, trn)
	if err != nil {
		return errors.New("checking grade duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used grade in your scale")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberAndWeight(scaleNumber, aggregateScaleItem.Weight, trn)
	if err != nil {
		return errors.New("checking weight duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used weight on this scale")
	}

	checkDuplicate, err = repositories.GetAggregateScaleByScaleNumberAndDescription(scaleNumber, aggregateScaleItem.Remark, trn)
	if err != nil {
		return errors.New("checking remark duplicate error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("you have already used remark on this scale")
	}

	newAggregateScale := models.AggregateScale{
		Weight:      aggregateScaleItem.Weight,
		GradeNumber: aggregateScaleItem.GradeNumber,
		ScaleNumber: scaleNumber,
		Code:        aggregateScaleItem.Code,
		Number:      uuid.NewV4().String(),
		MinValue:    aggregateScaleItem.MinValue,
		MaxValue:    aggregateScaleItem.MaxValue,
		Description: aggregateScaleItem.Remark,
	}

	err = repositories.SaveAggregateScale(newAggregateScale, trn)
	if err != nil {
		return errors.New("Saving grading scale error: " + err.Error())
	}

	return nil
}

func FetchGradingStandards(schoolNumber string) ([]responses.GradingStandardResponse, error) {
	var gradingStandardResponses = make([]responses.GradingStandardResponse, 0)

	gradingStandards, err := repositories.FindGradingStandardByCampusNumberAndCampusClassNumber(schoolNumber)
	if err != nil {
		return gradingStandardResponses, errors.New("getting grading standard error: " + err.Error())
	}

	for _, gradingStandard := range gradingStandards {
		gradingStandardResponses = append(gradingStandardResponses, createGradingStandardResponse(gradingStandard))
	}

	return gradingStandardResponses, nil
}

func createGradingStandardResponse(gradingStandard models.GradingStandard) responses.GradingStandardResponse {
	return responses.GradingStandardResponse{
		Number:      gradingStandard.Number,
		Name:        gradingStandard.Name,
		Code:        gradingStandard.Code,
		Description: gradingStandard.Description,
		Status:      gradingStandard.Status,
	}
}

func FindAllDetailsScales(schoolNumber string) ([]responses.DetailedScaleResponse, error) {

	var scaleResponses = make([]responses.DetailedScaleResponse, 0)

	scales, err := repositories.FindAllScales(schoolNumber)
	if err != nil {
		return scaleResponses, errors.New("error fetching scales: " + err.Error())
	}

	for _, scale := range scales {

		aggregateScale, err := repositories.FindAggregateScaleByScaleNumber(scale.Number)
		if err != nil {
			return scaleResponses, errors.New("error fetching aggregate scales : " + err.Error())
		}

		newResponse, err := createScaleResponse(scale, aggregateScale)
		if err != nil {
			return scaleResponses, err
		}

		scaleResponses = append(scaleResponses, newResponse)
	}

	return scaleResponses, nil
}

func createScaleResponse(scale models.Scale, aggregateScales []models.AggregateScale) (responses.DetailedScaleResponse, error) {

	minResponse := createMinResponse(scale)

	response := responses.DetailedScaleResponse{}
	response.ScaleResponse = minResponse

	var grades = make([]responses.AggregateScaleResponse, 0)

	for _, aggregateScale := range aggregateScales {

		grade, err := repository.FindGradeByNumber(aggregateScale.GradeNumber)
		if err != nil {
			return response, errors.New("error fetching aggregate scales : " + err.Error())
		}

		if grade.Number == "" {
			return response, errors.New("could not get grade by number: " + aggregateScale.GradeNumber)
		}

		grades = append(grades, responses.AggregateScaleResponse{
			MinValue:     aggregateScale.MinValue,
			MaxValue:     aggregateScale.MaxValue,
			Description:  aggregateScale.Description,
			Illustration: "",
			Weight:       aggregateScale.Weight,
			Name:         grade.Name,
			ShortName:    grade.ShortName,
			Status:       grade.Status,
			GradeNumber:  grade.Number,
			Code:         aggregateScale.Code,
		})
	}

	response.Grades = grades

	return response, nil
}

func createMinResponse(scale models.Scale) responses.ScaleResponse {
	return responses.ScaleResponse{
		ScaleNumber: scale.Number,
		ScaleName:   scale.Name,
		Count:       scale.Count,
	}
}

func FindAllScales(schoolNumber string) ([]responses.ScaleResponse, error) {

	var scaleResponses = make([]responses.ScaleResponse, 0)

	scales, err := repositories.FindAllScales(schoolNumber)
	if err != nil {
		return scaleResponses, errors.New("error fetching scales: " + err.Error())
	}

	for _, scale := range scales {
		scaleResponses = append(scaleResponses, createMinResponse(scale))
	}

	return scaleResponses, nil

}

func UpdateAggregateScale(aggregateScaleRequest requests.EditAggregateScaleRequest, currentUser string) error {

	scale, err := repositories.GetScaleByNumber(aggregateScaleRequest.ScaleNumber)
	if err != nil {
		return errors.New("checking scale number error: " + err.Error())
	}

	if scale.Number == "" {
		return errors.New("invalid scale number")
	}

	grades := aggregateScaleRequest.Grades

	if int64(len(grades)) != aggregateScaleRequest.Count {
		return errors.New("you should define the grades up to set count")
	}

	err = checkDuplicates(grades)
	if err != nil {
		return err
	}

	aggregateScales, err := repositories.FindAggregateScaleByScaleNumber(scale.Number)
	if err != nil {
		return errors.New("getting old aggregate scales error: " + err.Error())
	}

	trn := config.DB.Begin()

	for _, aggregateScale := range aggregateScales {
		err = repositories.RemoveAggregateScale(aggregateScale, trn)
		if err != nil {
			return errors.New("clearing old scale error: " + err.Error())
		}
	}

	scale.Name = aggregateScaleRequest.ScaleName
	scale.Count = aggregateScaleRequest.Count
	scale.UserNumber = currentUser

	_, err = repositories.SaveScale(scale, trn)
	if err != nil {
		return errors.New("updating base scale error: " + err.Error())
	}

	for _, grade := range grades {

		err = createAggregateScales(trn, grade, scale.Number)
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

func RemoveGradingStandard(gradingStandardNumber string) error {

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardNumber)
	if err != nil {
		return errors.New("checking grading standard number err: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return errors.New("invalid grading standard number")
	}

	classes, err := repositories.FindGradingStandardClassScaleByGradingStandardNumber(gradingStandard.Number)

	trn := config.DB.Begin()

	for _, class := range classes {
		err = repositories.RemoveGradingStandardClassScale(class, trn)
		if err != nil {
			trn.Rollback()
			return errors.New("remove class mapping error: " + err.Error())
		}
	}

	err = repositories.RemoveGradingStandard(gradingStandard, trn)
	if err != nil {
		trn.Rollback()
		return errors.New("this standard can not be removed, try updating its status to " + constants.StatusInactive)
	}

	if err = trn.Commit().Error; err != nil {
		return errors.New("remove finalization failed, try again")
	}

	return nil
}

func RemoveScale(scaleNumber string) error {
	scale, err := repositories.GetScaleByNumber(scaleNumber)
	if err != nil {
		return errors.New("checking scale number failed: " + err.Error())
	}

	if scale.Number == "" {
		return errors.New("invalid scale number")
	}

	aggregateScales, err := repositories.FindAggregateScaleByScaleNumber(scale.Number)
	if err != nil {
		return errors.New("getting scales error: " + err.Error())
	}

	scaleClasses, err := repositories.FindGradingStandardClassScaleByScaleNumber(scale.Number)
	if err != nil {
		return errors.New("checking classes attached error: " + err.Error())
	}

	classCount := len(scaleClasses)
	if classCount > 0 {
		return errors.New("this scale is being used on " + strconv.Itoa(classCount) + " classes")
	}

	trn := config.DB.Begin()

	for _, aggregateScale := range aggregateScales {
		err = repositories.RemoveAggregateScale(aggregateScale, trn)
		if err != nil {
			trn.Rollback()
			return errors.New("clearing scales error: " + err.Error())
		}
	}

	err = repositories.RemoveScale(scale, trn)
	if err != nil {
		trn.Rollback()
		return errors.New("removing scale failed: " + err.Error())
	}

	if err = trn.Commit().Error; err != nil {
		return errors.New("finalization failed, try again: " + err.Error())
	}

	return nil
}
