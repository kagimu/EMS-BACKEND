package service

import (
	"errors"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
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
	periodRepo    repository.StudyTimetablePeriodRepository
	slotRepo      repository.TimetableSlotRepository
	taskRepo      repository.TimetableTaskRepository
	timetableRepo repository.StudyTimetableRepository
}

func NewStudyTimetablePeriodService(
	periodRepo repository.StudyTimetablePeriodRepository,
	slotRepo repository.TimetableSlotRepository,
	taskRepo repository.TimetableTaskRepository,
	timetableRepo repository.StudyTimetableRepository,
) StudyTimetablePeriodService {
	return &studyTimetablePeriodService{
		periodRepo:    periodRepo,
		slotRepo:      slotRepo,
		taskRepo:      taskRepo,
		timetableRepo: timetableRepo,
	}
}

func (s *studyTimetablePeriodService) CreateStudyTimetablePeriod(request requests.CreateStudyTimetablePeriodRequest, userNumber string) error {
	if userNumber == "" {
		return errors.New("user number is required")
	}

	if len(request.Periods) == 0 {
		return errors.New("at least one study timetable period is required")
	}

	if request.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	if request.StreamNumber == "" {
		return errors.New("stream number is required")
	}

	studyTimetable := models.StudyTimetable{
		Number:            uuid.NewV4().String(),
		Name:              "StudyTimetablePeriod",
		Status:            constants.StatusDraft,
		UserNumber:        userNumber,
		StreamNumber:      request.StreamNumber,
		StudyPeriodNumber: request.StudyPeriodNumber,
	}

	if err := s.timetableRepo.SaveStudyTimetable(studyTimetable); err != nil {
		return errors.New("error saving study timetable: " + err.Error())
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
			StudyTimetableNumber: studyTimetable.Number,
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
