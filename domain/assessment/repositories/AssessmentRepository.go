package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber string, assessmentMetricNumber string) (assessmentModel.Assessment, error) {

	db := config.DB

	var assessment assessmentModel.Assessment

	if err := db.Scopes(filters.ByStudentNumber(studentNumber), filters.ByAssessmentMetricNumber(assessmentMetricNumber)).First(&assessment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessment, nil
		}

		return assessment, err
	}

	return assessment, nil
}

func SaveAssessment(assessment assessmentModel.Assessment) error {

	db := config.DB

	if err := db.Save(&assessment).Error; err != nil {
		return err
	}

	return nil

}

func FindAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber string, assessmentMetricNumber string) ([]assessmentModel.Assessment, error) {

	db := config.DB

	var assessment = make([]assessmentModel.Assessment, 0)

	if err := db.Scopes(filters.ByStudentNumber(studentNumber), filters.ByAssessmentMetricNumber(assessmentMetricNumber)).Find(&assessment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessment, nil
		}

		return assessment, err
	}

	return assessment, nil
}

func GetAssessmentByNumber(assessmentNumber string, trn ...*gorm.DB) (assessmentModel.Assessment, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessment assessmentModel.Assessment

	if err := db.Where("number = ?", assessmentNumber).First(&assessment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessment, nil
		}

		return assessment, err
	}

	return assessment, nil
}
