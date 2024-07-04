package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

type InterviewRepo interface {
	RemoveAssessmentMetricFromInterview()
	CreateInterview(i *models.Interview) error
	UpdateInterview(i *models.Interview) error
	DeleteInterview(i *models.Interview) error
	FindInterviewById(id string) (*models.Interview, error)
	FindAllInterviews() ([]*models.Interview, error)
}

// implement interview repository
type InterviewRepoImpl struct {
	DB *gorm.DB
}

func NewInterviewRepoImpl(db *gorm.DB) *InterviewRepoImpl {
	return &InterviewRepoImpl{DB: db}
}

func (i *InterviewRepoImpl) CreateInterview(interview *models.Interview) error {
	return i.DB.Create(interview).Error
}

func (i *InterviewRepoImpl) UpdateInterview(interview *models.Interview) error {
	return i.DB.Save(interview).Error
}

func (i *InterviewRepoImpl) DeleteInterview(interview *models.Interview) error {
	return i.DB.Delete(interview).Error
}

func (i *InterviewRepoImpl) FindInterviewByNumber(number string) (*models.Interview, error) {
	var interview models.Interview
	err := i.DB.Where("number = ?", number).First(&interview).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}

func (i *InterviewRepoImpl) FindAllInterviews() ([]*models.Interview, error) {
	var interviews []*models.Interview
	err := i.DB.Find(&interviews).Error
	if err != nil {
		return nil, err
	}
	return interviews, nil
}
