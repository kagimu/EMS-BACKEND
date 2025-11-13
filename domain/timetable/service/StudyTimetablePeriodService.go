package service

import (
	"errors"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type StudyTimetablePeriodService interface {
	CreateStudyTimetablePeriod(request requests.CreateStudyTimetablePeriodRequest, userNumber string) error
	UpdateStudyTimetablePeriod(periodNumber string, request requests.UpdateStudyTimetablePeriodRequest) error
	DeleteStudyTimetablePeriod(periodNumber string) error
	GetStudyTimetablePeriodByNumber(periodNumber string) (models.StudyTimetablePeriod, error)
	ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber string) ([]models.StudyTimetablePeriod, error)
}

type studyTimetablePeriodService struct {
	periodRepo repository.StudyTimetablePeriodRepository
	slotRepo   repository.TimetableSlotRepository
	taskRepo   repository.TimetableTaskRepository
}

func NewStudyTimetablePeriodService(
	periodRepo repository.StudyTimetablePeriodRepository,
	slotRepo repository.TimetableSlotRepository,
	taskRepo repository.TimetableTaskRepository,
) StudyTimetablePeriodService {
	return &studyTimetablePeriodService{
		periodRepo: periodRepo,
		slotRepo:   slotRepo,
		taskRepo:   taskRepo,
	}
}

func (s *studyTimetablePeriodService) CreateStudyTimetablePeriod(request requests.CreateStudyTimetablePeriodRequest, userNumber string) error {
	if userNumber == "" {
		return errors.New("user number is required")
	}

	if len(request.Periods) == 0 {
		return errors.New("at least one study timetable period is required")
	}

	for _, periodItem := range request.Periods {
		if periodItem.Day == "" {
			return errors.New("period day is required")
		}
		if periodItem.SlotNumber == "" {
			return errors.New("period slot number is required")
		}
		if periodItem.TaskNumber == "" {
			return errors.New("period task number is required")
		}

		slot, err := s.slotRepo.FindTimetableSlotByNumber(periodItem.SlotNumber)
		if err != nil {
			return errors.New("error finding timetable slot: " + err.Error())
		}
		if slot.Number == "" {
			return errors.New("timetable slot not found")
		}

		task, err := s.taskRepo.FindTimetableTaskByNumber(periodItem.TaskNumber)
		if err != nil {
			return errors.New("error finding timetable task: " + err.Error())
		}
		if task.Number == "" {
			return errors.New("timetable task not found")
		}

		period := models.StudyTimetablePeriod{
			Number:               uuid.NewV4().String(),
			StudyTimetableNumber: request.StudyTimetableNumber,
			StreamNumber:         request.StreamNumber,
			Day:                  periodItem.Day,
			DayCode:              mapDayToCode(periodItem.Day),
			StartTime:            slot.StartTime,
			EndTime:              slot.EndTime,
			TaskType:             task.Type,
			UserNumber:           userNumber,
		}

		if err := s.periodRepo.SaveStudyTimetablePeriod(period); err != nil {
			return err
		}
	}

	return nil
}

func (s *studyTimetablePeriodService) UpdateStudyTimetablePeriod(periodNumber string, request requests.UpdateStudyTimetablePeriodRequest) error {
	if periodNumber == "" {
		return errors.New("period number is required")
	}

	period, err := s.periodRepo.FindStudyTimetablePeriodByNumber(periodNumber)
	if err != nil {
		return errors.New("error finding study timetable period: " + err.Error())
	}

	if period.Number == "" {
		return errors.New("study timetable period not found")
	}

	if request.StudyTimetableNumber != "" {
		period.StudyTimetableNumber = request.StudyTimetableNumber
	}
	if request.StreamNumber != "" {
		period.StreamNumber = request.StreamNumber
	}
	if request.Day != "" {
		period.Day = request.Day
		period.DayCode = mapDayToCode(request.Day)
	}
	if request.DayCode != nil {
		period.DayCode = *request.DayCode
	}
	if request.StartTime != "" {
		period.StartTime = request.StartTime
	}
	if request.EndTime != "" {
		period.EndTime = request.EndTime
	}
	if request.TaskType != "" {
		period.TaskType = request.TaskType
	}

	return s.periodRepo.UpdateStudyTimetablePeriod(&period)
}

func (s *studyTimetablePeriodService) DeleteStudyTimetablePeriod(periodNumber string) error {
	if periodNumber == "" {
		return errors.New("period number is required")
	}

	period, err := s.periodRepo.FindStudyTimetablePeriodByNumber(periodNumber)
	if err != nil {
		return errors.New("error finding study timetable period: " + err.Error())
	}

	if period.Number == "" {
		return errors.New("study timetable period not found")
	}

	return s.periodRepo.DeleteStudyTimetablePeriod(period)
}

func (s *studyTimetablePeriodService) GetStudyTimetablePeriodByNumber(periodNumber string) (models.StudyTimetablePeriod, error) {
	if periodNumber == "" {
		return models.StudyTimetablePeriod{}, errors.New("period number is required")
	}

	period, err := s.periodRepo.FindStudyTimetablePeriodByNumber(periodNumber)
	if err != nil {
		return models.StudyTimetablePeriod{}, errors.New("error finding study timetable period: " + err.Error())
	}

	if period.Number == "" {
		return models.StudyTimetablePeriod{}, errors.New("study timetable period not found")
	}

	return period, nil
}

func (s *studyTimetablePeriodService) ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber string) ([]models.StudyTimetablePeriod, error) {
	return s.periodRepo.FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber)
}

func mapDayToCode(day string) int32 {
	switch strings.ToLower(day) {
	case "monday":
		return 1
	case "tuesday":
		return 2
	case "wednesday":
		return 3
	case "thursday":
		return 4
	case "friday":
		return 5
	case "saturday":
		return 6
	case "sunday":
		return 7
	default:
		return 0
	}
}
