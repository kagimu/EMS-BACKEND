package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	aggregateScaleRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	assessmentScaleModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	schoolRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	subjectRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func SaveAssessmentArea(payload requests.AssessmentAreaRequest, currentUser string) error {

	trn := config.DB.Begin()

	for _, subjectNumber := range payload.Subjects {
		err := processSaveAssessmentArea(subjectNumber, payload.CurriculumClassNumber, currentUser, trn)
		if err != nil {
			trn.Rollback()
			return err
		}
	}

	if err := trn.Commit().Error; err != nil {
		return errors.New("Failed to commit transaction, error: " + err.Error())
	}

	return nil
}

func processSaveAssessmentArea(subjectNumber string, curriculumClassNumber string, currentUser string, trn *gorm.DB) error {

	var otherName string

	subject, err := subjectRepositories.GetSubjectByNumber(subjectNumber, trn)
	if err != nil {
		return errors.New("checking subject error: " + err.Error())
	}

	if subject.Number == "" {
		otherName = subjectNumber
		subjectNumber = ""
	} else {
		otherName = ""
	}

	//validate for composite unique curriculum_class_number and subject number
	if subjectNumber != "" {
		assessmentArea, err := repositories.FindAssessmentAreaBySubjectNumberAndCurriculumClassNumber(subjectNumber, curriculumClassNumber, trn)
		if err != nil {
			return err
		}

		if assessmentArea.Number != "" {
			fmt.Println("assessment area already exists for the specified curriculum class and subject")
			return nil
		}
	}

	if otherName != "" {
		assessmentArea, err := repositories.FindAssessmentAreaByNameAndCurriculumClassNumber(otherName, curriculumClassNumber, trn)
		if err != nil {
			return err
		}

		if assessmentArea.Number != "" {
			fmt.Println("assessment area already exists for the specified curriculum class and subject")
			return nil
		}
	}

	saveAssessmentArea := models.AssessmentArea{
		Number:                uuid.NewV4().String(),
		CurriculumClassNumber: curriculumClassNumber,
		Name:                  otherName,
		UserNumber:            currentUser,
		SubjectNumber:         subjectNumber,
	}

	_, err = repositories.SaveAssessmentArea(saveAssessmentArea, trn)
	if err != nil {
		return errors.New("failed to create assessment area")
	}

	return nil

}

func UpdateAssessmentArea(payload requests.UpdateAssessmentAreaRequest, currentUser string) error {

	assessmentArea, err := repositories.GetAssessmentAreaByNumber(payload.AssessmentAreaNumber)
	if err != nil {
		return errors.New("checking assessment area error: " + err.Error())
	}

	if assessmentArea.Number == "" {
		return errors.New("invalid assessment area number")
	}

	assessmentByName, err := repositories.FindAssessmentAreaByNameAndCurriculumClassNumber(payload.Name, assessmentArea.CurriculumClassNumber)
	if err != nil {
		return errors.New("checking assessment area name error: " + err.Error())
	}

	if assessmentByName.Number != "" && assessmentByName.Number != assessmentArea.Number {
		return errors.New("assessment area name already used")
	}

	assessmentArea.Name = payload.Name
	assessmentArea.UserNumber = currentUser

	_, err = repositories.SaveAssessmentArea(assessmentArea)
	if err != nil {
		return errors.New("failed to update assessment area")
	}

	return nil
}

func FetchAssessmentAreasByGradingStandard(schoolNumber, assessmentAreaNumber, gradingStandardNumber string) ([]responses.GradingStandardAssessmentAreaResponse, error) {

	var assessmentAreaResponses = make([]responses.GradingStandardAssessmentAreaResponse, 0)

	assessmentAreaGradingStandards, err := repositories.FindAssessmentAreaGradingStandardsByGradingStandardNumber(schoolNumber, assessmentAreaNumber, gradingStandardNumber)
	if err != nil {
		return assessmentAreaResponses, errors.New("getting assessment area grading standard error: " + err.Error())
	}

	for _, assessmentAreaGradingStandard := range assessmentAreaGradingStandards {

		assessmentArea, gradingStandard, err := findAssessmentAreaAndGradingStandard(assessmentAreaGradingStandard.AssessmentAreaNumber, assessmentAreaGradingStandard.GradingStandardNumber)
		if err != nil {
			return assessmentAreaResponses, err
		}

		assessmentAreaResponses = append(assessmentAreaResponses, responses.GradingStandardAssessmentAreaResponse{
			Code:                assessmentAreaGradingStandard.Code,
			GradingStandardName: gradingStandard.Name,
			AssessmentAreaName:  assessmentArea.Name,
		})
	}

	return assessmentAreaResponses, nil
}

func SaveAssessmentAreaGradingStandard(assessmentAreaGradingStandardRequest aggregateScaleRequests.AssessmentAreaGradingStandardRequest) error {

	school, err := schoolRepositories.GetSchoolByNumber(assessmentAreaGradingStandardRequest.SchoolNumber)
	if err != nil {
		return errors.New("checking school number error: " + err.Error())
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	checkDuplicate, err := repositories.GetAssessmentAreaGradingStandardsByCodeAndSchoolNumber(assessmentAreaGradingStandardRequest.Code, school.Number)
	if err != nil {
		return errors.New("checking duplicate for assessment area code and school error: " + err.Error())
	}

	if checkDuplicate.Number != "" {
		return errors.New("code already used")
	}

	assessmentArea, gradingStandard, err := findAssessmentAreaAndGradingStandard(assessmentAreaGradingStandardRequest.AssessmentAreaNumber, assessmentAreaGradingStandardRequest.GradingStandardNumber)
	if err != nil {
		return err
	}

	assessmentAreaGradingStandard := models.AssessmentAreaGradingStandard{
		Number:                uuid.NewV4().String(),
		Code:                  assessmentAreaGradingStandardRequest.Code,
		AssessmentAreaNumber:  assessmentArea.Number,
		GradingStandardNumber: gradingStandard.Number,
		SchoolNumber:          school.Number,
	}

	err = repositories.SaveAssessmentAreaGradingStandard(assessmentAreaGradingStandard)
	if err != nil {
		return errors.New("failed saving grading standard for the assessment area")
	}

	return nil
}

func findAssessmentAreaAndGradingStandard(assessmentAreaNumber string, gradingStandardNumber string) (models.AssessmentArea, assessmentScaleModels.GradingStandard, error) {

	var assessmentArea models.AssessmentArea
	var gradingStandard assessmentScaleModels.GradingStandard

	assessmentArea, err := repositories.GetAssessmentAreaByNumber(assessmentAreaNumber)
	if err != nil {
		return assessmentArea, gradingStandard, errors.New("checking assessment area error: " + err.Error())
	}

	if assessmentArea.Number == "" {
		return assessmentArea, gradingStandard, errors.New("invalid assessment area number")
	}

	gradingStandard, err = assessmentRepositories.GetGradingStandardByNumber(gradingStandardNumber)
	if err != nil {
		return assessmentArea, gradingStandard, errors.New("checking grading standard error: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return assessmentArea, gradingStandard, errors.New("invalid grading standard number")
	}

	return assessmentArea, gradingStandard, nil
}

func FindAllAssessmentAreas(curriculumClassNumber, classNumber, curriculumNumber, levelNumber, curriculumLevelNumber string) ([]responses.AssessmentAreaResponse, error) {

	var assessmentAreasResponse = make([]responses.AssessmentAreaResponse, 0)

	all, err := repositories.FindAllAssessmentAreas(curriculumClassNumber, classNumber, curriculumNumber, levelNumber, curriculumLevelNumber)
	if err != nil {
		return assessmentAreasResponse, errors.New("Fetching assessment areas error: " + err.Error())
	}

	for _, assessmentArea := range all {

		name := assessmentArea.Name

		if assessmentArea.SubjectNumber != "" {
			subject, err := subjectRepositories.GetSubjectByNumber(assessmentArea.SubjectNumber)
			if err != nil {
				fmt.Println("Subject error: " + err.Error())
			} else {
				name = subject.Name
			}
		}

		curriculumClass, _ := repositories2.GetCurriculumClass(assessmentArea.CurriculumClassNumber)

		curriculumLevel, _ := repositories3.GetCurriculumLevelByNumber(curriculumClass.CurriculumLevelNumber)

		curriculum, _ := models2.GetCurriculumByNumber(curriculumLevel.CurriculumNumber)

		assessmentAreasResponse = append(assessmentAreasResponse, responses.AssessmentAreaResponse{
			Name:                  name,
			Number:                assessmentArea.Number,
			UserNumber:            assessmentArea.UserNumber,
			SubjectNumber:         assessmentArea.SubjectNumber,
			CurriculumClassNumber: assessmentArea.CurriculumClassNumber,
			CurriculumNumber:      curriculumLevel.CurriculumNumber,
			CurriculumName:        curriculum.Name,
		})
	}

	return assessmentAreasResponse, nil
}
