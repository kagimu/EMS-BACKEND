package services

import (
	"errors"
	"fmt"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentMetricRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	assessmentRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	enrollmentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	reportSettingsRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/api/responses"
	"gorm.io/gorm"
)

// GetReport returns students with their scores per assessment schedule for the given assessment area (and optional filters).
func GetReport(req requests.StudentAssessmentAreaReportRequest) (*responses.StudentAssessmentAreaReportResponse, error) {
	//assessmentArea, err := assessmentAreaRepo.GetAssessmentAreaByNumber(req.AssessmentAreaNumber)
	//if err != nil {
	return nil, nil
	//}
	//if assessmentArea.Number == "" {
	//	return nil, errors.New("assessment area not found")
	//}
	//
	//metrics, err := assessmentMetricRepo.NewAssessmentMetricRepository().AllAssessmentMetric(
	//	req.AssessmentAreaNumber,
	//	req.StudyPeriodNumber,
	//	req.CampusClassNumber,
	//	"", "", "",
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("fetching assessment metrics: %w", err)
	//}
	//if len(metrics) == 0 {
	//	return &responses.StudentAssessmentAreaReportResponse{
	//		AssessmentAreaNumber: assessmentArea.Number,
	//		AssessmentAreaName:   assessmentArea.Name,
	//		Schedules:            []responses.AssessmentScheduleSummary{},
	//		Students:             []responses.StudentWithScores{},
	//	}, nil
	//}
	//
	//metricNumbers := make([]string, 0, len(metrics))
	//campusClassNumbers := make(map[string]struct{})
	//scheduleSummaries := make([]responses.AssessmentScheduleSummary, 0, len(metrics))
	//for i := range metrics {
	//	metricNumbers = append(metricNumbers, metrics[i].Number)
	//	campusClassNumbers[metrics[i].CampusClassNumber] = struct{}{}
	//	scheduleSummaries = append(scheduleSummaries, responses.AssessmentScheduleSummary{
	//		AssessmentMetricNumber: metrics[i].Number,
	//		Description:            metrics[i].Description,
	//		ScheduledOn:            metrics[i].ScheduledOn.Format(time.RFC3339),
	//		MaxScore:               metrics[i].MaxScore,
	//	})
	//}
	//
	//classList := make([]string, 0, len(campusClassNumbers))
	//for c := range campusClassNumbers {
	//	classList = append(classList, c)
	//}
	//
	//var enrollments []enrollmentModels.Enrollment
	//db := config.DB.Where("campus_class_number IN ? AND status = ?", classList, constants.StatusActive)
	//if req.StudyPeriodNumber != "" {
	//	db = db.Where("study_period_number = ?", req.StudyPeriodNumber)
	//}
	//if err := db.Preload("Student").Preload("CampusClass").Find(&enrollments).Error; err != nil {
	//	return nil, fmt.Errorf("fetching enrollments: %w", err)
	//}
	//
	//studentNumbers := make([]string, 0, len(enrollments))
	//studentMap := make(map[string]struct {
	//	student studentModels.Student
	//	class   string
	//})
	//for i := range enrollments {
	//	studentNumbers = append(studentNumbers, enrollments[i].StudentNumber)
	//	className := ""
	//	if enrollments[i].CampusClass.Number != "" {
	//		className = enrollments[i].CampusClass.ClassName
	//	}
	//	studentMap[enrollments[i].StudentNumber] = struct {
	//		student studentModels.Student
	//		class   string
	//	}{student: enrollments[i].Student, class: className}
	//}
	//
	//assessments, err := assessmentRepo.FindAssessmentsByMetricNumbersAndStudentNumbers(metricNumbers, studentNumbers)
	//if err != nil {
	//	return nil, fmt.Errorf("fetching assessments: %w", err)
	//}
	//
	//scoreMap := make(map[string]map[string]assessmentModel.Assessment)
	//for _, a := range assessments {
	//	if scoreMap[a.StudentNumber] == nil {
	//		scoreMap[a.StudentNumber] = make(map[string]assessmentModel.Assessment)
	//	}
	//	scoreMap[a.StudentNumber][a.AssessmentMetricNumber] = a
	//}
	//
	//studentsResp := make([]responses.StudentWithScores, 0, len(studentMap))
	//for studentNumber, info := range studentMap {
	//	scores := make([]responses.AssessmentScheduleScore, 0, len(metrics))
	//	for _, m := range metrics {
	//		row := responses.AssessmentScheduleScore{
	//			AssessmentMetricNumber: m.Number,
	//			Description:            m.Description,
	//			ScheduledOn:            m.ScheduledOn.Format(time.RFC3339),
	//			MaxScore:               m.MaxScore,
	//		}
	//		if a, ok := scoreMap[studentNumber][m.Number]; ok {
	//			row.ActualScore = a.ActualScore
	//			row.Score = a.Score
	//			row.Comment = a.Comment
	//		}
	//		scores = append(scores, row)
	//	}
	//	studentsResp = append(studentsResp, responses.StudentWithScores{
	//		StudentNumber:   studentNumber,
	//		Surname:         info.student.Surname,
	//		SecondName:      info.student.SecondName,
	//		OtherName:       info.student.OtherName,
	//		ADN:             info.student.ADN,
	//		APN:             info.student.APN,
	//		CampusClassName: info.class,
	//		Scores:          scores,
	//	})
	//}
	//
	//return &responses.StudentAssessmentAreaReportResponse{
	//	AssessmentAreaNumber: assessmentArea.Number,
	//	AssessmentAreaName:   assessmentArea.Name,
	//	Schedules:            scheduleSummaries,
	//	Students:             studentsResp,
	//}, nil

}

func GetAllStudentScoresReportByClassNumberAndStudyPeriod(reqCampusClassNumber, reqStudyPeriodNumber, campusNumber, curriculumNumber, campusCurriculumNumber string) (*[]responses.ClassStudentAssessmentReportResponse, error) {

	CampusClasses, err := repositories2.GetCampusClassByClassNumberAndCampusNumber(reqCampusClassNumber, campusNumber, curriculumNumber, campusCurriculumNumber)
	if err != nil {
		return nil, fmt.Errorf("fetching CampusClassByClassNumberAndCampusNumber: %w", err)
	}

	campusClassNumber := CampusClasses[0].Number

	var enrollments []enrollmentModels.Enrollment

	if err := config.DB.Where("campus_class_number  =? AND status = ? AND study_period_number = ? ", campusClassNumber, constants.StatusActive, reqStudyPeriodNumber).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("fetching enrollments: %w", err)
	}

	allAssessmentSchedulesByClassAndStudyPeriod, err := assessmentMetricRepo.NewAssessmentMetricRepository().FindAssessmentMetricsByDistinctAssessmentAreasNumbers(
		"",
		reqStudyPeriodNumber,
		campusClassNumber,
		"", "", "",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch fetching assessment metrics: %w", err)
	}

	studentsResp := make([]responses.ClassStudentAssessmentReportResponse, 0)

	for _, enrollment := range enrollments {

		as := make([]responses.AssessmentArea, 0)
		for _, amv := range allAssessmentSchedulesByClassAndStudyPeriod {

			var AssessmentArea responses.AssessmentArea
			AssessmentArea.Name = amv.AssessmentArea.Name

			averageScore := 0.0
			assessmentCount := 0
			var schedules = make([]responses.Schedule, 0)
			assessmentMetric, err := assessmentMetricRepo.NewAssessmentMetricRepository().AllAssessmentMetric(
				amv.AssessmentAreaNumber,
				reqStudyPeriodNumber,
				campusClassNumber,
				"", "", "",
			)
			if err != nil {
				return nil, fmt.Errorf("fetching assessment assessmentMetric: %w", err)
			}

			for _, amcv := range assessmentMetric {
				rca, err := reportSettingsRepo.GetReportColumnAssessmentByAssessmentMetricNumber(amcv.Number)
				if err != nil {
					if !errors.Is(gorm.ErrRecordNotFound, err) {
						return nil, fmt.Errorf("fetching report column assessment: %w", err)
					}

				}
				if rca.Number == "" {
					continue
				}

				assessmentCount++

				assessments, err := assessmentRepo.FindAssessmentByStudentNumberAndAssessmentMetricNumber(enrollment.StudentNumber, amcv.Number)
				if err != nil {
					return nil, fmt.Errorf("fetching assessments: %w", err)
				}
				if len(assessments) > 0 {

					schedules = append(schedules, responses.Schedule{Name: amcv.Description, Score: assessments[0].Score, MaxScore: amcv.MaxScore})
					AssessmentArea.Scores = schedules

					averageScore += (assessments[0].Score / amcv.MaxScore)
				} else {
					averageScore += 0
				}
			}
			averageScore = (averageScore / float64(assessmentCount)) * 100

			AssessmentArea.Average = averageScore

			as = append(as, AssessmentArea)
		}
		student, err := repositories.FindStudentByNumber(enrollment.StudentNumber)
		if err != nil {
			return nil, err
		}
		studentsResp = append(studentsResp, responses.ClassStudentAssessmentReportResponse{
			StudentNumber: enrollment.StudentNumber,
			StudentName:   student.Surname + " " + student.SecondName + " " + student.OtherName,
			Assessments:   as,
		})
	}

	return &studentsResp, nil
}
