package reportSettings

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

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

	trn := config.DB.Begin()

	for _, reportColumn := range reportColumnRequest.Columns {

		err = processReportColumn(reportColumn, currentUser, reportSetting)
		if err != nil {
			trn.Rollback()
			return err
		}

	}

	if err := trn.Commit(); err != nil {
		return errors.New("failed to commit report columns transaction: " + err.Error.Error())
	}

	return nil
}

func processReportColumn(reportColumnRequest requests.ReportColumnRequestItem, currentUser string, reportSetting models.ReportSetting) error {

	checkName, err := reporsitory.GetReportColumnByNameAndReportSettingsNumber(reportColumnRequest.Name, reportSetting.Number)
	if err != nil {
		return errors.New("Error checking report column by name: " + reportColumnRequest.Name)
	}

	if checkName.Name != "" {
		return errors.New("Column name already used: " + reportColumnRequest.Name)
	}

	checkByCode, err := reporsitory.GetReportColumnByCodeAndReportSettingsNumber(reportColumnRequest.Code, reportSetting.Number)
	if err != nil {
		return errors.New("Error checking column code: " + err.Error())
	}

	if checkByCode.Number != "" {
		return errors.New("Column code already used: " + reportColumnRequest.Name)
	}

	reportColumns, err := reporsitory.FindReportColumnsByReportSettingsNumber(reportSetting.Number)
	if err != nil {
		return errors.New("Error getting report columns: " + err.Error())
	}

	savedReportColumn := int64(len(reportColumns))

	if savedReportColumn >= reportSetting.Columns {
		return errors.New("report columns count reached")
	}

	totalPercentageShare, err := reporsitory.CountReportSettingTotalPercentageScores(reportSetting.Number)
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

	err = reporsitory.SaveReportColumn(newReportColumn)
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

	if err := trn.Commit(); err != nil {
		return errors.New("Failed to finalize processing: " + err.Error.Error())
	}

	return nil
}

func addAssessmentToReportColumn(assessmentMetricNumber string, reportColumn models.ReportColumn, assessmentAreaNumber string, currentUser string, trn *gorm.DB) error {

	assessmentMetric, err := reporsitory.GetReportColumnAssessmentByAssessmentMetricNumber(assessmentMetricNumber)
	if err != nil {
		return errors.New("Error checking assessment metric number: " + assessmentMetricNumber)
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
