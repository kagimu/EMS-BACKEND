package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	var scheduled models.ScheduledNarrative

	if narrativeRequest.NarrativeScheduleNumber != "" {

		checkScheduled, err := repositories.FindScheduledNarrativeByNumber(narrativeRequest.NarrativeScheduleNumber)
		if err != nil {
			return err
		}

		if checkScheduled.Number == "" {
			return errors.New("Invalid schedule selected ")
		}

		if checkScheduled.ScheduledOn.Before(time.Now()) {
			return errors.New("This schedule is not yet due ")
		}

		scheduled = checkScheduled
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		err := repositories.SaveNarrative(narrative, tx)
		if err != nil {
			return err
		}

		if scheduled.Number != "" {

			checkDueScheduled, err := repositories.FindNarrativeByScheduleNumberAndStudentNumber(scheduled.Number, narrative.StudentNumber)
			if err != nil {
				return err
			}

			if checkDueScheduled.Number != "" {
				return errors.New("Assessment for this learner already entered ")
			}

			err = repositories.SaveNarrativeAndSchedule(models.NarrativeSchedule{ScheduledNarrativeNumber: scheduled.Number, NarrativeNumber: narrative.Number}, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

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

func (narrativeService NarrativeService) SaveSchedule(scheduledNarrativeRequest requests.ScheduledNarrativeRequest, user *models2.User) error {
	scheduledOn, err := time.Parse(time.DateOnly, scheduledNarrativeRequest.ScheduledOn)
	if err != nil {
		return errors.New(fmt.Sprint("Passing scheduled on date error: ", err))
	}

	schedule := models.ScheduledNarrative{
		ScheduledOn:          scheduledOn,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CampusClassNumber:    scheduledNarrativeRequest.CampusClassNumber,
		AssessmentAreaNumber: scheduledNarrativeRequest.AssessmentAreaNumber,
		StudyPeriodNumber:    scheduledNarrativeRequest.StudyPeriodNumber,
		UserNumber:           user.Number,
	}

	err = repositories.SaveScheduledNarrative(schedule)
	if err != nil {
		return errors.New(fmt.Sprint("Error saving scheduled narrative assessment: ", err))
	}

	return nil
}

func (narrativeService NarrativeService) AllSchedules(studyPeriodNumber, campusClassNumber, assessmentAreaNumber string) ([]responses.ScheduledNarrativeResponse, error) {

	data, err := repositories.FindAllScheduledNarratives(studyPeriodNumber, campusClassNumber, assessmentAreaNumber, time.Now())
	if err != nil {
		return nil, err
	}

	var scheduledNarratives = make([]responses.ScheduledNarrativeResponse, 0)

	for _, narrative := range data {
		scheduledNarratives = append(scheduledNarratives, responses.ScheduledNarrativeResponse{
			Number:               narrative.Number,
			ScheduledOn:          narrative.ScheduledOn.Format(time.DateOnly),
			CampusClassNumber:    narrative.CampusClassNumber,
			AssessmentAreaNumber: narrative.AssessmentAreaNumber,
			StudyPeriodNumber:    narrative.StudyPeriodNumber,
		})
	}

	return scheduledNarratives, nil
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
