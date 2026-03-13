package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentMetricRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	assessmentRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assessmentAreaRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	enrollmentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/api/responses"
)

// GetReport returns students with their scores per assessment schedule for the given assessment area (and optional filters).
func GetReport(req requests.StudentAssessmentAreaReportRequest) (*responses.StudentAssessmentAreaReportResponse, error) {
	assessmentArea, err := assessmentAreaRepo.GetAssessmentAreaByNumber(req.AssessmentAreaNumber)
	if err != nil {
		return nil, err
	}
	if assessmentArea.Number == "" {
		return nil, errors.New("assessment area not found")
	}

	metrics, err := assessmentMetricRepo.NewAssessmentMetricRepository().AllAssessmentMetric(
		req.AssessmentAreaNumber,
		req.StudyPeriodNumber,
		req.CampusClassNumber,
		"", "", "",
	)
	if err != nil {
		return nil, fmt.Errorf("fetching assessment metrics: %w", err)
	}
	if len(metrics) == 0 {
		return &responses.StudentAssessmentAreaReportResponse{
			AssessmentAreaNumber: assessmentArea.Number,
			AssessmentAreaName:   assessmentArea.Name,
			Schedules:           []responses.AssessmentScheduleSummary{},
			Students:            []responses.StudentWithScores{},
		}, nil
	}

	metricNumbers := make([]string, 0, len(metrics))
	campusClassNumbers := make(map[string]struct{})
	scheduleSummaries := make([]responses.AssessmentScheduleSummary, 0, len(metrics))
	for i := range metrics {
		metricNumbers = append(metricNumbers, metrics[i].Number)
		campusClassNumbers[metrics[i].CampusClassNumber] = struct{}{}
		scheduleSummaries = append(scheduleSummaries, responses.AssessmentScheduleSummary{
			AssessmentMetricNumber: metrics[i].Number,
			Description:            metrics[i].Description,
			ScheduledOn:            metrics[i].ScheduledOn.Format(time.RFC3339),
			MaxScore:               metrics[i].MaxScore,
		})
	}

	classList := make([]string, 0, len(campusClassNumbers))
	for c := range campusClassNumbers {
		classList = append(classList, c)
	}

	var enrollments []enrollmentModels.Enrollment
	db := config.DB.Where("campus_class_number IN ? AND status = ?", classList, constants.StatusActive)
	if req.StudyPeriodNumber != "" {
		db = db.Where("study_period_number = ?", req.StudyPeriodNumber)
	}
	if err := db.Preload("Student").Preload("CampusClass").Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("fetching enrollments: %w", err)
	}

	studentNumbers := make([]string, 0, len(enrollments))
	studentMap := make(map[string]struct {
		student studentModels.Student
		class   string
	})
	for i := range enrollments {
		studentNumbers = append(studentNumbers, enrollments[i].StudentNumber)
		className := ""
		if enrollments[i].CampusClass.Number != "" {
			className = enrollments[i].CampusClass.ClassName
		}
		studentMap[enrollments[i].StudentNumber] = struct {
			student studentModels.Student
			class   string
		}{student: enrollments[i].Student, class: className}
	}

	assessments, err := assessmentRepo.FindAssessmentsByMetricNumbersAndStudentNumbers(metricNumbers, studentNumbers)
	if err != nil {
		return nil, fmt.Errorf("fetching assessments: %w", err)
	}

	scoreMap := make(map[string]map[string]assessmentModel.Assessment)
	for _, a := range assessments {
		if scoreMap[a.StudentNumber] == nil {
			scoreMap[a.StudentNumber] = make(map[string]assessmentModel.Assessment)
		}
		scoreMap[a.StudentNumber][a.AssessmentMetricNumber] = a
	}

	studentsResp := make([]responses.StudentWithScores, 0, len(studentMap))
	for studentNumber, info := range studentMap {
		scores := make([]responses.AssessmentScheduleScore, 0, len(metrics))
		for _, m := range metrics {
			row := responses.AssessmentScheduleScore{
				AssessmentMetricNumber: m.Number,
				Description:           m.Description,
				ScheduledOn:          m.ScheduledOn.Format(time.RFC3339),
				MaxScore:             m.MaxScore,
			}
			if a, ok := scoreMap[studentNumber][m.Number]; ok {
				row.ActualScore = a.ActualScore
				row.Score = a.Score
				row.Comment = a.Comment
			}
			scores = append(scores, row)
		}
		studentsResp = append(studentsResp, responses.StudentWithScores{
			StudentNumber:   studentNumber,
			Surname:         info.student.Surname,
			SecondName:      info.student.SecondName,
			OtherName:       info.student.OtherName,
			ADN:             info.student.ADN,
			APN:             info.student.APN,
			CampusClassName: info.class,
			Scores:          scores,
		})
	}

	return &responses.StudentAssessmentAreaReportResponse{
		AssessmentAreaNumber: assessmentArea.Number,
		AssessmentAreaName:   assessmentArea.Name,
		Schedules:           scheduleSummaries,
		Students:            studentsResp,
	}, nil
}
