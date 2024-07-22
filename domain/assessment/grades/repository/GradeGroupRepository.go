package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"gorm.io/gorm"
)

type GradeGroupRepository interface {
	Create(gradeGroup *models.GradeGroup) error
	Update(gradeGroup *models.GradeGroup) error
	Delete(gradeGroup *models.GradeGroup) error
	FindAllBySchoolNumber(schoolNumber string) ([]models.GradeGroup, error)
	FindGradeGroupByNumber(gradeGroupNumber string) (models.GradeGroup, error)
}

type GradeGroupRepositoryImpl struct {
	DB *gorm.DB
}

func (g GradeGroupRepositoryImpl) FindGradeGroupByNumber(gradeGroupNumber string) (models.GradeGroup, error) {
	var gradeGroup models.GradeGroup
	if err := g.DB.Where("number = ?", gradeGroupNumber).First(&gradeGroup).Error; err != nil {
		return gradeGroup, err
	}
	return gradeGroup, nil
}

func (g GradeGroupRepositoryImpl) FindAllBySchoolNumber(schoolNumber string) ([]models.GradeGroup, error) {
	gradeGroups := make([]models.GradeGroup, 0)
	if err := g.DB.Where("school_number = ?", schoolNumber).Find(&gradeGroups).Error; err != nil {
		return nil, err
	}
	return gradeGroups, nil
}

func (g GradeGroupRepositoryImpl) Create(gradeGroup *models.GradeGroup) error {
	return g.DB.Create(gradeGroup).Error
}

func (g GradeGroupRepositoryImpl) Update(gradeGroup *models.GradeGroup) error {
	return g.DB.Save(gradeGroup).Error
}

func (g GradeGroupRepositoryImpl) Delete(gradeGroup *models.GradeGroup) error {
	return g.DB.Delete(gradeGroup).Error
}

func NewGradeGroupRepositoryImpl(db *gorm.DB) GradeGroupRepository {
	return &GradeGroupRepositoryImpl{DB: db}
}
