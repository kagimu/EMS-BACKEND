package reportSettings

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	uuid "github.com/satori/go.uuid"
)

func SaveReportColumn(reportColumnRequest requests.ReportColumnRequest, currentUser string) error {

	reportSetting, err := reporsitory.GetReportSettingByNumber(reportColumnRequest.ReportSettingNumber)
	if err != nil {
		return errors.New("Error checking report setting by number: " + reportColumnRequest.ReportSettingNumber)
	}

	if reportSetting.Number == "" {
		return errors.New("invalid report setting number")
	}

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
		return errors.New("this report column should cover upto percentage 100%")
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
