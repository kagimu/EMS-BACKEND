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
	ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber string) ([]models.StudyTimetablePeriod, error)
	GetStructuredTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber string) (responses.TimetablePeriodsResponse, error)
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

func (s *studyTimetablePeriodService) GetStructuredTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber string) (responses.TimetablePeriodsResponse, error) {
	// Get all periods
	periods, err := s.periodRepo.FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber)
	if err != nil {
		return responses.TimetablePeriodsResponse{}, err
	}

	if len(periods) == 0 {
		return responses.TimetablePeriodsResponse{
			StudyTimetableNumber: studyTimetableNumber,
			Days:                 []responses.DayPeriods{},
		}, nil
	}

	// Build a map of task types to tasks for quick lookup
	taskTypeMap := make(map[string]models.TimetableTask)
	taskTypes := make(map[string]bool)
	for _, period := range periods {
		if !taskTypes[period.TaskType] {
			taskTypes[period.TaskType] = true
		}
	}

	// Fetch all tasks that match the task types we found
	// Note: This assumes we can get tasks by type, which might need adjustment
	// For now, we'll query all tasks and filter by type
	allTasks, err := s.taskRepo.FindAllTimetableTasks("")
	if err == nil {
		for _, task := range allTasks {
			if taskTypes[task.Type] {
				// Use the first task found for this type
				if _, exists := taskTypeMap[task.Type]; !exists {
					taskTypeMap[task.Type] = task
				}
			}
		}
	}

	// Build a map of slots by slot number (we'll need to match by time if we don't have slot number stored)
	// For now, we'll format the time slot directly from period times

	// Group periods by day
	dayMap := make(map[string][]models.StudyTimetablePeriod)
	for _, period := range periods {
		dayKey := fmt.Sprintf("%s_%d", period.Day, period.DayCode)
		dayMap[dayKey] = append(dayMap[dayKey], period)
	}

	// Convert to response structure
	var days []responses.DayPeriods
	dayOrder := []string{"monday_1", "tuesday_2", "wednesday_3", "thursday_4", "friday_5", "saturday_6", "sunday_7"}

	for _, dayKey := range dayOrder {
		periodsForDay, exists := dayMap[dayKey]
		if !exists {
			continue
		}

		// Sort periods by start time
		sort.Slice(periodsForDay, func(i, j int) bool {
			return periodsForDay[i].StartTime < periodsForDay[j].StartTime
		})

		var periodItems []responses.StudyTimetablePeriodItem
		for _, period := range periodsForDay {
			task, taskFound := taskTypeMap[period.TaskType]
			taskName := period.TaskType
			taskNumber := ""
			if taskFound {
				taskName = task.Name
				taskNumber = task.Number
			}

			// Format time slot (e.g., "8:00 AM - 9:00 AM")
			timeSlot := formatTimeSlot(period.StartTime, period.EndTime)

			periodItem := responses.StudyTimetablePeriodItem{
				PeriodNumber: period.Number,
				TaskName:     taskName,
				TaskType:     period.TaskType,
				TaskNumber:   taskNumber,
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
				Day:     periodsForDay[0].Day,
				DayCode: periodsForDay[0].DayCode,
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
