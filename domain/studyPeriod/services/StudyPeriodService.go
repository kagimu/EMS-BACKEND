package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	studyModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type StudyPeriodService interface {
	NewStudyPeriodName(studyPeriodNameRequest requests.StudyPeriodNameRequest, currentUserNumber string) error
	UpdateStudyPeriodName(studyPeriodNameRequest requests.UpdateStudyPeriodNameRequest, currentUserNumber string) error
	FindAllStudyPeriodNames(schoolNumber string) ([]responses.StudyPeriodNameResponse, error)
	RemoveStudyPeriodName(studyPeriodNameNumber string) error
}

type studyPeriodService struct {
}

func NewStudyPeriodService() StudyPeriodService {
	return &studyPeriodService{}
}

func (s *studyPeriodService) NewStudyPeriodName(studyPeriodNameRequest requests.StudyPeriodNameRequest, currentUserNumber string) error {

	checkName, err := repositories.GetStudyPeriodNameByName(studyPeriodNameRequest.SchoolNumber, studyPeriodNameRequest.Name)
	if err != nil {
		return errors.New("Checking name duplicate error " + err.Error())
	}

	if checkName.Number != "" {
		return errors.New("Name already exists ")
	}

	err = repositories.SaveStudyPeriodName(studyModels.StudyPeriodName{
		Number:       uuid.NewV4().String(),
		Name:         studyPeriodNameRequest.Name,
		SchoolNumber: studyPeriodNameRequest.SchoolNumber,
		UserNumber:   currentUserNumber,
		Status:       constants.StatusActive,
	})

	if err != nil {
		return errors.New("Saving study period name error: " + err.Error())
	}

	return nil
}

func (s *studyPeriodService) FindAllStudyPeriodNames(schoolNumber string) ([]responses.StudyPeriodNameResponse, error) {

	var response = make([]responses.StudyPeriodNameResponse, 0)

	allNames, err := repositories.FindStudyPeriodName(schoolNumber)
	if err != nil {
		return response, errors.New("Fetching names error: " + err.Error())
	}

	for _, studyPeriodName := range allNames {
		response = append(response, createNameResponse(studyPeriodName))
	}

	return response, nil
}

func createNameResponse(studyPeriodName studyModels.StudyPeriodName) responses.StudyPeriodNameResponse {
	return responses.StudyPeriodNameResponse{
		DateAdded:             studyPeriodName.CreatedAt.Format(time.DateTime),
		StudyPeriodNameNumber: studyPeriodName.Number,
		Name:                  studyPeriodName.Name,
	}
}

func (s *studyPeriodService) UpdateStudyPeriodName(studyPeriodNameRequest requests.UpdateStudyPeriodNameRequest, currentUserNumber string) error {

	studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriodNameRequest.StudyPeriodNameNumber)
	if err != nil {
		return errors.New("Checking study period number error: " + err.Error())
	}

	checkName, err := repositories.GetStudyPeriodNameByName(studyPeriodNameRequest.SchoolNumber, studyPeriodNameRequest.Name)
	if err != nil {
		return errors.New("Checking name duplicate error " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != studyPeriodName.Number {
		return errors.New("Name already exists ")
	}

	studyPeriodName.UserNumber = currentUserNumber
	studyPeriodName.Name = studyPeriodNameRequest.Name

	err = repositories.SaveStudyPeriodName(studyPeriodName)

	if err != nil {
		return errors.New("Updating study period name error: " + err.Error())
	}

	return nil
}

func (s *studyPeriodService) RemoveStudyPeriodName(studyPeriodNameNumber string) error {

	studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriodNameNumber)
	if err != nil {
		return errors.New("Checking study period number error: " + err.Error())
	}

	err = repositories.RemoveStudyPeriodName(studyPeriodName)
	if err != nil {
		return errors.New("Removing study period name error: " + err.Error())
	}

	return nil
}
