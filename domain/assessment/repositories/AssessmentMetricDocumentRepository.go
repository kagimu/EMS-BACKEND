package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

type ScriptMarkingGuideRepository interface {
	SaveScriptMarkingGuide(scriptMarkingGuide models.ScriptMarkingGuide) error
	UseTx(tx *gorm.DB) *gorm.DB
}

type scriptMarkingGuideRepository struct {
	Db *gorm.DB
}

func (s *scriptMarkingGuideRepository) UseTx(tx *gorm.DB) *gorm.DB {
	s.Db = tx
	return s.Db
}

func (s *scriptMarkingGuideRepository) SaveScriptMarkingGuide(scriptMarkingGuide models.ScriptMarkingGuide) error {
	return s.Db.Save(&scriptMarkingGuide).Error
}

func NewScriptMarkingGuideRepository() ScriptMarkingGuideRepository {
	return &scriptMarkingGuideRepository{
		Db: config.DB,
	}
}
