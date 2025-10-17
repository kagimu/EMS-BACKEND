package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	repositories11 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repositories10 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/repositories"
	repositories9 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/repositories"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/responses"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	repositories8 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/repositories"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

		uploads, _ := repositories.FindNarrativeScriptMarkingGuideByNumber(narrative.Number)

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
			ScriptNumber:         uploads.ScriptNumber,
			MarkingGuideNumber:   uploads.MarkingGuideNumber,
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

	if dateTime.Before(studyPeriod.StudyPeriodStartDate) {
		return errors.New("Reference date should be after the study period starts ")
	}

	if dateTime.After(studyPeriod.StudyPeriodEndDate) {
		return errors.New("Reference date should be withing the study period ")
	}

	narrative := models.Narrative{
		UserNumber:              userNumber,
		StudentNumber:           narrativeRequest.StudentNumber,
		CampusClassNumber:       narrativeRequest.CampusClassNumber,
		AssessmentAreaNumber:    narrativeRequest.AssessmentAreaNumber,
		Description:             narrativeRequest.Description,
		UpdatedAt:               time.Now(),
		StudyPeriodNumber:       narrativeRequest.StudyPeriodNumber,
		ReferenceDate:           dateTime,
		ReferenceTime:           narrativeRequest.ReferenceTime,
		CurrentStatus:           narrativeRequest.CurrentStatus,
		Status:                  narrativeRequest.Status,
		SuggestedSolution:       narrativeRequest.SuggestedSolution,
		NarrativeScheduleNumber: narrativeRequest.NarrativeScheduleNumber,
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

		if time.Now().Before(checkScheduled.ScheduledOn) {
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

			if scheduledOn.Before(studyPeriod.StudyPeriodStartDate) {
				return errors.New(fmt.Sprint("The scheduled on date should not be before the study period starts. "))
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

func (narrativeService NarrativeService) NarrativeReportPrintPreview(narrative requests.NarrativeRequest, startDate, endDate time.Time, user *models2.User) (responses.NarrativeResponsePrintPreview, error) {

	student, err := repositories6.FindStudentByNumber(narrative.StudentNumber)
	if err != nil {
		return responses.NarrativeResponsePrintPreview{}, err
	}

	if student.Number == "" {
		return responses.NarrativeResponsePrintPreview{}, errors.New("Invalid study number ")
	}

	studyPeriod, _ := repositories2.GetStudyPeriodByNumber(narrative.StudyPeriodNumber)
	spn, _ := repositories2.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)
	class, _ := repositories4.GetClassByCampusClassNumber(narrative.CampusClassNumber)

	school, _ := repositories7.GetSchoolByNumber(spn.SchoolNumber)
	campus, _ := repositories11.GetCampusByCampusClassNumber(narrative.CampusClassNumber)

	response := responses.NarrativeResponsePrintPreview{
		Class:                class.Name,
		StudyPeriod:          spn.Name,
		DateOfConference:     time.Now().Format("2006-01-02"),
		CurrentAdministrator: user.FirstName + " " + user.LastName,
		RegistrationNumber:   student.ADN,
		ChildFirstName:       student.Surname,
		ChildSecondName:      student.SecondName,
		School: responses2.SchoolResponse{
			Name:             school.Name,
			Number:           school.Number,
			Logo:             school.Logo,
			Motto:            school.Motto,
			PrincipalName:    school.PrincipalName,
			CommencementDate: school.CommencementDate.Format("2006-01-02"),
		},
		Campus: responses3.CampusResponse{
			Name:             campus.Name,
			Number:           campus.Number,
			Phone:            campus.Phone,
			Email:            campus.Email,
			Location:         campus.Location,
			CommencementDate: campus.CommencementDate,
		},
	}

	narrationData, err := repositories.GetNarrativesGroupedByAssessmentArea(narrative.StudentNumber, narrative.StudyPeriodNumber, "", startDate, endDate)
	if err != nil {
		return response, err
	}

	var narrativeResponseData = make([]responses.NarrativeResponseData, 0)

	for assessmentAreaNumber, narrationList := range narrationData {

		assessmentArea, err := repositories3.GetAssessmentAreaByNumber(assessmentAreaNumber)
		if err != nil {
			return response, err
		}

		subject, err := repositories8.GetSubjectByNumber(assessmentArea.SubjectNumber)
		if err != nil {
			return response, err
		}

		subjectCluster, err := repositories9.GetSubjectClusterByNumber(subject.SubjectClusterNumber)
		if err != nil {
			return response, err
		}

		learningArea, err := repositories10.GetLearningAreaByNumber(subjectCluster.LearningAreaNumber)
		if err != nil {
			return response, err
		}

		var narrativeResponseRows = make([]responses.NarrativeResponseRow, 0)

		for _, narration := range narrationList {
			narrativeResponseRows = append(narrativeResponseRows, responses.NarrativeResponseRow{
				Status:                  narration.Status,
				CurrentStatus:           narration.CurrentStatus,
				Narration:               narration.Description,
				SuggestedSolution:       narration.SuggestedSolution,
				ReferenceDate:           narration.ReferenceDate.Format("2006-01-02"),
				ReferenceTime:           narration.ReferenceTime,
				CreatedAt:               narration.CreatedAt.Format("2006-01-02"),
				NarrativeScheduleNumber: narration.NarrativeScheduleNumber,
			})
		}

		narrativeResponseData = append(narrativeResponseData, responses.NarrativeResponseData{
			AssessmentAreaNumber:  assessmentAreaNumber,
			AssessmentArea:        subject.Name,
			SubjectCode:           strconv.FormatInt(subject.Code, 10),
			SubjectCluster:        subjectCluster.Name,
			LearningArea:          learningArea.Name,
			NarrativeResponseRows: narrativeResponseRows,
		})

	}

	response.NarrativeResponseData = narrativeResponseData

	return response, nil
}

func CreateNarrativeResponse(narrative models.Narrative) responses.NarrativeResponse {

	user, _ := repository.GetUserByNumber(narrative.UserNumber)
	studyPeriod, _ := repositories2.GetStudyPeriodByNumber(narrative.StudyPeriodNumber)
	class, _ := repositories4.GetClassByCampusClassNumber(narrative.CampusClassNumber)
	academicYear, _ := repositories5.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
	assessmentArea, _ := repositories3.GetAssessmentAreaByNumber(narrative.AssessmentAreaNumber)
	spn, _ := repositories2.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	row := responses.NarrativeResponseRow{
		Teacher:                 user.FirstName + " " + user.LastName,
		CreatedAt:               narrative.CreatedAt.Format(time.DateTime),
		Narration:               narrative.Description,
		ReferenceTime:           narrative.ReferenceTime,
		ReferenceDate:           narrative.ReferenceDate.Format(time.DateOnly),
		Status:                  narrative.Status,
		CurrentStatus:           narrative.CurrentStatus,
		SuggestedSolution:       narrative.SuggestedSolution,
		NarrativeScheduleNumber: narrative.NarrativeScheduleNumber,
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

func UploadNarrativeMarkingGuide(ctx *gin.Context, scheduledNarrativeNumber, currentUserNumber string) error {

	narrative, err := repositories.FindScheduledNarrativeByNumber(scheduledNarrativeNumber)
	if err != nil {
		return errors.New("Error getting narrative details: " + err.Error())
	}

	if narrative.Number == "" {
		return errors.New("Invalid narrative number ")
	}

	scriptFile, err := services2.UploadAttachment(ctx, "script-file")
	if err != nil {
		return errors.New("uploading narrative script file error: " + err.Error())
	}

	markingGuideFile, err := services2.UploadAttachment(ctx, "marking-guide-file")
	if err != nil {
		return errors.New("uploading narrative marking guide error: " + err.Error())
	}

	//scriptMarkingGuideRepository.UseTx(config.DB)
	enty := models.ScheduledNarrativeScriptMarkingGuide{
		ScriptNumber:             scriptFile,
		MarkingGuideNumber:       markingGuideFile,
		UserNumber:               currentUserNumber,
		ScheduledNarrativeNumber: scheduledNarrativeNumber,
	}

	oldScripts, err := repositories.FindNarrativeScriptMarkingGuideByNumber(scheduledNarrativeNumber)
	if err != nil {
		return errors.New("Error checking old narrative script and marking guide error: " + err.Error())
	}

	if oldScripts.ScheduledNarrativeNumber != "" {
		enty.CreatedAt = oldScripts.CreatedAt
	}

	enty.UpdatedAt = oldScripts.UpdatedAt

	err = repositories.SaveScheduledNarrativeScriptMarkingGuide(enty)

	if err != nil {
		return errors.New("storing narrative script marking guide error: " + err.Error())
	}

	return nil
}

func NewNarrativeService() *NarrativeService {
	return &NarrativeService{}
}
