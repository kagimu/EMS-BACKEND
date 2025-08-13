package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

type subjectGroupRepository struct {
	DB *gorm.DB
}

func (r *subjectGroupRepository) FindAssessmentAreasByGroupNumber(number string) ([]models.GroupAssessmentArea, error) {
	var assessmentAreas = make([]models.GroupAssessmentArea, 0)
	err := r.DB.Where("group_number = ?", number).Find(&assessmentAreas).Error
	return assessmentAreas, err
}

func (r *subjectGroupRepository) DeleteGroup(number string) error {
	return r.DB.Delete(&models.SubjectGroup{}, "number = ?", number).Error
}

func (r *subjectGroupRepository) FindSubjectGroupByNumber(number string) (*models.SubjectGroup, error) {
	var sg models.SubjectGroup
	err := r.DB.Where("number = ?", number).First(&sg).Error
	return &sg, err
}

func (r *subjectGroupRepository) FindSubjectGroupsByReportSettingNumber(number string) ([]models.SubjectGroup, error) {
	var sgs []models.SubjectGroup
	err := r.DB.Where("report_settings_number = ?", number).Find(&sgs).Error
	return sgs, err
}

func (r *subjectGroupRepository) SaveGroupAssessmentArea(group models.GroupAssessmentArea) error {
	return r.DB.Create(&group).Error
}

type SubjectGroupRepository interface {
	Save(group models.SubjectGroup) error
	SaveGroupAssessmentArea(group models.GroupAssessmentArea) error
	Update(group models.SubjectGroup) error
	FindSubjectGroupByNumber(number string) (*models.SubjectGroup, error)
	FindSubjectGroupsByReportSettingNumber(number string) ([]models.SubjectGroup, error)
	DeleteGroup(number string) error
	FindAssessmentAreasByGroupNumber(number string) ([]models.GroupAssessmentArea, error)
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
