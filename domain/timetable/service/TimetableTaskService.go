package service

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type TimetableTaskService interface {
	CreateTimetableTask(request requests.TimetableTaskRequest, userNumber string) error
	UpdateTimetableTask(taskNumber string, request requests.UpdateTimetableTaskRequest) error
	DeleteTimetableTask(taskNumber string) error
	GetTimetableTaskByNumber(taskNumber string) (models.TimetableTask, error)
	GetAllTimetableTasks(classNumber string) ([]models.TimetableTask, error)
}

type timetableTaskService struct {
	repo repository.TimetableTaskRepository
}

func NewTimetableTaskService(repo repository.TimetableTaskRepository) TimetableTaskService {
	return &timetableTaskService{repo: repo}
}

func (s *timetableTaskService) CreateTimetableTask(request requests.TimetableTaskRequest, userNumber string) error {
	if len(request.Tasks) == 0 {
		return errors.New("at least one task definition is required")
	}

	if len(request.ClassNumbers) == 0 {
		return errors.New("at least one class number is required")
	}

	for _, classNumber := range request.ClassNumbers {
		if classNumber == "" {
			return errors.New("class number cannot be empty")
		}

		for _, taskItem := range request.Tasks {
			if taskItem.Name == "" {
				return errors.New("task name cannot be empty")
			}
			if taskItem.Type == "" {
				return errors.New("task type cannot be empty")
			}

			task := models.TimetableTask{
				Number:      uuid.NewV4().String(),
				Name:        taskItem.Name,
				Type:        taskItem.Type,
				UserNumber:  userNumber,
				ClassNumber: classNumber,
			}

			if err := s.repo.SaveTimetableTask(task); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *timetableTaskService) UpdateTimetableTask(taskNumber string, request requests.UpdateTimetableTaskRequest) error {
	task, err := s.repo.FindTimetableTaskByNumber(taskNumber)
	if err != nil {
		return errors.New("error finding timetable task: " + err.Error())
	}

	if task.Number == "" {
		return errors.New("timetable task not found")
	}

	if request.Name != "" {
		task.Name = request.Name
	}
	if request.Type != "" {
		task.Type = request.Type
	}
	if request.ClassNumber != "" {
		task.ClassNumber = request.ClassNumber
	}

	return s.repo.UpdateTimetableTask(&task)
}

func (s *timetableTaskService) DeleteTimetableTask(taskNumber string) error {
	task, err := s.repo.FindTimetableTaskByNumber(taskNumber)
	if err != nil {
		return errors.New("error finding timetable task: " + err.Error())
	}

	if task.Number == "" {
		return errors.New("timetable task not found")
	}

	return s.repo.DeleteTimetableTask(task)
}

func (s *timetableTaskService) GetTimetableTaskByNumber(taskNumber string) (models.TimetableTask, error) {
	task, err := s.repo.FindTimetableTaskByNumber(taskNumber)
	if err != nil {
		return models.TimetableTask{}, errors.New("error finding timetable task: " + err.Error())
	}

	if task.Number == "" {
		return models.TimetableTask{}, errors.New("timetable task not found")
	}

	return task, nil
}

func (s *timetableTaskService) GetAllTimetableTasks(classNumber string) ([]models.TimetableTask, error) {
	return s.repo.FindAllTimetableTasks(classNumber)
}
