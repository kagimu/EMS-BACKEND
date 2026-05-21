package service

import (
	"errors"

	campusRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	schoolRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	levelRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type TimetableTaskService interface {
	CreateTimetableTask(request requests.TimetableTaskRequest, userNumber string) error
	UpdateTimetableTask(taskNumber string, request requests.UpdateTimetableTaskRequest) error
	DeleteTimetableTask(taskNumber string) error
	GetTimetableTaskByNumber(taskNumber string) (models.TimetableTask, error)
	GetAllTimetableTasks(classNumber string) ([]responses.TimetableTaskResponse, error)
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
				Number:                 uuid.NewV4().String(),
				Name:                   taskItem.Name,
				Code:                   taskItem.Code,
				Type:                   taskItem.Type,
				UserNumber:             userNumber,
				ClassNumber:            classNumber,
				SchoolNumber:           request.SchoolNumber,
				CampusNumber:           request.CampusNumber,
				CampusCurriculumNumber: request.CampusCurriculumNumber,
				LevelNumber:            request.LevelNumber,
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

	name := request.Name
	code := request.Code
	taskType := request.Type
	classNumber := request.ClassNumber

	if len(request.Tasks) > 0 {
		if name == "" {
			name = request.Tasks[0].Name
		}
		if code == 0 {
			code = request.Tasks[0].Code
		}
		if taskType == "" {
			taskType = request.Tasks[0].Type
		}
	}

	if len(request.ClassNumbers) > 0 && classNumber == "" {
		classNumber = request.ClassNumbers[0]
	}

	if name != "" {
		task.Name = name
	}
	if code != 0 {
		task.Code = code
	}
	if taskType != "" {
		task.Type = taskType
	}
	if classNumber != "" {
		task.ClassNumber = classNumber
	}
	if request.SchoolNumber != "" {
		task.SchoolNumber = request.SchoolNumber
	}
	if request.CampusNumber != "" {
		task.CampusNumber = request.CampusNumber
	}
	if request.CampusCurriculumNumber != "" {
		task.CampusCurriculumNumber = request.CampusCurriculumNumber
	}
	if request.LevelNumber != "" {
		task.LevelNumber = request.LevelNumber
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

func (s *timetableTaskService) GetAllTimetableTasks(classNumber string) ([]responses.TimetableTaskResponse, error) {
	tasks, err := s.repo.FindAllTimetableTasks(classNumber)
	if err != nil {
		return nil, err
	}

	var response []responses.TimetableTaskResponse
	for _, task := range tasks {
		res := responses.TimetableTaskResponse{
			Number:           task.Number,
			Name:             task.Name,
			Type:             task.Type,
			Code:             task.Code,
			UserNumber:       task.UserNumber,
			ClassNumber:      task.ClassNumber,
			StreamNumber:     task.StreamNumber,
			SchoolNumber:     task.SchoolNumber,
			CampusNumber:     task.CampusNumber,
			CurriculumNumber: task.CampusCurriculumNumber,
			LevelNumber:      task.LevelNumber,
		}

		// Fallback derivation for legacy tasks where DB fields are empty
		if res.SchoolNumber == "" || res.CampusNumber == "" || res.CurriculumNumber == "" || res.LevelNumber == "" {
			if task.ClassNumber != "" {
				campusClass, err := repositories.GetCampusClassByNumber(task.ClassNumber)
				if err == nil && campusClass.Number != "" {
					if res.CampusNumber == "" {
						res.CampusNumber = campusClass.CampusNumber
					}

					campus, err := campusRepositories.GetCampusByNumber(campusClass.CampusNumber)
					if err == nil && campus.Number != "" {
						if res.SchoolNumber == "" {
							res.SchoolNumber = campus.SchoolNumber
						}
					}

					campusLevel, err := levelRepositories.GetCampusLevelByNumber(campusClass.CampusLevelNumber)
					if err == nil && campusLevel.Number != "" {
						curriculumLevel, err := levelRepositories.GetCurriculumLevelByNumber(campusLevel.CurriculumLevelNumber)
						if err == nil && curriculumLevel.Number != "" {
							if res.CurriculumNumber == "" {
								res.CurriculumNumber = curriculumLevel.CurriculumNumber
							}
							if res.LevelNumber == "" {
								res.LevelNumber = curriculumLevel.LevelNumber
							}
						}
					}
				}
			}
		}

		// Query readable names/labels
		if res.ClassNumber != "" {
			campusClass, err := repositories.GetCampusClassByNumber(res.ClassNumber)
			if err == nil && campusClass.Number != "" {
				res.Classes = campusClass.ClassName
				res.Class = campusClass.ClassName
			}
		}

		if res.CampusNumber != "" {
			campus, err := campusRepositories.GetCampusByNumber(res.CampusNumber)
			if err == nil && campus.Number != "" {
				res.Campus = campus.Name
			}
		}

		if res.SchoolNumber != "" {
			school, err := schoolRepositories.GetSchoolByNumber(res.SchoolNumber)
			if err == nil && school.Number != "" {
				res.School = school.Name
			}
		}

		if res.LevelNumber != "" {
			level, err := levelRepositories.GetLevelByNumber(res.LevelNumber)
			if err == nil && level.Number != "" {
				res.Level = level.Name
			}
		}

		response = append(response, res)
	}

	return response, nil
}
