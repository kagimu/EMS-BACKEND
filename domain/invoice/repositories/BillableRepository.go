package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SaveBillable(billable models.Billable, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&billable).Error
}

func FindBillable(schoolNumber, studyPeriodNumber, classNumber, campusNumber, academicYearNumber, status, billType string, trn ...*gorm.DB) ([]models.Billable, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billable []models.Billable

	if err := db.
		Joins("JOIN study_periods ON study_periods.number = billables.study_period_number").
		Scopes(filters.ByStatus(status), filters.ByAnyField("bill_type", billType), filters.ByCampusNumber(campusNumber), filters.ByAnyField("academic_year_number", academicYearNumber), filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByAnyField("class_number", classNumber)).
		Distinct("billables.*").
		Where("school_number = ?", schoolNumber).Find(&billable).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billable, nil
		}

		return nil, err
	}

	return billable, nil
}

func RemoveBillable(billable models.Billable, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&billable).Error
}

func GetBillableByNumber(billableNumber string, trn ...*gorm.DB) (models.Billable, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billable models.Billable

	if err := db.
		Where("number = ?", billableNumber).First(&billable).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billable, nil
		}

		return billable, err
	}

	return billable, nil
}

func GetDuplicatedBill(studyPeriodNumber, classNumber, campusNumber, billType, modeOfStudyNumber, status string, trn ...*gorm.DB) (models.Billable, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billable models.Billable

	if err := db.
		Scopes(filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByAnyField("mode_of_study_number", modeOfStudyNumber), filters.ByAnyField("bill_type", billType), filters.ByAnyField("class_number", classNumber)).
		Where("status = ? AND campus_number = ?", status, campusNumber).First(&billable).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billable, nil
		}

		return billable, err
	}

	return billable, nil
}

func GetBillableByInvoiceNumber(invoiceNumber string, trn ...*gorm.DB) (models.Billable, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billable models.Billable

	if err := db.
		Joins("JOIN billable_invoices On billable_invoices.invoice_number = invoices.number").
		Where("invoice_number = ?", invoiceNumber).First(&billable).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billable, nil
		}

		return billable, err
	}

	return billable, nil
}
