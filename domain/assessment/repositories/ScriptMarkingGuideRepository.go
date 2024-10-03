package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

type ScriptMarkingGuideRepository interface {
	SaveScriptMarkingGuide(scriptMarkingGuide models.ScriptMarkingGuide) error
	UseTx(tx *gorm.DB) *gorm.DB
	AllScriptMarkingGuide(assessmentMetricNumber string) ([]models.ScriptMarkingGuide, error)
}

type scriptMarkingGuideRepository struct {
	Db *gorm.DB
}

func (scriptMarkingGuideRepository *scriptMarkingGuideRepository) AllScriptMarkingGuide(assessmentMetricNumber string) ([]models.ScriptMarkingGuide, error) {

	var scriptMarkingGuides []models.ScriptMarkingGuide
	if err := scriptMarkingGuideRepository.Db.
		Where("assessment_metric_number = ?", assessmentMetricNumber).
		Find(&scriptMarkingGuides).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.ScriptMarkingGuide{}, nil
		}

		return []models.ScriptMarkingGuide{}, err
	}

	return scriptMarkingGuides, nil
}

func (scriptMarkingGuideRepository *scriptMarkingGuideRepository) UseTx(tx *gorm.DB) *gorm.DB {
	scriptMarkingGuideRepository.Db = tx
	return scriptMarkingGuideRepository.Db
}

func (scriptMarkingGuideRepository *scriptMarkingGuideRepository) SaveScriptMarkingGuide(scriptMarkingGuide models.ScriptMarkingGuide) error {
	return scriptMarkingGuideRepository.Db.Save(&scriptMarkingGuide).Error
}

func NewScriptMarkingGuideRepository() ScriptMarkingGuideRepository {
	return &scriptMarkingGuideRepository{
		Db: config.DB,
	}
}
