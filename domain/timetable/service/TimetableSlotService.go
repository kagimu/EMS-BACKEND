package service

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type TimetableSlotService interface {
	CreateTimetableSlot(request requests.TimetableSlotRequest) error
	UpdateTimetableSlot(slotNumber string, request requests.UpdateTimetableSlotRequest) error
	DeleteTimetableSlot(slotNumber string) error
	GetTimetableSlotByNumber(slotNumber string) (models.TimetableSlots, error)
	GetAllTimetableSlots(schoolNumber, levelNumber string) ([]models.TimetableSlots, error)
}

type timetableSlotService struct {
	repo repository.TimetableSlotRepository
}

func NewTimetableSlotService(repo repository.TimetableSlotRepository) TimetableSlotService {
	return &timetableSlotService{repo: repo}
}

func (s *timetableSlotService) CreateTimetableSlot(request requests.TimetableSlotRequest) error {
	if len(request.Slots) == 0 {
		return errors.New("at least one slot is required")
	}

	for _, slotItem := range request.Slots {
		if slotItem.StartTime == "" || slotItem.EndTime == "" {
			return errors.New("start time and end time cannot be empty")
		}

		slot := models.TimetableSlots{
			Number:       uuid.NewV4().String(),
			Name:         request.Name,
			Type:         request.Type,
			SchoolNumber: request.SchoolNumber,
			LevelNumber:  request.LevelNumber,
			StartTime:    slotItem.StartTime,
			EndTime:      slotItem.EndTime,
		}

		if err := s.repo.SaveTimetableSlot(slot); err != nil {
			return err
		}
	}

	return nil
}

func (s *timetableSlotService) UpdateTimetableSlot(slotNumber string, request requests.UpdateTimetableSlotRequest) error {
	slot, err := s.repo.FindTimetableSlotByNumber(slotNumber)
	if err != nil {
		return errors.New("error finding timetable slot: " + err.Error())
	}

	if slot.Number == "" {
		return errors.New("timetable slot not found")
	}

	if request.Name != "" {
		slot.Name = request.Name
	}
	if request.Type != "" {
		slot.Type = request.Type
	}
	if request.SchoolNumber != "" {
		slot.SchoolNumber = request.SchoolNumber
	}
	if request.LevelNumber != "" {
		slot.LevelNumber = request.LevelNumber
	}
	if request.StartTime != "" {
		slot.StartTime = request.StartTime
	}
	if request.EndTime != "" {
		slot.EndTime = request.EndTime
	}

	return s.repo.UpdateTimetableSlot(&slot)
}

func (s *timetableSlotService) DeleteTimetableSlot(slotNumber string) error {
	slot, err := s.repo.FindTimetableSlotByNumber(slotNumber)
	if err != nil {
		return errors.New("error finding timetable slot: " + err.Error())
	}

	if slot.Number == "" {
		return errors.New("timetable slot not found")
	}

	return s.repo.DeleteTimetableSlot(slot)
}

func (s *timetableSlotService) GetTimetableSlotByNumber(slotNumber string) (models.TimetableSlots, error) {
	slot, err := s.repo.FindTimetableSlotByNumber(slotNumber)
	if err != nil {
		return models.TimetableSlots{}, errors.New("error finding timetable slot: " + err.Error())
	}

	if slot.Number == "" {
		return models.TimetableSlots{}, errors.New("timetable slot not found")
	}

	return slot, nil
}

func (s *timetableSlotService) GetAllTimetableSlots(schoolNumber, levelNumber string) ([]models.TimetableSlots, error) {
	return s.repo.FindAllTimetableSlots(schoolNumber, levelNumber)
}
