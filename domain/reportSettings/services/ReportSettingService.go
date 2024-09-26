package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

var assessmentRepository = repositories2.NewAssessmentMetricRepository()

func SaveReportColumn(reportColumnRequest requests.ReportColumnRequest, currentUser string) error {

	reportSetting, err := reporsitory.GetReportSettingByNumber(reportColumnRequest.ReportSettingNumber)
	if err != nil {
		return errors.New("Error checking report setting by number: " + reportColumnRequest.ReportSettingNumber)
	}

	if reportSetting.Number == "" {
		return errors.New("invalid report setting number")
	}

	if reportSetting.Columns < int64(len(reportColumnRequest.Columns)) {
		return errors.New("you should define upto the specified columns")
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for _, reportColumn := range reportColumnRequest.Columns {
			return processReportColumn(reportColumn, currentUser, reportSetting, tx)
		}
		return nil
	})

	return err
}

func processReportColumn(reportColumnRequest requests.ReportColumnRequestItem, currentUser string, reportSetting models.ReportSetting, tx *gorm.DB) error {

	checkName, err := reporsitory.GetReportColumnByNameAndReportSettingsNumber(reportColumnRequest.Name, reportSetting.Number, tx)
	if err != nil {
		return errors.New("Error checking report column by name: " + reportColumnRequest.Name)
	}

	if checkName.Number != "" {
		return errors.New(fmt.Sprintf("Column name '%s' already used on report setting with code '%d'", reportColumnRequest.Name, reportSetting.Code))
	}

	checkByCode, err := reporsitory.GetReportColumnByCodeAndReportSettingsNumber(reportColumnRequest.Code, reportSetting.Number, tx)
	if err != nil {
		return errors.New("Error checking column code: " + err.Error())
	}

	if checkByCode.Number != "" {
		return errors.New("Column code already used on: " + reportColumnRequest.Name)
	}

	reportColumns, err := reporsitory.FindReportColumnsByReportSettingsNumber(reportSetting.Number, tx)
	if err != nil {
		return errors.New("Error getting report columns: " + err.Error())
	}

	savedReportColumn := int64(len(reportColumns))

	if savedReportColumn >= reportSetting.Columns {
		return errors.New("report columns count reached")
	}

	totalPercentageShare, err := reporsitory.CountReportSettingTotalPercentageScores(reportSetting.Number, tx)
	if err != nil {
		return errors.New("Error getting report setting total percentage scores: " + err.Error())
	}

	savedColumnsPercentageShare := totalPercentageShare + reportColumnRequest.PercentageShare

	if (savedColumnsPercentageShare) > 100 {
		return errors.New("report settings columns exceeds 100%")
	}

	remainingPercentageScore := 100 - savedColumnsPercentageShare

	if remainingPercentageScore > reportColumnRequest.PercentageShare && (savedReportColumn+1) >= reportSetting.Columns {
		return errors.New("the report column '" + reportColumnRequest.Name + "' should cover upto percentage")
	}

	newReportColumn := models.ReportColumn{
		Number:              uuid.NewV4().String(),
		ReportSettingNumber: reportSetting.Number,
		Name:                reportColumnRequest.Name,
		PercentageShare:     reportColumnRequest.PercentageShare,
		UserNumber:          currentUser,
		Code:                reportColumnRequest.Code,
	}

	err = reporsitory.SaveReportColumn(newReportColumn, tx)
	if err != nil {
		return errors.New("Error saving report column: " + err.Error())
	}

	return nil
}

func SaveAssessmentToReportColumn(reportColumnAssessmentRequest requests.ReportColumnAssessmentRequest, currentUser string) error {

	reportColumn, err := reporsitory.GetReportColumnByNumber(reportColumnAssessmentRequest.ReportColumnNumber)
	if err != nil {
		return errors.New("Error checking report column by number: " + reportColumnAssessmentRequest.ReportColumnNumber)
	}

	if reportColumn.Number == "" {
		return errors.New("invalid report column number")
	}

	trn := config.DB.Begin()

	for _, assessmentMetricNumber := range reportColumnAssessmentRequest.AssessmentMetrics {
		err = addAssessmentToReportColumn(assessmentMetricNumber, reportColumn, reportColumnAssessmentRequest.AssessmentAreaNumber, currentUser, trn)
		if err != nil {
			trn.Rollback()
			return err
		}
	}

	if err := trn.Commit().Error; err != nil {
		return errors.New("Failed to finalize processing: " + err.Error())
	}

	return nil
}

func addAssessmentToReportColumn(assessmentMetricNumber string, reportColumn models.ReportColumn, assessmentAreaNumber string, currentUser string, trn *gorm.DB) error {

	assessmentMetric, err := reporsitory.GetReportColumnAssessmentByAssessmentMetricNumber(assessmentMetricNumber)
	if err != nil {
		return errors.New("Error checking assessment metric number: " + err.Error())
	}

	if assessmentMetric.Number != "" {
		return errors.New("The assessment metric number is already used: " + assessmentAreaNumber)
	}

	reportColumnAssessment := models.ReportColumnAssessment{
		Number:                 uuid.NewV4().String(),
		ReportColumnNumber:     reportColumn.Number,
		AssessmentAreaNumber:   assessmentAreaNumber,
		AssessmentMetricNumber: assessmentMetricNumber,
		UserNumber:             currentUser,
	}

	err = reporsitory.SaveReportColumnAssessment(reportColumnAssessment, trn)
	if err != nil {
		return errors.New("Error saving report column assessment: " + err.Error())
	}

	return nil
}

func GetReportColumns(reportSettingNumber string) ([]responses.ReportColumnResponse, error) {

	reportSetting, err := reporsitory.GetReportSettingByNumber(reportSettingNumber)
	if err != nil {
		return nil, errors.New("Error getting report setting: " + err.Error())
	}

	if reportSetting.Number == "" {
		return nil, errors.New("invalid report settings number")
	}

	columns, err := reporsitory.FindReportColumnsByReportSettingsNumber(reportSetting.Number)
	if err != nil {
		return nil, errors.New("Error getting report columns: " + err.Error())
	}

	var reportColumns = make([]responses.ReportColumnResponse, 0)

	for _, column := range columns {

		reportColumn := responses.ReportColumnResponse{
			Name:            column.Name,
			Number:          column.Number,
			PercentageShare: column.PercentageShare,
			Code:            column.Code,
		}

		reportColumnAssessments, err := reporsitory.FindReportColumnAssessmentsByReportColumnNumber(column.Number)
		if err != nil {
			return nil, errors.New("Error getting report column assessments: " + err.Error())
		}

		var reportColumnAssessmentResponse []responses.ReportColumnAssessmentResponse

		for _, reportColumnAssessment := range reportColumnAssessments {

			assessmentArea, err := repositories.GetAssessmentAreaByNumber(reportColumnAssessment.AssessmentAreaNumber)
			if err != nil {
				return nil, errors.New("Error getting report column assessment area: " + err.Error())
			}

			if assessmentArea.Number == "" {
				return nil, errors.New("invalid report column assessment area number")
			}

			assessmentSchedule, err := assessmentRepository.FindAssessmentMetricByNumber(reportColumnAssessment.AssessmentMetricNumber)
			if err != nil {
				return nil, errors.New("Error getting report column assessment metric number: " + err.Error())
			}

			if assessmentSchedule.Number == "" {
				return nil, errors.New("invalid report column assessment metric number")
			}

			metric, err := repositories2.GetMetricByNumber(assessmentSchedule.MetricNumber)
			if err != nil {
				return nil, errors.New("getting metric by number error: " + err.Error())
			}

			reportColumnAssessmentResponse = append(reportColumnAssessmentResponse, responses.ReportColumnAssessmentResponse{
				Number:           reportColumnAssessment.Number,
				AssessmentArea:   assessmentArea.Name,
				AssessmentMetric: metric.Name,
				ScheduledOn:      assessmentSchedule.ScheduledOn.Format(time.DateOnly),
			})
		}

		reportColumn.Assessments = reportColumnAssessmentResponse

		reportColumns = append(reportColumns, reportColumn)

	}

	return reportColumns, nil
}

func AllReportColumnAssessments(reportColumnNumber, assessmentAreaNumber, assessmentMetricNumber, academicYearNumber, studyPeriodNumber, campusClassNumber, reportSettingNumber string) ([]responses.ReportColumnAssessmentResponse, error) {

	var columnAssessmentResponse = make([]responses.ReportColumnAssessmentResponse, 0)

	allAssessments, err := reporsitory.AllReportColumnAssessments(reportColumnNumber, assessmentAreaNumber, assessmentMetricNumber, academicYearNumber, studyPeriodNumber, campusClassNumber, reportSettingNumber)
	if err != nil {
		return columnAssessmentResponse, errors.New("Getting report column assessments error " + err.Error())
	}

	for _, columnAssessment := range allAssessments {
		columnAssessmentResponse = append(columnAssessmentResponse, CreateReportColumnAssessmentResponse(columnAssessment))
	}

	return columnAssessmentResponse, nil
}

func CreateReportColumnAssessmentResponse(reportColumnAssessment models.ReportColumnAssessment) responses.ReportColumnAssessmentResponse {

	assessmentArea, _ := repositories.GetAssessmentAreaByNumber(reportColumnAssessment.AssessmentAreaNumber)

	assessmentSchedule, _ := assessmentRepository.FindAssessmentMetricByNumber(reportColumnAssessment.AssessmentMetricNumber)

	metric, _ := repositories2.GetMetricByNumber(assessmentSchedule.MetricNumber)

	reportColumn, _ := reporsitory.GetReportColumnByNumber(reportColumnAssessment.ReportColumnNumber)

	reportSetting, _ := reporsitory.GetReportSettingByNumber(reportColumn.ReportSettingNumber)

	class, _ := repositories3.GetClassByCampusClassNumber(reportSetting.CampusClassNumber)
	sp, _ := repositories4.GetStudyPeriodByNumber(reportSetting.StudyPeriodNumber)
	ay, _ := repositories5.GetAcademicYearByNumber(reportSetting.AcademicYearNumber)
	spn, _ := repositories4.GetStudyPeriodNameByStudyPeriodNumber(sp.Number)

	return responses.ReportColumnAssessmentResponse{
		Number:           reportColumnAssessment.Number,
		AssessmentArea:   assessmentArea.Name,
		AssessmentMetric: metric.Name,
		ReportColumn:     reportColumn.Name,
		ScheduledOn:      assessmentSchedule.ScheduledOn.Format(time.DateTime),
		Class:            class.Name,
		StudyPeriod:      spn.Name,
		AcademicYear:     ay.AcademicYearName,
	}
}
