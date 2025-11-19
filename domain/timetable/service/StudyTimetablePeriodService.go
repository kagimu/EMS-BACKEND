package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type StudyTimetablePeriodService interface {
	CreateStudyTimetablePeriod(request requests.CreateStudyTimetablePeriodRequest, userNumber string) error
	UpdateStudyTimetablePeriod(periodNumber string, request requests.UpdateStudyTimetablePeriodRequest) error
	DeleteStudyTimetablePeriod(periodNumber string) error
	GetStudyTimetablePeriodByNumber(periodNumber string) (models.StudyTimetablePeriod, error)
	ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day string) ([]models.StudyTimetablePeriod, error)
	GetStructuredTimetablePeriods(studyTimetableNumber, streamNumber, day string) (responses.TimetablePeriodsResponse, error)
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

func (s *studyTimetablePeriodService) ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day string) ([]models.StudyTimetablePeriod, error) {
	return s.periodRepo.FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, day)
}

func (s *studyTimetablePeriodService) GetStructuredTimetablePeriods(studyTimetableNumber, streamNumber, day string) (responses.TimetablePeriodsResponse, error) {
	// Get all periods

	// Convert to response structure
	var days []responses.DayPeriods
	dayOrder := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}

	for _, dayKey := range dayOrder {
		periodsForDay, err := s.periodRepo.FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, dayKey)
		if err != nil {
			return responses.TimetablePeriodsResponse{}, err
		}

		// Sort periods by start time
		sort.Slice(periodsForDay, func(i, j int) bool {
			return periodsForDay[i].StartTime < periodsForDay[j].StartTime
		})

		var periodItems = make([]responses.StudyTimetablePeriodItem, 0)
		for _, period := range periodsForDay {

			// Format time slot (e.g., "8:00 AM - 9:00 AM")
			timeSlot := formatTimeSlot(period.StartTime, period.EndTime)

			periodItem := responses.StudyTimetablePeriodItem{
				PeriodNumber: period.Number,
				TaskName:     period.Task,
				TaskType:     period.TaskType,
				TaskNumber:   period.Number,
				StartTime:    period.StartTime,
				EndTime:      period.EndTime,
				TimeSlot:     timeSlot,
				SlotNumber:   "", // We don't store slot number in period, would need model update
				StreamNumber: period.StreamNumber,
				UserNumber:   period.UserNumber,
			}
			periodItems = append(periodItems, periodItem)
		}

		if len(periodItems) > 0 {
			dayPeriods := responses.DayPeriods{
				Day:     dayKey,
				DayCode: mapDayToCode(dayKey),
				Periods: periodItems,
			}
			days = append(days, dayPeriods)
		}
	}

	return responses.TimetablePeriodsResponse{
		StudyTimetableNumber: studyTimetableNumber,
		Days:                 days,
	}, nil
}

func formatTimeSlot(startTime, endTime string) string {
	return fmt.Sprintf("%s - %s", startTime, endTime)
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
