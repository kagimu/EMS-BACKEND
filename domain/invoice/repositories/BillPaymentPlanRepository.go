package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SaveBillPaymentPlan(plan models.BillPaymentPlan, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	return db.Save(&plan).Error
}

func SaveBillPaymentPlanInstallment(installment models.BillPaymentPlanInstallment, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	return db.Save(&installment).Error
}

func FindBillPaymentPlanByNumber(number string, trn ...*gorm.DB) (models.BillPaymentPlan, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var plan models.BillPaymentPlan
	if err := db.
		Preload("Installments").
		Where("number = ?", number).
		First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BillPaymentPlan{}, nil
		}
		return models.BillPaymentPlan{}, err
	}

	return plan, nil
}

func FindBillPaymentPlans(planNumber, billableNumber, studentNumber, status string, trn ...*gorm.DB) ([]models.BillPaymentPlan, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var plans []models.BillPaymentPlan
	if err := db.
		Scopes(
			filters.ByAnyField("number", planNumber),
			filters.ByAnyField("billable_number", billableNumber),
			filters.ByAnyField("student_number", studentNumber),
			filters.ByAnyField("status", status),
		).
		Preload("Installments").
		Order("created_at desc").
		Find(&plans).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]models.BillPaymentPlan, 0), nil
		}
		return nil, err
	}

	return plans, nil
}

func DeleteBillPaymentPlanInstallmentsByPlanNumber(planNumber string, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Where("bill_payment_plan_number = ?", planNumber).Delete(&models.BillPaymentPlanInstallment{}).Error
}

func FindActiveBillPaymentPlanByStudentAndBillable(studentNumber, billableNumber string, trn ...*gorm.DB) (models.BillPaymentPlan, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var plan models.BillPaymentPlan
	if err := db.Where("student_number = ? AND billable_number = ? AND status = ?", studentNumber, billableNumber, "ACTIVE").First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BillPaymentPlan{}, nil
		}
		return models.BillPaymentPlan{}, err
	}

	return plan, nil
}

