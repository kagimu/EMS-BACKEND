package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"gorm.io/gorm"
)

type TimetableTaskRepository interface {
	SaveTimetableTask(task models.TimetableTask) error
	UpdateTimetableTask(task *models.TimetableTask) error
	DeleteTimetableTask(task models.TimetableTask) error
	FindTimetableTaskByNumber(number string) (models.TimetableTask, error)
	FindTimetableTaskByID(id uint) (models.TimetableTask, error)
	FindAllTimetableTasks(classNumber string) ([]models.TimetableTask, error)
}

type timetableTaskRepository struct {
	DB *gorm.DB
}

func NewTimetableTaskRepository(db *gorm.DB) TimetableTaskRepository {
	return &timetableTaskRepository{DB: db}
}

func (r *timetableTaskRepository) SaveTimetableTask(task models.TimetableTask) error {
	return r.DB.Save(&task).Error
}

func (r *timetableTaskRepository) UpdateTimetableTask(task *models.TimetableTask) error {
	return r.DB.Save(task).Error
}

func (r *timetableTaskRepository) DeleteTimetableTask(task models.TimetableTask) error {
	return r.DB.Unscoped().Delete(&task).Error
}

func (r *timetableTaskRepository) FindTimetableTaskByNumber(number string) (models.TimetableTask, error) {
	var task models.TimetableTask

	if err := r.DB.Where("number = ?", number).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TimetableTask{}, nil
		}
		return models.TimetableTask{}, err
	}

	return task, nil
}

func (r *timetableTaskRepository) FindTimetableTaskByID(id uint) (models.TimetableTask, error) {
	var task models.TimetableTask

	if err := r.DB.Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TimetableTask{}, nil
		}
		return models.TimetableTask{}, err
	}

	return task, nil
}

func (r *timetableTaskRepository) FindAllTimetableTasks(classNumber string) ([]models.TimetableTask, error) {
	var tasks []models.TimetableTask

	query := r.DB
	if classNumber != "" {
		query = query.Where("class_number = ?", classNumber)
	}

	err := query.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
