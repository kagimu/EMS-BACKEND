package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

type subjectGroupRepository struct {
	DB *gorm.DB
}

func (r *subjectGroupRepository) SaveGroupAssessmentArea(group models.GroupAssessmentArea) error {
	return r.DB.Create(&group).Error
}

type SubjectGroupRepository interface {
	Save(group models.SubjectGroup) error
	SaveGroupAssessmentArea(group models.GroupAssessmentArea) error
	Update(group models.SubjectGroup) error
}

func (r *subjectGroupRepository) Save(group models.SubjectGroup) error {
	return r.DB.Create(&group).Error
}

func (r *subjectGroupRepository) Update(group models.SubjectGroup) error {
	return r.DB.Save(&group).Error
}

func NewSubjectGroupRepository(db *gorm.DB) SubjectGroupRepository {
	return &subjectGroupRepository{
		DB: db,
	}
}
