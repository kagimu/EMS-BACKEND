package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

func GetCurriculumClassFromCampusClassNumber(campusClassNumber string, trn ...*gorm.DB) (models.CurriculumClass, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var curriculumClass models.CurriculumClass

	if err := db.Where("campus_class_number = ?").Find(&curriculumClass).Error; err != nil {
		return curriculumClass, err
	}
}
