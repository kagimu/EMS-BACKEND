package service

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
)

type StudyTimetableService interface {
	GetStudyTimetableByNumber(timetableNumber string) (models.StudyTimetable, error)
	GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber string) ([]responses.StudyTimetables, error)
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

func (s *studyTimetableService) GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber string) ([]responses.StudyTimetables, error) {
	timetables, err := s.repo.FindAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, classNumber, streamNumber)
	if err != nil {
		return nil, err
	}

	var response []responses.StudyTimetables
	for _, timetable := range timetables {
		response = append(response, responses.StudyTimetables{
			Name:        timetable.Name,
			Number:      timetable.Number,
			Status:      timetable.Status,
			School:      timetable.SchoolNumber,
			Campus:      timetable.CampusNumber,
			CLass:       timetable.ClassNumber,
			Stream:      timetable.StreamNumber,
			StudyPeriod: timetable.StudyPeriodNumber,
			AcademicYear: func() string {
				if timetable.StudyPeriod.Number != "" {
					return timetable.StudyPeriod.AcademicYearNumber
				}
				return ""
			}(),
		})
	}

	return response, nil
}
