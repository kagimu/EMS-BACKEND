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
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
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
	return narrativeService.doSaveSchedule(scheduledNarrativeRequest, user, "")
}

func (narrativeService NarrativeService) AllSchedules(studyPeriodNumber, classNumber, campusNumber, curriculumNumber, assessmentAreaNumber string) ([]responses.ScheduledNarrativeResponse, error) {

	campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber, "")
	if err != nil {
		return nil, err
	}

	if campusClass.Number == "" {
		return nil, errors.New("Could not find campus class ")
	}

	studyPeriod, err := repositories2.GetStudyPeriodByNumber(studyPeriodNumber)
	if err != nil {
		return nil, err
	}

	data, err := repositories.FindAllScheduledNarratives(studyPeriodNumber, campusClass.Number, assessmentAreaNumber, studyPeriod.StudyPeriodStartDate, studyPeriod.StudyPeriodEndDate)
	if err != nil {
		return nil, err
	}

	var scheduledNarratives = make([]responses.ScheduledNarrativeResponse, 0)

	for _, narrative := range data {

		name := "Study Period"

		studyPeriodName, err := repositories2.GetStudyPeriodNameByStudyPeriodNumber(narrative.StudyPeriodNumber)
		if err != nil {
			name = err.Error()
		} else {
			name = studyPeriodName.Name
		}

		assessmentArea, _ := repositories3.GetAssessmentAreaByNumber(narrative.AssessmentAreaNumber)

		class, _ := repositories4.GetClassByCampusClassNumber(narrative.CampusClassNumber)

		scheduledNarratives = append(scheduledNarratives, responses.ScheduledNarrativeResponse{
			Number:               narrative.Number,
			ScheduledOn:          narrative.ScheduledOn.Format(time.DateOnly),
			CampusClassNumber:    narrative.CampusClassNumber,
			AssessmentAreaNumber: narrative.AssessmentAreaNumber,
			StudyPeriodNumber:    narrative.StudyPeriodNumber,
			Description:          narrative.Description,
			StudyPeriod:          name,
			AssessmentArea:       assessmentArea.Name,
			Class:                class.Name,
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

	if dateTime.After(time.Now().UTC()) {
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
		CurrentStatus:        narrativeRequest.CurrentStatus,
		Status:               narrativeRequest.Status,
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

func (narrativeService NarrativeService) NarrativeDetails(narrativeNumber string) (responses.NarrativeResponse, error) {
	narrative, err := repositories.FindNarrativeByNumber(narrativeNumber)
	if err != nil {
		return responses.NarrativeResponse{}, errors.New(fmt.Sprint(fmt.Sprintf("Error getting narrative details: %s", err)))
	}

	if narrative.Number == "" {
		return responses.NarrativeResponse{}, errors.New("Invalid narrative number ")
	}

	return CreateNarrativeResponse(narrative), nil
}

func (narrativeService NarrativeService) DeleteNarrative(narrativeNumber string) error {
	narrative, err := repositories.FindNarrativeByNumber(narrativeNumber)
	if err != nil {
		return err
	}

	err = repositories.DeleteNarrative(narrative)
	if err != nil {
		return errors.New(fmt.Sprint(fmt.Sprintf("Error deleting narrative details: %s", err.Error())))
	}

	return nil
}

func (narrativeService NarrativeService) DeleteNarrativeSchedule(scheduleNumber string) error {
	data, err := repositories.GetScheduledNarrative(scheduleNumber)
	if err != nil {
		return err
	}

	if data.Number == "" {
		return errors.New("Invalid narrative schedule number ")
	}

	allSubmissions, err := repositories.FindAllScheduledNarrativeSubmissions(data.Number)
	if err != nil {
		return errors.New(fmt.Sprint(fmt.Sprintf("Error finding all submissions: %s", err)))
	}

	size := len(allSubmissions)

	if size > 0 {
		return errors.New(fmt.Sprintf("This schedule has '%d' submissions, can not be deleted.", size))
	}

	err = repositories.DeleteScheduledNarrative(data)
	if err != nil {
		return errors.New(fmt.Sprint(fmt.Sprintf("Error deleting narrative schedule: %s", err.Error())))
	}

	return nil
}

func (narrativeService NarrativeService) UpdateNarrativeSchedule(scheduledNarrativeRequest requests.ScheduledNarrativeRequest, scheduleNumber string, user *models2.User) error {
	return narrativeService.doSaveSchedule(scheduledNarrativeRequest, user, scheduleNumber)
}

func (narrativeService NarrativeService) doSaveSchedule(scheduledNarrativeRequest requests.ScheduledNarrativeRequest, user *models2.User, scheduleNumber string) error {

	assessmentArea, err := repositories3.GetAssessmentAreaByNumber(scheduledNarrativeRequest.AssessmentAreaNumber)
	if err != nil {
		return err
	}
	if assessmentArea.Number == "" {
		return errors.New("Invalid assessment area number ")
	}

	campusClass, err := services.ResolveCampusClass(scheduledNarrativeRequest.ClassNumber, scheduledNarrativeRequest.CampusNumber, scheduledNarrativeRequest.CurriculumNumber, "")
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

	createdAt := time.Now()

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, scheduledOn := range scheduledNarrativeRequest.ScheduledOn {

			scheduledOn, err := dateparse.ParseAny(scheduledOn)
			if err != nil {
				return errors.New(fmt.Sprint("Passing scheduled on date error: ", err))
			}

			if scheduledOn.Before(time.Now()) {
				return errors.New(fmt.Sprint("The scheduled on date should not be in the past date. "))
			}

			if scheduledOn.After(studyPeriod.StudyPeriodEndDate) {
				return errors.New(fmt.Sprint(fmt.Sprintf("Schedule date exceeds the study period end date of '%s'", studyPeriod.StudyPeriodEndDate.Format(time.DateOnly))))
			}

			schedule := models.ScheduledNarrative{
				ScheduledOn:          scheduledOn,
				CreatedAt:            createdAt,
				UpdatedAt:            createdAt,
				CampusClassNumber:    campusClass.Number,
				AssessmentAreaNumber: scheduledNarrativeRequest.AssessmentAreaNumber,
				StudyPeriodNumber:    scheduledNarrativeRequest.StudyPeriodNumber,
				UserNumber:           user.Number,
				Description:          scheduledNarrativeRequest.Description,
				Number:               uuid.NewV4().String(),
			}

			if scheduleNumber != "" {
				oldSchedule, err := repositories.FindScheduledNarrativeByNumber(scheduleNumber, tx)
				if err != nil {
					return err
				}

				if oldSchedule.Number == "" {
					return errors.New("Invalid narrative schedule number ")
				}

				schedule.Number = oldSchedule.Number
				schedule.CreatedAt = oldSchedule.CreatedAt
			}

			err = repositories.SaveScheduledNarrative(schedule, tx)
			if err != nil {
				return errors.New(fmt.Sprint("Error saving scheduled narrative assessment: ", err))
			}

		}

		return err
	})

	return err
}

func CreateNarrativeResponse(narrative models.Narrative) responses.NarrativeResponse {

	user, _ := repository.GetUserByNumber(narrative.UserNumber)
	studyPeriod, _ := repositories2.GetStudyPeriodByNumber(narrative.StudyPeriodNumber)
	class, _ := repositories4.GetClassByCampusClassNumber(narrative.CampusClassNumber)
	academicYear, _ := repositories5.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
	assessmentArea, _ := repositories3.GetAssessmentAreaByNumber(narrative.AssessmentAreaNumber)
	spn, _ := repositories2.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	row := responses.NarrativeResponseRow{
		Teacher:       user.FirstName + " " + user.LastName,
		CreatedAt:     narrative.CreatedAt.Format(time.DateTime),
		Narration:     narrative.Description,
		ReferenceTime: narrative.ReferenceTime,
		ReferenceDate: narrative.ReferenceDate.Format(time.DateOnly),
		Status:        narrative.Status,
		CurrentStatus: narrative.CurrentStatus,
	}

	return responses.NarrativeResponse{
		Number:               narrative.Number,
		AssessmentAreaNumber: narrative.AssessmentAreaNumber,
		StudyPeriod:          spn.Name,
		Class:                class.Name,
		AssessmentArea:       assessmentArea.Name,
		AcademicYear:         academicYear.AcademicYearName,
		NarrativeResponseRow: row,
	}
}

func NewNarrativeService() *NarrativeService {
	return &NarrativeService{}
}
