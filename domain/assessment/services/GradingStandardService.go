package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentGradeRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/request"
	assessmentGradeResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/responses"
	assessmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	classResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
)

func SaveGradingStandard(gradingStandardRequest requests.GradingStandardRequest, currentUser string) error {

	selectedClasses := gradingStandardRequest.Classes

	err := checkDuplicatesInGradingStandardClassesRequest(selectedClasses)
	if err != nil {
		return err
	}

	err = config.DB.Transaction(func(trn *gorm.DB) error {

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

		checkName, err := repositories.FindGradingStandardByNameAndSchoolNumber(gradingStandardRequest.Name, gradingStandardRequest.SchoolNumber, trn)
		if err != nil {
			return errors.New("Checking grading standard by name error: " + err.Error())
		}

		if checkName.Number != "" {
			return errors.New("name already used")
		}

		checkName, err = repositories.FindGradingStandardByCodeAndSchoolNumber(gradingStandardRequest.Code, gradingStandardRequest.SchoolNumber, trn)
		if err != nil {
			return errors.New("Checking grading standard by code error: " + err.Error())
		}

		if checkName.Number != "" {
			return errors.New("code already used")
		}

		err = repositories.SaveGradingStandard(gradingStandard, trn)
		if err != nil {
			trn.Rollback()
			return errors.New("Saving grading standard error: " + err.Error())
		}

		for _, gradingStandardClassRequest := range selectedClasses {
			err = processGradingStandardClass(gradingStandardClassRequest, gradingStandard.Number, trn)
			if err != nil {
				trn.Rollback()
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errors.New("commiting all changes error: " + err.Error())
	}

	return nil
}

func processGradingStandardClass(gradingStandardClassRequest requests.GradingStandardClasses, gradingStandardNumber string, trn *gorm.DB) error {
	class, err := classRepositories.GetClassByNumber(gradingStandardClassRequest.ClassNumber, trn)
	if err != nil {
		return errors.New("getting class from campus class number error: " + err.Error())
	}

	if class.Number == "" {
		return errors.New("campus class number has no class")
	}

	classNewScale, err := repositories.GetScaleByNumber(gradingStandardClassRequest.ScaleNumber, trn)
	if err != nil {
		return errors.New("error checking scale: " + err.Error())
	}

	if classNewScale.Number == "" {
		return errors.New("invalid scale number")
	}

	checkClassAndScale, err := repositories.GetGradingStandardClassScaleByClassNumberAndScaleNumber(class.Number, gradingStandardClassRequest.ScaleNumber, trn)
	if err != nil {
		return errors.New("checking class and scale duplicate error: " + err.Error())
	}

	if checkClassAndScale.ClassNumber != "" && checkClassAndScale.GradingStandardNumber != gradingStandardNumber {
		return errors.New(classNewScale.Name + " is already set to " + class.Name)
	}

	checkClassAndScale, err = repositories.GetGradingStandardClassScaleByClassNumberAndScaleNumber(class.Number, "", trn)
	if err != nil {
		return errors.New("checking only class duplicate error: " + err.Error())
	}

	if checkClassAndScale.ClassNumber != "" {

		classCurrentScale, _ := repositories.GetScaleByNumber(checkClassAndScale.ScaleNumber, trn)

		if classCurrentScale.Count != classNewScale.Count {
			return errors.New(classNewScale.Name + " can not be used on " + class.Name + ", they differ in size with " + classCurrentScale.Name)
		}

	}

	checkClassAndScale, err = repositories.GetGradingStandardClassScaleByClassNumberAndScaleNumber(class.Number, gradingStandardClassRequest.ScaleNumber, trn)
	if err != nil {
		return errors.New("checking class duplicate error: " + err.Error())
	}

	if checkClassAndScale.ClassNumber == "" {
		return errors.New("'" + classNewScale.Name + "' is already added to class: " + class.Name)
	}

	err = repositories.SaveGradingStandardClassScale(assessmentGradeModels.GradingStandardClassScale{
		ClassNumber:           class.Number,
		ScaleNumber:           gradingStandardClassRequest.ScaleNumber,
		GradingStandardNumber: gradingStandardNumber,
	}, trn)

	if err != nil {
		return errors.New("Saving grading standard classes error: " + err.Error())
	}

	return nil
}

func UpdateGradingStandard(gradingStandardRequest requests.UpdateGradingStandardRequest, currentUser string) error {

	selectedClasses := gradingStandardRequest.Classes

	err := checkDuplicatesInGradingStandardClassesRequest(selectedClasses)
	if err != nil {
		return err
	}

	trn := config.DB.Begin()

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardRequest.GradingStandardNumber, trn)
	if err != nil {
		return errors.New("Checking grading standard error: " + err.Error())
	}

	checkName, err := repositories.FindGradingStandardByNameAndSchoolNumber(gradingStandardRequest.Name, gradingStandardRequest.SchoolNumber, trn)
	if err != nil {
		return errors.New("Checking grading standard by name error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != gradingStandard.Number {
		return errors.New("name already used on another standard")
	}

	checkName, err = repositories.FindGradingStandardByCodeAndSchoolNumber(gradingStandardRequest.Code, gradingStandardRequest.SchoolNumber, trn)
	if err != nil {
		return errors.New("Checking grading standard by code error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != gradingStandard.Number {
		return errors.New("code already used on another standard")
	}

	oldMappings, err := repositories.FindGradingStandardClassScaleByGradingStandardNumber(gradingStandard.Number, trn)
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

	for _, oldMapping := range oldMappings {
		err = repositories.RemoveGradingStandardClassScale(oldMapping, trn)
		if err != nil {
			return errors.New("clearing old mapping error: " + err.Error())
		}
	}

	err = repositories.SaveGradingStandard(gradingStandard, trn)
	if err != nil {
		return errors.New("Updating grading standard error: " + err.Error())
	}

	for _, gradingStandardClassRequest := range selectedClasses {
		err = processGradingStandardClass(gradingStandardClassRequest, gradingStandard.Number, trn)
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

	err = config.DB.Transaction(func(trn *gorm.DB) error {

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
				return err
			}
		}

		return nil
	})

	if err != nil {
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
	var totalRanges float64

	for i, grade := range grades {
		totalRanges += (grade.MaxValue - grade.MinValue) + 1

		shouldReject := ""

		for _, cleanGrade := range cleanGrades {

			if cleanGrade.GradeNumber == grade.GradeNumber {
				gr, err := repository.FindGradeByNumber(grade.GradeNumber)
				if err != nil {
					return errors.New("error fetching grade: " + err.Error())
				}
				shouldReject = fmt.Sprintf("grade %s already exists", gr.Name)

			}

		}

		if shouldReject == "" {
			cleanGrades[i] = grade
		} else {
			return errors.New(shouldReject)
		}

	}

	if (totalRanges - 1) != 100 {
		return errors.New("please ensure there are no ranges skipped ,the total range should add up to 100")
	}

	return nil
}

func checkDuplicatesInGradingStandardClassesRequest(selectedClasses []requests.GradingStandardClasses) error {

	cleanClassNumbers := make([]requests.GradingStandardClasses, len(selectedClasses))

	for _, selectedClass := range selectedClasses {

		for _, cleanClassNumber := range cleanClassNumbers {

			if cleanClassNumber.ClassNumber == selectedClass.ClassNumber {
				class, err := classRepositories.GetClassByNumber(selectedClass.ClassNumber)
				if err != nil {
					return errors.New("error fetching grade: " + err.Error())
				}
				return errors.New(fmt.Sprintf("%s is selected multiple times", class.Name))
			}

		}

		cleanClassNumbers = append(cleanClassNumbers, selectedClass)

	}

	return nil
}

func checkDuplicatesInClasses(selectedClasses []string) error {

	cleanClassNumbers := make([]string, len(selectedClasses))

	for _, selectedClass := range selectedClasses {

		for _, cleanClassNumber := range cleanClassNumbers {

			if cleanClassNumber == selectedClass {
				class, err := classRepositories.GetClassByNumber(selectedClass)
				if err != nil {
					return errors.New("error fetching grade: " + err.Error())
				}
				return errors.New(fmt.Sprintf("%s is selected multiple times", class.Name))
			}

		}

		cleanClassNumbers = append(cleanClassNumbers, selectedClass)

	}

	return nil
}

func Overlap(aggregateScaleItem1, aggregateScaleItem2 requests.AggregateScaleItem) bool {
	return aggregateScaleItem1.MinValue <= aggregateScaleItem2.MaxValue && aggregateScaleItem2.MinValue <= aggregateScaleItem1.MaxValue
}

func OverlapInGrades(aggregateScaleItem1, aggregateScaleItem2 assessmentGradeRequests.AssessmentGradeItemRequest) bool {
	return aggregateScaleItem1.MinimumPoints <= aggregateScaleItem2.MaximumPoints && aggregateScaleItem2.MinimumPoints <= aggregateScaleItem1.MaximumPoints
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

func CheckOverlapsInGrades(assessmentGradeItemRequests []assessmentGradeRequests.AssessmentGradeItemRequest) error {

	scaleLength := len(assessmentGradeItemRequests)

	for i := 0; i < scaleLength; i++ {

		for j := i + 1; j < scaleLength; j++ {

			if OverlapInGrades(assessmentGradeItemRequests[i], assessmentGradeItemRequests[j]) {
				return errors.New("the scale has overlapping entry grades")
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

func FetchGradingStandards(schoolNumber, classNumber, scaleNumber string) ([]responses.GradingStandardResponse, error) {
	var gradingStandardResponses = make([]responses.GradingStandardResponse, 0)

	gradingStandards, err := repositories.FindGradingStandardByCampusNumberAndCampusClassNumber(schoolNumber, classNumber, scaleNumber)
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

func GetGradingStandardDetails(gradingStandardNumber string) (responses.DetailedGradingStandardResponse, error) {

	var gradingStandardResponse responses.DetailedGradingStandardResponse

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardNumber)
	if err != nil {
		return gradingStandardResponse, errors.New("checking grading standard number err: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return gradingStandardResponse, errors.New("invalid grading standard number")
	}

	gradingStandardResponse.GradingStandardResponse = createGradingStandardResponse(gradingStandard)

	classes, err := repositories.FindGradingStandardClassScaleByGradingStandardNumber(gradingStandard.Number)

	var gradingStandardClassesResponses = make([]responses.GradingStandardClassesResponse, 0)

	for _, class := range classes {

		gradingClassResponse, err := createGradingStandardClassesResponse(class)
		if err != nil {
			return gradingStandardResponse, err
		}

		gradingStandardClassesResponses = append(gradingStandardClassesResponses, gradingClassResponse)
	}

	gradingStandardResponse.StandardClasses = gradingStandardClassesResponses

	return gradingStandardResponse, nil
}

func createGradingStandardClassesResponse(gradingScaleClass assessmentGradeModels.GradingStandardClassScale) (responses.GradingStandardClassesResponse, error) {

	var gradingStandardClassesResponse responses.GradingStandardClassesResponse

	scale, err := repositories.GetScaleByNumber(gradingScaleClass.ScaleNumber)
	if err != nil {
		return gradingStandardClassesResponse, errors.New("getting scale details error: " + err.Error())
	}

	if scale.Number == "" {
		return gradingStandardClassesResponse, errors.New("scale not determined")
	}

	aggregateScales, err := repositories.FindAggregateScaleByScaleNumber(scale.Number)
	if err != nil {
		return gradingStandardClassesResponse, errors.New("getting aggregate scales error: " + err.Error())
	}

	detailedScale, err := createScaleResponse(scale, aggregateScales)
	if err != nil {
		return gradingStandardClassesResponse, err
	}

	class, err := classRepositories.GetClassByNumber(gradingScaleClass.ClassNumber)
	if err != nil {
		return gradingStandardClassesResponse, errors.New("getting classes error: " + err.Error())
	}

	if class.Number == "" {
		return gradingStandardClassesResponse, errors.New("could not determine class by number")
	}

	gradingStandardClassesResponse.Scale = detailedScale

	gradingStandardClassesResponse.Class = services.CreateClassResponse(class)

	return gradingStandardClassesResponse, nil
}

func SaveAssessmentGrade(assessmentGradeRequest assessmentGradeRequests.AssessmentGradeRequest, currentUserNumber string) error {

	err := checkDuplicatesInClasses(assessmentGradeRequest.Classes)
	if err != nil {
		return err
	}

	err = CheckOverlapsInGrades(assessmentGradeRequest.AssessmentGradeItems)
	if err != nil {
		return errors.New(err.Error())
	}

	err = config.DB.Transaction(func(trn *gorm.DB) error {

		saveAssessmentGrade := assessmentGradeModels.AssessmentGrade{
			Status:              constants.StatusActive,
			Number:              uuid.NewV4().String(),
			UserNumber:          currentUserNumber,
			AssessmentAreaCount: assessmentGradeRequest.AssessmentAreaCount,
			SchoolNumber:        assessmentGradeRequest.SchoolNumber,
		}

		err = repositories.SaveAssessmentGrade(saveAssessmentGrade, trn)
		if err != nil {
			return errors.New("saving assessment grade error: " + err.Error())
		}

		for _, classNumber := range assessmentGradeRequest.Classes {
			err = processAssessmentGradeClasses(classNumber, saveAssessmentGrade, trn)
			if err != nil {
				return err
			}
		}

		for _, assessmentGradeScaleItem := range assessmentGradeRequest.AssessmentGradeItems {
			err = processAssessmentGradeScaleItem(assessmentGradeScaleItem, saveAssessmentGrade.Number, trn)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errors.New("failed finalizing assessment grade, error: " + err.Error())
	}

	return nil
}

func processAssessmentGradeClasses(classNumber string, saveAssessmentGrade assessmentGradeModels.AssessmentGrade, trn *gorm.DB) error {

	class, err := classRepositories.GetClassByNumber(classNumber, trn)
	if err != nil {
		return errors.New("checking class error: " + err.Error())
	}

	if class.Number == "" {
		return errors.New("invalid class number selected")
	}

	checkClassDuplicate, err := repositories.GetAssessmentGradeClassByGradeNumberAndClassNumber(saveAssessmentGrade.Number, class.Number, trn)
	if err != nil {
		return errors.New("checking class double entry error: " + err.Error())
	}

	if checkClassDuplicate.ClassNumber != "" {
		return errors.New("Class '" + class.Name + "' already selected")
	}

	checkWithClass, err := repositories.GetAssessmentGradeByClassNumberAndStatus(class.Number, constants.StatusActive)
	if err != nil {
		return errors.New("checking class with active assessment grade error: " + err.Error())
	}

	if checkWithClass.Number != "" && checkWithClass.Number != saveAssessmentGrade.Number {
		return errors.New("Class '" + class.Name + "' already has an active assessment grade")
	}

	err = repositories.SaveAssessmentGradeClass(assessmentModel.AssessmentGradeClass{
		AssessmentGradeNumber: saveAssessmentGrade.Number,
		ClassNumber:           class.Number,
	}, trn)

	if err != nil {
		return errors.New("attaching class error: " + err.Error())
	}

	return nil
}

func processAssessmentGradeScaleItem(assessmentGradeScaleItem assessmentGradeRequests.AssessmentGradeItemRequest, assessmentGradeNumber string, trn *gorm.DB) error {
	checkDuplicateEntries, err := repositories.GetAssessmentGradeScaleByCodeAndGradeNumber(assessmentGradeScaleItem.Code, assessmentGradeNumber, trn)
	if err != nil {
		return errors.New("checking assessment grades error : " + err.Error())
	}

	if checkDuplicateEntries.Number != "" {
		return errors.New("code already used")
	}

	checkDuplicateEntries, err = repositories.GetAssessmentGradeScaleByNameAndGradeNumber(assessmentGradeScaleItem.Name, assessmentGradeNumber, trn)
	if err != nil {
		return errors.New("checking assessment grades name error : " + err.Error())
	}

	if checkDuplicateEntries.Number != "" {
		return errors.New("name " + assessmentGradeScaleItem.Name + " already used")
	}

	err = repositories.SaveAssessmentGradeScale(assessmentGradeModels.AssessmentGradeScale{
		Code:                  assessmentGradeScaleItem.Code,
		AssessmentGradeNumber: assessmentGradeNumber,
		MaximumPoints:         assessmentGradeScaleItem.MaximumPoints,
		MinimumPoints:         assessmentGradeScaleItem.MinimumPoints,
		Description:           assessmentGradeScaleItem.Description,
		Name:                  assessmentGradeScaleItem.Name,
		Remark:                assessmentGradeScaleItem.Remarks,
		Abbreviation:          assessmentGradeScaleItem.Abbreviation,
		Number:                uuid.NewV4().String(),
	}, trn)

	if err != nil {
		return errors.New("failed saving assessment grade scale item: " + err.Error())
	}

	return nil
}

func GetAssessmentGrades(classNumber, status string) ([]assessmentGradeResponses.AssessmentGradeResponse, error) {

	var assessmentGradeResponse = make([]assessmentGradeResponses.AssessmentGradeResponse, 0)

	assessmentGrades, err := repositories.FindAssessmentGradeByClassNumber(classNumber, status)
	if err != nil {
		return assessmentGradeResponse, errors.New("checking assessment grades by scale error: " + err.Error())
	}

	for _, assessmentGrade := range assessmentGrades {
		gradeScale, err := repositories.FindAssessmentGradeScaleByGradeNumber(assessmentGrade.Number)
		if err != nil {
			return assessmentGradeResponse, errors.New("checking assessment scale error: " + err.Error())
		}

		classes, err := classRepositories.FindClassesByAssessmentGradeNumber(assessmentGrade.Number)
		if err != nil {
			return assessmentGradeResponse, errors.New("checking assessment grade classes error: " + err.Error())
		}

		assessmentGradeResponse = append(assessmentGradeResponse, createAssessmentGradeResponse(assessmentGrade, gradeScale, classes))
	}

	return assessmentGradeResponse, nil
}

func createAssessmentGradeResponse(assessmentGrade assessmentGradeModels.AssessmentGrade, assessmentGradeScales []assessmentGradeModels.AssessmentGradeScale, classes []classModels.Class) assessmentGradeResponses.AssessmentGradeResponse {

	assessmentGradeResponse := assessmentGradeResponses.AssessmentGradeResponse{
		AssessmentAreaCount:   assessmentGrade.AssessmentAreaCount,
		AssessmentGradeNumber: assessmentGrade.Number,
	}

	var assessmentGradeItemResponses = make([]assessmentGradeResponses.AssessmentGradeItemResponse, 0)

	for _, assessmentGradeItem := range assessmentGradeScales {
		assessmentGradeItemResponses = append(assessmentGradeItemResponses, assessmentGradeResponses.AssessmentGradeItemResponse{
			Name:          assessmentGradeItem.Name,
			Abbreviation:  assessmentGradeItem.Abbreviation,
			MaximumPoints: assessmentGradeItem.MaximumPoints,
			MinimumPoints: assessmentGradeItem.MinimumPoints,
			Code:          assessmentGradeItem.Code,
			Description:   assessmentGradeItem.Description,
			Remarks:       assessmentGradeItem.Remark,
		})
	}

	var classResponse = make([]classResponses.ClassResponse, 0)

	for _, class := range classes {
		classResponse = append(classResponse, services.CreateClassResponse(class))
	}

	assessmentGradeResponse.Classes = classResponse
	assessmentGradeResponse.AssessmentGradeItems = assessmentGradeItemResponses
	return assessmentGradeResponse
}

func UpdateGradingStandardStatus(updateStatusRequest requests.UpdateGradingStatusStandardRequest, currentUserNumber string) error {

	gradingStandard, err := repositories.GetGradingStandardByNumber(updateStatusRequest.GradingStandardNumber)
	if err != nil {
		return errors.New("failed getting grading standard by number: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return errors.New("invalid grading standard number")
	}

	gradingStandard.Status = updateStatusRequest.Status
	gradingStandard.UserNumber = currentUserNumber
	err = repositories.SaveGradingStandard(gradingStandard)
	if err != nil {
		return errors.New("failed updating grading standard: " + err.Error())
	}

	return nil
}

func RemoveAssessmentGrade(assessmentGradeNumber string) error {
	assessmentGrade, err := repositories.FindAssessmentGradeByNumber(assessmentGradeNumber)
	if err != nil {
		return errors.New("checking assessment grade error: " + err.Error())
	}

	if assessmentGrade.Number == "" {
		return errors.New("invalid assessment grade number")
	}

	scales, err := repositories.FindAssessmentGradeScaleByGradeNumber(assessmentGrade.Number)
	if err != nil {
		return errors.New("checking assessment grade scale error: " + err.Error())
	}

	err = config.DB.Transaction(func(trn *gorm.DB) error {

		for _, scale := range scales {
			err = repositories.RemoveAssessmentGradeScale(scale, trn)
			if err != nil {
				return errors.New("clearing assessment grade scale error: " + err.Error())
			}
		}

		err = repositories.RemoveAssessmentGrade(assessmentGrade, trn)
		if err != nil {
			return errors.New("clearing assessment grade error: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return errors.New("removing assessment grade error: " + err.Error())
	}

	return nil
}

func EditAssessmentGrade(assessmentGradeRequest assessmentGradeRequests.EditAssessmentGradeRequest, currentUser string) error {

	assessmentGrade, err := repositories.FindAssessmentGradeByNumber(assessmentGradeRequest.AssessmentGradeNumber)
	if err != nil {
		return errors.New("checking assessment grade error: " + err.Error())
	}

	if assessmentGrade.Number == "" {
		return errors.New("invalid assessment grade number")
	}

	checkClassesError := checkDuplicatesInClasses(assessmentGradeRequest.Classes)
	if checkClassesError != nil {
		return checkClassesError
	}

	scales, err := repositories.FindAssessmentGradeScaleByGradeNumber(assessmentGrade.Number)
	if err != nil {
		return errors.New("checking assessment grade scale error: " + err.Error())
	}

	classes, err := repositories.FindAssessmentGradeClassByGradeNumber(assessmentGrade.Number)
	if err != nil {
		return errors.New("checking assessment grade classes error: " + err.Error())
	}

	if assessmentGradeRequest.Status != "" {

		if assessmentGradeRequest.Status != constants.StatusActive && assessmentGradeRequest.Status != constants.StatusInactive {
			return errors.New("the status should be ACTIVE or INACTIVE")
		}

		assessmentGrade.Status = assessmentGradeRequest.Status
	}

	assessmentGrade.UserNumber = currentUser
	assessmentGrade.AssessmentAreaCount = assessmentGradeRequest.AssessmentAreaCount

	err = config.DB.Transaction(func(trn *gorm.DB) error {

		err = repositories.SaveAssessmentGrade(assessmentGrade, trn)

		if err != nil {
			return errors.New("updating assessment grade meta error: " + err.Error())
		}

		for _, class := range classes {
			err = repositories.RemoveAssessmentGradeClass(class, trn)
			if err != nil {
				return errors.New("clearing previous class error: " + err.Error())
			}
		}

		for _, scale := range scales {
			err = repositories.RemoveAssessmentGradeScale(scale, trn)
			if err != nil {
				return errors.New("clearing previous scale error: " + err.Error())
			}
		}

		for _, newClass := range assessmentGradeRequest.Classes {
			err = processAssessmentGradeClasses(newClass, assessmentGrade, trn)
			if err != nil {
				return err
			}
		}

		for _, newScale := range assessmentGradeRequest.AssessmentGradeItems {
			err = processAssessmentGradeScaleItem(newScale, assessmentGrade.Number, trn)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errors.New("updating assessment grade error: " + err.Error())
	}

	return nil
}
