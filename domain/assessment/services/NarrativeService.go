package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	uuid "github.com/satori/go.uuid"
	"time"
)

type NarrativeService struct {
}

func (narrativeService NarrativeService) CreateNarrative(narrativeRequest requests.NarrativeRequest, userNumber string) error {

	narrative := models.Narrative{
		UserNumber:           userNumber,
		StudentNumber:        narrativeRequest.StudentNumber,
		CampusClassNumber:    narrativeRequest.CampusClassNumber,
		AssessmentAreaNumber: narrativeRequest.AssessmentAreaNumber,
		Description:          narrativeRequest.Description,
		Number:               uuid.NewV4().String(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err := repositories.SaveNarrative(narrative)
	if err != nil {
		return errors.New(fmt.Sprint("Error saving narrative assessment: ", err))
	}

	return nil
}

func (narrativeService NarrativeService) AllNarrative(request requests.NarrativeRequest, startDate, endDate time.Time) ([]responses.NarrativeResponse, error) {

	narrativeResponse := make([]responses.NarrativeResponse, 0)

	narratives, err := repositories.FindAllNarratives(request.StudentNumber, request.StudyPeriodNumber, request.CampusClassNumber, request.AssessmentAreaNumber, startDate, endDate)
	if err != nil {
		return narrativeResponse, err
	}

	for _, narrative := range narratives {
		narrativeResponse = append(narrativeResponse, CreateNarrativeResponse(narrative))
	}

	return narrativeResponse, nil
}

func CreateNarrativeResponse(narrative models.Narrative) responses.NarrativeResponse {

	user, _ := repository.GetUserByNumber(narrative.UserNumber)

	return responses.NarrativeResponse{
		Teacher:              user.FirstName + " " + user.LastName,
		Number:               narrative.Number,
		Date:                 narrative.CreatedAt.Format(time.DateTime),
		AssessmentAreaNumber: narrative.AssessmentAreaNumber,
	}
}

func NewNarrativeService() *NarrativeService {
	return &NarrativeService{}
}
