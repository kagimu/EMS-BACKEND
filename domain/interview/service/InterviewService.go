package service

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	assessmentAreaRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	enrollmentResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	interviewResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	interviewRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repositories"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repository"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var assessmentMetricRepository = assessmentRepositories.NewAssessmentMetricRepository()

func FetchInterviewAssessmentMetricByStudent(studentNumber, interviewNumber string) ([]interviewResponses.StudentInterviewResults, error) {

	var studentInterviewResults []interviewResponses.StudentInterviewResults

	student, err := studentRepositories.FindStudentByUserNumber(studentNumber)
	if err != nil {
		return studentInterviewResults, errors.New("checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return studentInterviewResults, errors.New("invalid student number")
	}

	assessments, err := assessmentRepositories.FindAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, "")
	if err != nil {
		return studentInterviewResults, errors.New("checking student assessments error: " + err.Error())
	}

	fmt.Println("Assessments attempted: " + strconv.Itoa(len(assessments)))

	interviewAssessmentMetrics, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumberAndCampusClassNumber(interviewNumber, student.CampusClassNumber)
	if err != nil {
		return studentInterviewResults, errors.New("getting interview assessment metric error: " + err.Error())
	}

	for _, interviewAssessmentMetric := range interviewAssessmentMetrics {

		assessmentMetric, err := assessmentMetricRepository.FindAssessmentMetricByNumber(interviewAssessmentMetric.AssessmentMetricNumber)
		if err != nil {
			return studentInterviewResults, errors.New("getting assessment metric error: " + err.Error())
		}

		assessmentAreaName, err := assessmentAreaRepositories.GetAssessmentAreaByNumber(assessmentMetric.AssessmentAreaNumber)
		if err != nil {
			return studentInterviewResults, errors.New("error getting assessment area: " + err.Error())
		}

		if assessmentAreaName.Number == "" {
			return studentInterviewResults, errors.New("no assessment area found")
		}

		class, err := classRepositories.GetClassByCampusClassNumber(assessmentMetric.CampusClassNumber)
		if err != nil {
			return studentInterviewResults, errors.New("error getting class: " + err.Error())
		}

		if class.Number == "" {
			return studentInterviewResults, errors.New("campus class number " + assessmentMetric.CampusClassNumber + " has no class")
		}

		assessment, err := assessmentRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber, assessmentMetric.Number)
		if err != nil {
			return studentInterviewResults, errors.New("error getting learner scores: " + err.Error())
		}

		gradingStandards, err := assessmentRepositories.FindGradingStandardByCampusNumberAndCampusClassNumber("", class.Number, "")
		if err != nil {
			return studentInterviewResults, errors.New("error getting " + class.Name + "'s grading standard: " + err.Error())
		}

		if len(gradingStandards) == 0 {
			return studentInterviewResults, errors.New("no grading standard set for: " + class.Name)
		}

		if len(gradingStandards) > 1 {
			return studentInterviewResults, errors.New("Class " + class.Name + " has more than one standard set, choose the default standard to use.")
		}

		aggregateScales, err := assessmentRepositories.FindAggregateScaleByClassNumberAndGradingStandardNumber(class.Number, gradingStandards[0].Number)
		if err != nil {
			return studentInterviewResults, errors.New("error getting aggregate scale: " + err.Error())
		}

		var comment string
		if assessment.Comment == "" {
			comment = "No comment"
		} else {
			comment = assessment.Comment
		}

		if len(aggregateScales) == 0 {
			studentInterviewResults = append(studentInterviewResults, createInterviewResult(assessmentAreaName.Name, comment, "Scale not set", assessment.ActualScore))

		} else {
			for _, aggregateScale := range aggregateScales {

				if assessment.ActualScore >= aggregateScale.MinValue && assessment.ActualScore <= aggregateScale.MaxValue {

					grade, err := repository.FindGradeByNumber(aggregateScale.GradeNumber)

					if err != nil {
						return studentInterviewResults, errors.New("error resolving assessment grade: " + err.Error())
					}

					if grade.Number == "" {
						return studentInterviewResults, errors.New("invalid assessment grade number: " + aggregateScale.GradeNumber)
					}

					studentInterviewResults = append(studentInterviewResults, createInterviewResult(assessmentAreaName.Name, comment, grade.Name, assessment.ActualScore))
				}
			}
		}
	}

	return studentInterviewResults, nil
}

func createInterviewResult(name, comment, grade string, score float64) interviewResponses.StudentInterviewResults {
	return interviewResponses.StudentInterviewResults{
		Name:    name,
		Comment: comment,
		Score:   score,
		Grade:   grade,
	}
}

func CreateInterview(interviewRequest requests.InterviewRequest, currentUser string) error {

	scheduledOn, err := commons.ConvertStringToDate(interviewRequest.InterviewDate)
	if err != nil {
		return err
	}

	campusClass, err := services2.ResolveCampusClass(interviewRequest.ClassNumber, interviewRequest.CampusNumber, interviewRequest.CurriculumNumber, interviewRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	tx := config.DB.Begin()

	var repo = repository2.NewInterviewRepoImpl(tx)

	interviewNumber := uuid.NewV4().String()
	saveInterview := models.Interview{
		Number:            interviewNumber,
		Message:           interviewRequest.Message,
		Venue:             interviewRequest.Venue,
		Name:              interviewRequest.Name,
		InterviewDate:     scheduledOn,
		UserNumber:        currentUser,
		CampusClassNumber: campusClass.Number,
	}

	if err := repo.SaveInterview(saveInterview); err != nil {
		tx.Rollback()
		return errors.New("creating interview error: " + err.Error())
	}

	err = processInterviewAssessments(saveInterview, interviewRequest.AssessmentMetrics, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("Failed to commit transaction " + err.Error())
	}

	return nil
}

func processInterviewAssessments(interview models.Interview, assessmentMetrics []string, tx *gorm.DB) error {

	for _, assessmentMetricNumber := range assessmentMetrics {

		assessmentMetric, err := assessmentMetricRepository.FindAssessmentMetricByNumber(assessmentMetricNumber, tx)
		if err != nil {
			return errors.New("checking assessment schedule error: " + err.Error())
		}

		if assessmentMetric.Number == "" {
			return errors.New("invalid selected assessment")
		}

		checkInterviewAssessmentMetric, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumberAndAssessmentMetricNumber(interview.Number, assessmentMetricNumber, tx)
		if err != nil {
			return errors.New("checking duplicates error: " + err.Error())
		}

		if checkInterviewAssessmentMetric.Number != "" {
			continue
		}

		interviewAssessmentMetric := models.InterviewAssessmentMetric{
			Number:                 uuid.NewV4().String(),
			InterviewNumber:        interview.Number,
			AssessmentMetricNumber: assessmentMetricNumber,
		}

		err = interviewRepositories.SaveInterviewAssessmentMetric(interviewAssessmentMetric, tx)
		if err != nil {
			return errors.New("Failed to create interview assessment metric: " + err.Error())
		}

	}

	return nil
}

func UpdateInterview(editInterviewRequest requests.EditInterviewRequest, currentUserNumber string) error {

	scheduledOn, err := commons.ConvertStringToDate(editInterviewRequest.InterviewDate)
	if err != nil {
		return err
	}

	oldSchedules, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumber(editInterviewRequest.InterviewNumber)
	if err != nil {
		return errors.New("getting old schedules error: " + err.Error())
	}

	campusClass, err := services2.ResolveCampusClass(editInterviewRequest.ClassNumber, editInterviewRequest.CampusNumber, editInterviewRequest.CurriculumNumber, editInterviewRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, oldSchedule := range oldSchedules {
			err = interviewRepositories.RemoveInterviewAssessmentMetric(oldSchedule, tx)
			if err != nil {
				return errors.New("removing previous schedule error: " + err.Error())
			}
		}

		var repo = repository2.NewInterviewRepoImpl(tx)

		interview, err := repo.FindInterviewByNumber(editInterviewRequest.InterviewNumber)
		if err != nil {
			return errors.New("checking interview by number error: " + err.Error())
		}

		if interview.Number == "" {
			return errors.New("invalid interview number")
		}

		checkName, err := repo.FindInterviewByName(interview.Name, campusClass.Number)
		if err != nil {
			return errors.New("checking duplicates by name error: " + err.Error())
		}

		if checkName.Number != "" && interview.Number != checkName.Number {
			return errors.New("name already used on another interview")
		}

		interview.Name = editInterviewRequest.Name
		interview.CampusClassNumber = campusClass.Number
		interview.InterviewDate = scheduledOn
		interview.UserNumber = currentUserNumber
		interview.Message = editInterviewRequest.Message
		interview.Venue = editInterviewRequest.Venue

		err = repo.SaveInterview(interview)
		if err != nil {
			return errors.New("updating interview error: " + err.Error())
		}

		err = processInterviewAssessments(interview, editInterviewRequest.AssessmentMetrics, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetInterviewDetails(interviewNumber string) (interviewResponses.InterviewResponseDetails, error) {

	var interviewResponseDetails interviewResponses.InterviewResponseDetails

	var repo = repository2.NewInterviewRepoImpl(config.DB)

	interview, err := repo.FindInterviewByNumber(interviewNumber)
	if err != nil {
		return interviewResponseDetails, errors.New("getting interview by number error: " + err.Error())
	}

	if interview.Number == "" {
		return interviewResponseDetails, errors.New("invalid interview number")
	}

	interViewSchedules, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumber(interview.Number)
	if err != nil {
		return interviewResponseDetails, errors.New("getting interview schedules error: " + err.Error())
	}

	interviewResponse, err := createInterviewResponse(interview)
	if err != nil {
		return interviewResponseDetails, err
	}

	interviewResponseDetails.InterviewResponse = interviewResponse

	var schedules = make([]responses.AssessmentMetricResponse, 0)

	assessmentMetricRepository := assessmentRepositories.NewAssessmentMetricRepository()

	for _, interviewSchedule := range interViewSchedules {
		assessment, err := assessmentMetricRepository.FindAssessmentMetricByNumber(interviewSchedule.AssessmentMetricNumber)
		if err != nil {
			return interviewResponseDetails, errors.New("resolving assessment schedule error: " + err.Error())
		}

		response, err := services.CreateAssessmentsResponse(assessment)
		if err != nil {
			return interviewResponseDetails, err
		}

		schedules = append(schedules, response)
	}

	interviewResponseDetails.AssessmentSchedules = schedules

	return interviewResponseDetails, nil
}

func createInterviewResponse(interview models.Interview) (interviewResponses.InterviewResponse, error) {

	campusClass, err := classRepositories.GetCampusClassByNumber(interview.CampusClassNumber)
	if err != nil {
		return interviewResponses.InterviewResponse{}, errors.New("checking campus class error: " + err.Error())
	}

	if campusClass.Number == "" {
		return interviewResponses.InterviewResponse{}, errors.New("invalid campus class number")
	}

	campus, err := repositories.GetCampusByNumber(campusClass.CampusNumber)
	if err != nil {
		return interviewResponses.InterviewResponse{}, errors.New("checking campus error: " + err.Error())
	}

	class, err := classRepositories.GetClassByCampusClassNumber(campusClass.Number)
	if err != nil {
		return interviewResponses.InterviewResponse{}, errors.New("getting class error: " + err.Error())
	}

	return interviewResponses.InterviewResponse{
		Class:             class.Name,
		CampusClassNumber: interview.CampusClassNumber,
		InterviewDate:     interview.InterviewDate.Format(time.DateOnly),
		Message:           interview.Message,
		Venue:             interview.Venue,
		Name:              interview.Name,
		InterviewNumber:   interview.Number,
		SchoolNumber:      campus.SchoolNumber,
		CampusNumber:      campusClass.CampusNumber,
	}, nil
}

func AllInterviews(campusClassNumber string) ([]interviewResponses.InterviewResponse, error) {

	var allInterviewResponses = make([]interviewResponses.InterviewResponse, 0)

	repo := repository2.NewInterviewRepoImpl(config.DB)

	interviews, err := repo.FindAllInterviews(campusClassNumber)
	if err != nil {
		return allInterviewResponses, errors.New("fetching all interviews error: " + err.Error())
	}

	for _, interview := range interviews {

		interviewResponse, err := createInterviewResponse(interview)
		if err != nil {
			return allInterviewResponses, err
		}

		allInterviewResponses = append(allInterviewResponses, interviewResponse)
	}

	return allInterviewResponses, nil
}

func InterviewStudents(assessmentMetricNumber, interviewNumber string) ([]responses2.SingleStudent, error) {

	var studentsList = make([]responses2.SingleStudent, 0)

	var repo = repository2.NewInterviewRepoImpl(config.DB)

	interview, err := repo.FindInterviewByNumber(interviewNumber)
	if err != nil {
		return studentsList, errors.New("error checking interview by number: " + err.Error())
	}

	if interview.Number == "" {
		return studentsList, errors.New("invalid interview number")
	}

	interviewStudents, err := repository2.FindInterviewStudentsByInterViewNumber(interview.Number)
	if err != nil {
		return studentsList, errors.New("finding students by interview error: " + err.Error())
	}

	var assessmentMetricRepo = assessmentRepositories.NewAssessmentMetricRepository()

	assessmentMetrics, err := assessmentMetricRepo.FindAssessmentMetricsByInterviewNumber(interviewNumber)
	if err != nil {
		return studentsList, errors.New("finding assessment metrics error: " + err.Error())
	}

	for _, interviewStudent := range interviewStudents {
		student := interviewStudent.Student

		var asts []responses.AssessmentResponse

		for _, assessmentMetric := range assessmentMetrics {

			assmt, err := assessmentRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(interviewStudent.StudentNumber, assessmentMetric.Number)
			if err != nil {
				return studentsList, errors.New("Failed to fetch assessments " + err.Error())
			}

			//GetAssessmentAreaByNumber

			assessmentArea, err := assessmentAreaRepositories.GetAssessmentAreaByNumber(assessmentMetric.AssessmentAreaNumber)
			if err != nil {
				return studentsList, errors.New("Failed to fetch assessment area " + err.Error())
			}

			asts = append(asts, responses.AssessmentResponse{
				Number:                 assmt.Number,
				AssessmentAreaName:     assessmentArea.Name,
				AssessmentMetricNumber: assessmentMetric.Number,
				Score:                  assmt.Score,
				ActualScore:            assmt.ActualScore,
				Comment:                assmt.Comment,
			})

		}

		var assessment assessmentModel.Assessment
		if assessmentMetricNumber != "" {
			assessment, err = assessmentRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, assessmentMetricNumber)
			if err != nil {
				return studentsList, errors.New("Failed to fetch assessments " + err.Error())
			}
		}
		studentsList = append(studentsList, responses2.SingleStudent{
			EnrolledStudent: enrollmentResponse.EnrolledStudent{
				Number:            student.Number,
				StudentOtherName:  student.OtherName,
				StudentSurname:    student.Surname,
				StudentSecondName: student.SecondName,
				Score:             assessment.Score,
				ActualScore:       assessment.ActualScore,
				Comment:           assessment.Comment,
				Status:            student.Status,
			},
			Assessments: asts,
		})
	}

	return studentsList, nil
}

func RemoveInterview(interviewNumber string) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		repo := repository2.NewInterviewRepoImpl(tx)

		interview, err := repo.FindInterviewByNumber(interviewNumber)
		if err != nil {
			return errors.New("error checking interview number: " + err.Error())
		}

		if interview.Number == "" {
			return errors.New("invalid interview number")
		}

		//FindInterviewStudentsByInterViewNumber
		students, err := repository2.FindInterviewStudentsByInterViewNumber(interview.Number, tx)
		if err != nil {
			return errors.New("checking interview students error: " + err.Error())
		}

		if len(students) > 0 {
			return errors.New("this interview has students added to it, you may need to remove them before deleting")
		}

		schedules, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumber(interview.Number, tx)
		if err != nil {
			return errors.New("getting old schedules error: " + err.Error())
		}

		for _, schedule := range schedules {
			err = interviewRepositories.RemoveInterviewAssessmentMetric(schedule, tx)
			if err != nil {
				return errors.New("failed removing schedule from interview: " + err.Error())
			}
		}

		err = repo.DeleteInterview(interview)
		if err != nil {
			return errors.New("failed to remove interview: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}
