package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"gorm.io/gorm"
)

type TimetableSlotRepository interface {
	SaveTimetableSlot(slot models.TimetableSlots) error
	UpdateTimetableSlot(slot *models.TimetableSlots) error
	DeleteTimetableSlot(slot models.TimetableSlots) error
	FindTimetableSlotByNumber(number string) (models.TimetableSlots, error)
	FindTimetableSlotByID(id uint) (models.TimetableSlots, error)
	FindAllTimetableSlots(schoolNumber, levelNumber string) ([]models.TimetableSlots, error)
}

type timetableSlotRepository struct {
	DB *gorm.DB
}

func NewTimetableSlotRepository(db *gorm.DB) TimetableSlotRepository {
	return &timetableSlotRepository{DB: db}
}

func (r *timetableSlotRepository) SaveTimetableSlot(slot models.TimetableSlots) error {
	return r.DB.Save(&slot).Error
}

func (r *timetableSlotRepository) UpdateTimetableSlot(slot *models.TimetableSlots) error {
	return r.DB.Save(slot).Error
}

func (r *timetableSlotRepository) DeleteTimetableSlot(slot models.TimetableSlots) error {
	return r.DB.Unscoped().Delete(&slot).Error
}

func (r *timetableSlotRepository) FindTimetableSlotByNumber(number string) (models.TimetableSlots, error) {
	var slot models.TimetableSlots

	if err := r.DB.Where("number = ?", number).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TimetableSlots{}, nil
		}
		return models.TimetableSlots{}, err
	}

	return slot, nil
}

func (r *timetableSlotRepository) FindTimetableSlotByID(id uint) (models.TimetableSlots, error) {
	var slot models.TimetableSlots

	if err := r.DB.Where("id = ?", id).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TimetableSlots{}, nil
		}
		return models.TimetableSlots{}, err
	}

	return slot, nil
}

func (r *timetableSlotRepository) FindAllTimetableSlots(schoolNumber, levelNumber string) ([]models.TimetableSlots, error) {
	var slots []models.TimetableSlots

	query := r.DB
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	if levelNumber != "" {
		query = query.Where("level_number = ?", levelNumber)
	}

	err := query.Find(&slots).Error
	if err != nil {
		return nil, err
	}

	return slots, nil
}
