package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func FindAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber string, assessmentMetricNumber string) (assessmentModel.Assessment, error) {

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
