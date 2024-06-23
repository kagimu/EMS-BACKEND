package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	uuid "github.com/satori/go.uuid"
)

func SaveAssessmentArea(payload requests.AssessmentAreaRequest, currentUser string) error {

	//validate for composite unique curriculum_class_number and subject number
	assessmentArea, err := repositories.FindAssessmentAreaBySubjectNumberAndCurriculumClassNumber(payload.SubjectNumber, payload.CurriculumClassNumber)
	if err != nil {
		return err
	}

	if assessmentArea.Number != "" {
		return errors.New("assessment area already exists for the specified curriculum class and subject")
	}

	saveAssessmentArea := models.AssessmentArea{
		Number:                uuid.NewV4().String(),
		CurriculumClassNumber: payload.CurriculumClassNumber,
		Name:                  payload.Name,
		UserNumber:            currentUser,
		SubjectNumber:         payload.SubjectNumber,
	}

	_, err = repositories.SaveAssessmentArea(saveAssessmentArea)
	if err != nil {
		return errors.New("failed to create assessment area")
	}

	return nil
}
