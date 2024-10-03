package requests

import (
	"errors"
)

type UpdateReportSettingsRequest struct {
	ReportSettingsRequest
	ReportSettingNumber string `json:"report_setting_number"`
}

type ReportSettingsRequest struct {
	Code                   int64  `json:"code"`
	AcademicYearNumber     string `json:"academic_year_number"`
	StudyPeriodNumber      string `json:"study_period_number"`
	ClassNumber            string `json:"class_number"`
	CampusNumber           string `json:"campus_number"`
	CurriculumNumber       string `json:"curriculum_number"`
	CampusCurriculumNumber string `json:"campus_curriculum_number"`
	CampusClassNumber      string `json:"campus_class_number"`
	Columns                int64  `json:"columns"`
}

type ReportColumnRequest struct {
	ReportSettingNumber string                    `json:"report_setting_number"`
	Columns             []ReportColumnRequestItem `json:"columns"`
}

type ReportColumnRequestItem struct {
	Code               int64   `json:"code"`
	Name               string  `json:"name"`
	PercentageShare    float64 `json:"percentage_share"`
	ReportColumnNumber string  `json:"report_column_number"`
}

type ReportColumnAssessmentRequest struct {
	Code                 int64    `json:"code"`
	ReportColumnNumber   string   `json:"report_column_number"`
	AssessmentAreaNumber string   `json:"assessment_area_number"`
	AssessmentMetrics    []string `json:"assessment_schedule"`
}

type ReportSummaryRequest struct {
	Code                int64   `json:"code"`
	ReportSettingNumber string  `json:"report_setting_number"`
	StudentNumber       string  `json:"student_number"`
	CampusClassNumber   string  `json:"campus_class_number"`
	StudyPeriodNumber   string  `json:"study_period_number"`
	Grade               string  `json:"grade"`
	Aggregate           float64 `json:"aggregate"`
	Comment             string  `json:"comment"`
}

type ReportScoreSummaryRequest struct {
	Code                 int64   `json:"code"`
	ReportSummaryNumber  string  `json:"report_summary_number"`
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	Aggregate            float64 `json:"aggregate"`
	Grade                string  `json:"grade"`
	Summation            float64 `json:"summation"`
}

type ReportScoreRequest struct {
	Code                 int64   `json:"code"`
	ReportSummaryNumber  string  `json:"report_summary_number"`
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	ReportColumnNumber   string  `json:"report_column_number"`
	Score                float64 `json:"score"`
	Aggregate            float64 `json:"aggregate"`
	Grade                string  `json:"grade"`
}

type GenerateReportSummaryRequest struct {
	StudentNumber     string `json:"student_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	CampusClassNumber string `json:"campus_class_number"`
}

func (reportSummary *ReportSummaryRequest) ValidateReportSummary() error {
	if reportSummary.ReportSettingNumber == "" {
		return errors.New("report setting is required")
	}
	if reportSummary.StudentNumber == "" {
		return errors.New("student is required")
	}
	if reportSummary.CampusClassNumber == "" {
		return errors.New("class is required")
	}
	if reportSummary.StudyPeriodNumber == "" {
		return errors.New("study period is required")
	}
	if reportSummary.Grade == "" {
		return errors.New("grade is required")
	}
	if reportSummary.Aggregate == 0 {
		return errors.New("aggregate is required")
	}

	if reportSummary.Code == 0 {
		return errors.New("code is required")

	}

	return nil

}

func (reportSetting *ReportSettingsRequest) ValidateReportSetting() error {
	if reportSetting.AcademicYearNumber == "" {
		return errors.New("academic year is required")
	}
	if reportSetting.StudyPeriodNumber == "" {
		return errors.New("study period is required")
	}

	if reportSetting.CampusNumber == "" {
		return errors.New("Choose a campus ")
	}

	if reportSetting.ClassNumber == "" {
		return errors.New("Choose a class ")
	}
	if reportSetting.Columns == 0 {
		return errors.New("number of columns is required")
	}
	if reportSetting.Code == 0 {
		return errors.New("code is required")
	}
	return nil
}

func (reportScoreSummary *ReportScoreSummaryRequest) ValidateReportScoreSummary() error {
	if reportScoreSummary.ReportSummaryNumber == "" {
		return errors.New("report summary is required")
	}
	if reportScoreSummary.AssessmentAreaNumber == "" {
		return errors.New("assessment area is required")
	}
	if reportScoreSummary.Aggregate == 0 {
		return errors.New("aggregate is required")
	}
	if reportScoreSummary.Grade == "" {
		return errors.New("grade is required")
	}
	if reportScoreSummary.Summation == 0 {
		return errors.New("summation is required")
	}
	if reportScoreSummary.Code == 0 {
		return errors.New("code is required")

	}
	return nil
}

func (reportScore *ReportScoreRequest) ValidateReportScore() error {
	if reportScore.ReportSummaryNumber == "" {
		return errors.New("report summary is required")
	}
	if reportScore.AssessmentAreaNumber == "" {
		return errors.New("assessment area is required")
	}
	if reportScore.ReportColumnNumber == "" {
		return errors.New("report column is required")
	}
	if reportScore.Score == 0 {
		return errors.New("score is required")
	}
	if reportScore.Aggregate == 0 {
		return errors.New("aggregate is required")
	}
	if reportScore.Grade == "" {
		return errors.New("grade is required")
	}

	if reportScore.Code == 0 {
		return errors.New("code is required")

	}
	return nil

}

func (reportColumnAssessment *ReportColumnAssessmentRequest) ValidateReportColumnAssessment() error {
	if reportColumnAssessment.ReportColumnNumber == "" {
		return errors.New("report column is required")
	}
	if reportColumnAssessment.AssessmentAreaNumber == "" {
		return errors.New("assessment area is required")
	}
	if len(reportColumnAssessment.AssessmentMetrics) == 0 {
		return errors.New("select some assessment schedules")
	}

	if reportColumnAssessment.Code == 0 {
		return errors.New("code is required")
	}

	return nil

}

func (reportColumnRequest *ReportColumnRequest) ValidateReportColumn() error {
	if reportColumnRequest.ReportSettingNumber == "" {
		return errors.New("report setting is required")
	}

	columns := reportColumnRequest.Columns
	if columns == nil {
		return errors.New("set columns list")
	}

	var totalShare float64

	for _, reportColumn := range columns {

		if reportColumn.Name == "" {
			return errors.New("name is required")
		}

		if reportColumn.PercentageShare <= 0 {
			return errors.New("percentage share is required")
		}

		if reportColumn.Code <= 0 {
			return errors.New("code is required")
		}

		totalShare += reportColumn.PercentageShare
	}

	if totalShare != 100 {
		return errors.New("Columns percentage total sum should be 100 ")
	}

	return nil
}

func (rpt *GenerateReportSummaryRequest) ValidateGenerateReportSummary() error {
	if rpt.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if rpt.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if rpt.CampusClassNumber == "" {
		return errors.New("class number is required")
	}
	return nil
}
