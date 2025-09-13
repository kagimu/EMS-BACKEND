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
		Order("billables.code ASC").
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

func GetEnrollmentBill(studyPeriodNumber, classNumber, campusNumber, billType, modeOfStudyNumber, curriculumNumber, status string, trn ...*gorm.DB) (models.Billable, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billable models.Billable

	if err := db.
		Scopes(filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByAnyField("mode_of_study_number", modeOfStudyNumber), filters.ByAnyField("bill_type", billType), filters.ByAnyField("class_number", classNumber)).
		Where("billables.status = ? AND billables.campus_number = ? AND billables.curriculum_number = ?", status, campusNumber, curriculumNumber).First(&billable).Error; err != nil {

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
		Joins("JOIN billable_invoices On billable_invoices.billable_number = billables.number").
		Where("billable_invoices.invoice_number = ?", invoiceNumber).First(&billable).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billable, nil
		}

		return billable, err
	}

	return billable, nil
}

func GetBillableTotal(billableNumber string) (float64, error) {

	db := config.DB

	var totalAmount float64

	if err := db.Model(&models.BillableItem{}).
		Where("billable_number = ? ", billableNumber).
		Select("COALESCE(sum(quantity * unit_price),0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}
