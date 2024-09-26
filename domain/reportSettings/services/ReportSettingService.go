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

	err = CheckReportColumnDuplicates(reportColumnRequest.Columns)
	if err != nil {
		return err
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for _, reportColumn := range reportColumnRequest.Columns {
			return processReportColumn(reportColumn, currentUser, reportSetting, tx)
		}
		return nil
	})

	return err
}

func CheckReportColumnDuplicates(reportColumnRequestItems []requests.ReportColumnRequestItem) error {

	cleanColumnsByName := make([]requests.ReportColumnRequestItem, len(reportColumnRequestItems))
	for _, reportColumn := range reportColumnRequestItems {

		for _, cleanColumn := range cleanColumnsByName {
			if reportColumn.Name == cleanColumn.Name {
				return errors.New(fmt.Sprintf("Column name %s is already used", reportColumn.Name))
			}
		}

		cleanColumnsByName = append(cleanColumnsByName, reportColumn)
	}

	cleanColumnsByCode := make([]requests.ReportColumnRequestItem, len(reportColumnRequestItems))
	for _, reportColumn := range reportColumnRequestItems {

		for _, cleanColumn := range cleanColumnsByCode {
			if reportColumn.Code == cleanColumn.Code {
				return errors.New(fmt.Sprintf("Column code %d on %s is already used on %s", cleanColumn.Code, cleanColumn.Name, reportColumn.Name))
			}
		}

		cleanColumnsByCode = append(cleanColumnsByCode, reportColumn)
	}

	return nil
}

func processReportColumn(reportColumnRequest requests.ReportColumnRequestItem, currentUser string, reportSetting models.ReportSetting, tx *gorm.DB) error {

	oldColumn := models.ReportColumn{}

	if reportColumnRequest.ReportColumnNumber != "" {
		oldColumn_, err := reporsitory.GetReportColumnByNumber(reportColumnRequest.ReportColumnNumber)
		if err != nil {
			return errors.New("Error checking report column by number: " + err.Error())
		}

		if oldColumn_.Number == "" {
			return errors.New("Invalid report column number provided ")
		}

		oldColumn = oldColumn_
	}

	if reportColumnRequest.ReportColumnNumber == "" {
		newReportColumn := models.ReportColumn{
			Number:              uuid.NewV4().String(),
			ReportSettingNumber: reportSetting.Number,
			Name:                reportColumnRequest.Name,
			PercentageShare:     reportColumnRequest.PercentageShare,
			UserNumber:          currentUser,
			Code:                reportColumnRequest.Code,
		}

		fmt.Println("Saving new report column name: " + reportColumnRequest.Name)
		err := reporsitory.SaveReportColumn(newReportColumn, tx)

		if err != nil {
			return errors.New("Error saving report column: " + err.Error())
		}
	} else {

		oldColumn.Name = reportColumnRequest.Name
		oldColumn.PercentageShare = reportColumnRequest.PercentageShare
		oldColumn.UserNumber = currentUser
		oldColumn.Code = reportColumnRequest.Code

		fmt.Println("Updating report column name: " + reportColumnRequest.Name)

		err := reporsitory.SaveReportColumn(oldColumn, tx)

		if err != nil {
			return errors.New("Error updating report column: " + err.Error())
		}
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
