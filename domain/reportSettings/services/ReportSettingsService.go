package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	uuid "github.com/satori/go.uuid"
)

func SaveReportSettings(reportSettingsRequest requests.ReportSettingsRequest, currentUserNumber string) error {

	existingReportSettings, err := reporsitory.FindDuplicateReportSettings(reportSettingsRequest.StudyPeriodNumber, reportSettingsRequest.AcademicYearNumber, reportSettingsRequest.CampusClassNumber, constants.StatusActive)

	if err != nil {
		return errors.New("error checking report duplicate: " + err.Error())
	}

	if existingReportSettings.Number != "" {
		return errors.New("Report settings for this class  already exists")
	}

	// Save report settings
	reportSettings := models.ReportSetting{
		Number:             uuid.NewV4().String(),
		StudyPeriodNumber:  reportSettingsRequest.StudyPeriodNumber,
		AcademicYearNumber: reportSettingsRequest.AcademicYearNumber,
		CampusClassNumber:  reportSettingsRequest.CampusClassNumber,
		Columns:            reportSettingsRequest.Columns,
		Status:             constants.StatusActive,
		UserNumber:         currentUserNumber,
		Code:               reportSettingsRequest.Code,
	}

	if err = reporsitory.SaveReportSettings(reportSettings); err != nil {
		return errors.New("Failed to create report settings" + err.Error())
	}

	return nil
}
