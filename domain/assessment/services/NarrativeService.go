package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"github.com/araddon/dateparse"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type NarrativeService struct {
}

func (narrativeService NarrativeService) CreateNarrative(narrativeRequest requests.NarrativeRequest, userNumber string) error {
	return narrativeService.doCreateNarrative(narrativeRequest, "", userNumber)
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
	scheduledOn, err := dateparse.ParseAny(scheduledNarrativeRequest.ScheduledOn)
	if err != nil {
		return errors.New(fmt.Sprint("Passing scheduled on date error: ", err))
	}

	if scheduledOn.After(time.Now()) {
		return errors.New(fmt.Sprint("The scheduled on date should not be in the past date. "))
	}

	assessmentArea, err := repositories3.GetAssessmentAreaByNumber(scheduledNarrativeRequest.AssessmentAreaNumber)
	if err != nil {
		return err
	}
	if assessmentArea.Number == "" {
		return errors.New("Invalid assessment area number ")
	}

	campusClass, err := repositories4.GetCampusClassByNumber(scheduledNarrativeRequest.CampusClassNumber)
	if err != nil {
		return err
	}
	if campusClass.Number == "" {
		return errors.New("Invalid campus class number ")
	}

	studyPeriod, _ := repositories2.GetStudyPeriodByNumber(scheduledNarrativeRequest.StudyPeriodNumber)
	if studyPeriod.Number == "" {
		return errors.New("Invalid study period number ")
	}

	schedule := models.ScheduledNarrative{
		ScheduledOn:          scheduledOn,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CampusClassNumber:    scheduledNarrativeRequest.CampusClassNumber,
		AssessmentAreaNumber: scheduledNarrativeRequest.AssessmentAreaNumber,
		StudyPeriodNumber:    scheduledNarrativeRequest.StudyPeriodNumber,
		UserNumber:           user.Number,
		Description:          scheduledNarrativeRequest.Description,
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
			Description:          narrative.Description,
		})
	}

	return scheduledNarratives, nil
}

func (narrativeService NarrativeService) UpdateNarrativeEntry(narrativeRequest requests.NarrativeRequest, narrativeNumber, userNumber string) error {
	return narrativeService.doCreateNarrative(narrativeRequest, narrativeNumber, userNumber)
}

func (narrativeService NarrativeService) doCreateNarrative(narrativeRequest requests.NarrativeRequest, narrativeNumber string, userNumber string) error {
	studyPeriod, err := repositories2.GetStudyPeriodByNumber(narrativeRequest.StudyPeriodNumber)
	if err != nil {
		return err
	}
	if studyPeriod.Number == "" {
		return errors.New("Invalid study period number ")
	}

	assessmentArea, err := repositories3.GetAssessmentAreaByNumber(narrativeRequest.AssessmentAreaNumber)
	if err != nil {
		return err
	}
	if assessmentArea.Number == "" {
		return errors.New("Invalid assessment area number ")
	}

	campusClass, err := repositories4.GetCampusClassByNumber(narrativeRequest.CampusClassNumber)
	if err != nil {
		return err
	}
	if campusClass.Number == "" {
		return errors.New("Invalid campus class number ")
	}

	dateTimeString := narrativeRequest.ReferenceDate + " " + narrativeRequest.ReferenceTime + ":00"
	dateTime, err := dateparse.ParseAny(dateTimeString)
	if err != nil {
		return errors.New(fmt.Sprint("Passing reference date time error: ", err))
	}

	if dateTime.After(time.Now()) {
		return errors.New("Date and time should not be in the future ")
	}

	/*referenceDate, err := time.Parse(time.DateOnly, narrativeRequest.ReferenceDate)
	if err != nil {
		return errors.New(fmt.Sprint("Passing reference date error: ", err))
	}*/

	narrative := models.Narrative{
		UserNumber:           userNumber,
		StudentNumber:        narrativeRequest.StudentNumber,
		CampusClassNumber:    narrativeRequest.CampusClassNumber,
		AssessmentAreaNumber: narrativeRequest.AssessmentAreaNumber,
		Description:          narrativeRequest.Description,
		UpdatedAt:            time.Now(),
		StudyPeriodNumber:    narrativeRequest.StudyPeriodNumber,
		ReferenceDate:        dateTime,
		ReferenceTime:        narrativeRequest.ReferenceTime,
	}

	createdAt := time.Now()

	if narrativeNumber == "" {
		narrative.Number = uuid.NewV4().String()
		narrative.CreatedAt = createdAt
	} else {
		checkNarrative, err := repositories.FindNarrativeByNumber(narrativeNumber)
		if err != nil {
			return err
		}

		if checkNarrative.Number == "" {
			return errors.New("Invalid narrative number ")
		}

		narrative.CreatedAt = checkNarrative.CreatedAt
		narrative.Number = checkNarrative.Number
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

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		err := repositories.SaveNarrative(narrative, tx)
		if err != nil {
			return err
		}

		if scheduled.Number != "" {

			duplicateScheduleNarrative, err := repositories.FindNarrativeByScheduleNarrativeNumberAndStudentNumber(scheduled.Number, narrative.StudentNumber)
			if err != nil {
				return err
			}

			if duplicateScheduleNarrative.Number != "" && duplicateScheduleNarrative.Number != narrative.Number {
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

func CreateNarrativeResponse(narrative models.Narrative) responses.NarrativeResponse {

	user, _ := repository.GetUserByNumber(narrative.UserNumber)
	studyPeriod, _ := repositories2.GetStudyPeriodByNumber(narrative.StudyPeriodNumber)
	class, _ := repositories4.GetClassByCampusClassNumber(narrative.CampusClassNumber)
	academicYear, _ := repositories5.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
	assessmentArea, _ := repositories3.GetAssessmentAreaByNumber(narrative.AssessmentAreaNumber)
	spn, _ := repositories2.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	return responses.NarrativeResponse{
		Teacher:              user.FirstName + " " + user.LastName,
		Number:               narrative.Number,
		CreatedAt:            narrative.CreatedAt.Format(time.DateTime),
		AssessmentAreaNumber: narrative.AssessmentAreaNumber,
		StudyPeriod:          spn.Name,
		Class:                class.Name,
		Narration:            narrative.Description,
		AssessmentArea:       assessmentArea.Name,
		AcademicYear:         academicYear.AcademicYearName,
		ReferenceTime:        narrative.ReferenceTime,
		ReferenceDate:        narrative.ReferenceDate.Format(time.DateOnly),
	}
}

func NewNarrativeService() *NarrativeService {
	return &NarrativeService{}
}
