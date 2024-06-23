package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func SaveAssessmentArea(payload requests.AssessmentAreaRequest, currentUser string) error {

	trn := config.DB.Begin()

	for _, subjectNumber := range payload.SubjectNumber {
		err := processSaveAssessmentArea(subjectNumber, payload.CurriculumClassNumber, payload.Name, currentUser, trn)
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

func processSaveAssessmentArea(subjectNumber string, curriculumClassNumber string, otherName string, currentUser string, trn *gorm.DB) error {

	//validate for composite unique curriculum_class_number and subject number
	assessmentArea, err := repositories.FindAssessmentAreaBySubjectNumberAndCurriculumClassNumber(subjectNumber, curriculumClassNumber, trn)
	if err != nil {
		return err
	}

	if assessmentArea.Number != "" {
		fmt.Println("assessment area already exists for the specified curriculum class and subject")
		return nil
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
