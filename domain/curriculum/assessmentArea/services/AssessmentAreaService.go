package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func SaveAssessmentArea(payload requests.AssessmentAreaRequest, currentUser string) error {

	trn := config.DB.Begin()

	for _, subjectNumber := range payload.SubjectNumber {
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

	subject, err := repositories2.GetSubjectByNumber(subjectNumber, trn)
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
