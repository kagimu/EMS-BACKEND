package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

type InterviewRepo interface {
	SaveInterview(i models.Interview) error
	UpdateInterview(i *models.Interview) error
	DeleteInterview(i models.Interview) error
	FindInterviewByName(name, campusClassNumber string) (models.Interview, error)
	FindInterviewByNumber(interviewNumber string) (models.Interview, error)
	FindAllInterviews(campusClassNumber string) ([]models.Interview, error)
}

// implement interview repository
type InterviewRepoImpl struct {
	DB *gorm.DB
}

func (i *InterviewRepoImpl) FindInterviewByName(name, campusClassNumber string) (models.Interview, error) {
	var interview models.Interview

	if err := i.DB.Where("name = ? AND campus_class_number = ?", name, campusClassNumber).First(&interview).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Interview{}, nil
		}

		return models.Interview{}, err

	}

	return interview, nil
}

func NewInterviewRepoImpl(db *gorm.DB) InterviewRepo {
	return &InterviewRepoImpl{DB: db}
}

func (i *InterviewRepoImpl) SaveInterview(interview models.Interview) error {
	return i.DB.Save(&interview).Error
}

func (i *InterviewRepoImpl) UpdateInterview(interview *models.Interview) error {
	return i.DB.Save(interview).Error
}

func (i *InterviewRepoImpl) DeleteInterview(interview models.Interview) error {
	return i.DB.Unscoped().Delete(&interview).Error
}

func (i *InterviewRepoImpl) FindInterviewByNumber(number string) (models.Interview, error) {
	var interview models.Interview

	if err := i.DB.Where("number = ?", number).First(&interview).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Interview{}, nil
		}

		return models.Interview{}, err
	}

	return interview, nil
}

func (i *InterviewRepoImpl) FindAllInterviews(campusClassNumber string) ([]models.Interview, error) {
	var interviews []models.Interview

	err := i.DB.Where("campus_class_number = ?", campusClassNumber).Find(&interviews).Error
	if err != nil {
		return nil, err
	}
	return interviews, nil
}
