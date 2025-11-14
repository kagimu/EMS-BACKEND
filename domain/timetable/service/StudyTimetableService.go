package service

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
)

type StudyTimetableService interface {
	GetStudyTimetableByNumber(timetableNumber string) (models.StudyTimetable, error)
	GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber string) ([]models.StudyTimetable, error)
}

type studyTimetableService struct {
	repo repository.StudyTimetableRepository
}

func NewStudyTimetableService(repo repository.StudyTimetableRepository) StudyTimetableService {
	return &studyTimetableService{repo: repo}
}

func (s *studyTimetableService) GetStudyTimetableByNumber(timetableNumber string) (models.StudyTimetable, error) {
	if timetableNumber == "" {
		return models.StudyTimetable{}, errors.New("timetable number is required")
	}

	timetable, err := s.repo.FindStudyTimetableByNumber(timetableNumber)
	if err != nil {
		return models.StudyTimetable{}, errors.New("error finding study timetable: " + err.Error())
	}

	if timetable.Number == "" {
		return models.StudyTimetable{}, errors.New("study timetable not found")
	}

	return timetable, nil
}

func (s *studyTimetableService) GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber string) ([]models.StudyTimetable, error) {
	return s.repo.FindAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber)
}
