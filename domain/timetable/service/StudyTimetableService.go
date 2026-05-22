package service

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
)

type StudyTimetableService interface {
	GetStudyTimetableByNumber(timetableNumber string) (models.StudyTimetable, error)
	GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, campusCurriculumNumber, academicYearNumber, campusClassNumber, classNumber, streamNumber string) ([]responses.StudyTimetables, error)
	DeleteStudyTimetable(timetableNumber string) error
}

type studyTimetableService struct {
	repo       repository.StudyTimetableRepository
	periodRepo repository.StudyTimetablePeriodRepository
}

func NewStudyTimetableService(repo repository.StudyTimetableRepository, periodRepo repository.StudyTimetablePeriodRepository) StudyTimetableService {
	return &studyTimetableService{repo: repo, periodRepo: periodRepo}
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

func (s *studyTimetableService) GetAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, campusCurriculumNumber, academicYearNumber, campusClassNumber, classNumber, streamNumber string) ([]responses.StudyTimetables, error) {
	var err error
	timetables, err := s.repo.FindAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, campusCurriculumNumber, academicYearNumber, campusClassNumber, classNumber, streamNumber)
	if err != nil {
		return nil, err
	}

	var response []responses.StudyTimetables

	for _, timetable := range timetables {
		var schoolName string
		school, err := repositories.GetSchoolByNumber(timetable.SchoolNumber)

		if err == nil {
			schoolName = school.Name
		}

		var campusName string
		campus, err := repositories2.GetCampusByNumber(timetable.CampusNumber)

		if err == nil {
			campusName = campus.Name
		}

		// Resolve class_number via campus_class → curriculum_class → class
		var className, classNumber string
		campusClass, err := repositories3.GetCampusClassByNumber(timetable.CampusClassNumber)
		if err == nil {
			className = campusClass.ClassName
			classNumber = campusClass.ClassNumber
		}

		// Resolve curriculum_number via campus_curriculum → school_curriculum
		var curriculumNumber string
		if timetable.CampusCurriculumNumber != "" {
			campusCurriculum, err := repositories2.GetCampusCurriculumByNumber(timetable.CampusCurriculumNumber)
			if err == nil && campusCurriculum.SchoolCurriculumNumber != "" {
				schoolCurriculum, err := repositories.GetSchoolCurriculumByNumber(campusCurriculum.SchoolCurriculumNumber)
				if err == nil {
					curriculumNumber = schoolCurriculum.CurriculumNumber
				}
			}
		}

		var streamName string
		stream, err := repositories3.GetStreamByNumber(timetable.StreamNumber)
		if err == nil {
			streamName = stream.Name
		}

		studyPeriod, err := repositories4.GetStudyPeriodByNumber(timetable.StudyPeriodNumber)

		if err != nil {
			return nil, err
		}

		studyPeriodName, err := repositories4.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

		if err != nil {
			return nil, err
		}

		academicYear, err := repositories5.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)

		if err != nil {
			return nil, err
		}

		response = append(response, responses.StudyTimetables{
			Name:                   timetable.Name,
			Number:                 timetable.Number,
			Status:                 timetable.Status,
			School:                 schoolName,
			Campus:                 campusName,
			CLass:                  className,
			Stream:                 streamName,
			StudyPeriod:            studyPeriodName.Name,
			AcademicYear:           academicYear.AcademicYearName,
			// Number fields
			SchoolNumber:           timetable.SchoolNumber,
			CampusNumber:           timetable.CampusNumber,
			LevelNumber:            timetable.LevelNumber,
			ClassNumber:            classNumber,
			StreamNumber:           timetable.StreamNumber,
			AcademicYearNumber:     timetable.AcademicYearNumber,
			CampusCurriculumNumber: timetable.CampusCurriculumNumber,
			CurriculumNumber:       curriculumNumber,
			StudyPeriodNumber:      timetable.StudyPeriodNumber,
		})
	}

	return response, nil
}

func (s *studyTimetableService) DeleteStudyTimetable(timetableNumber string) error {
	if timetableNumber == "" {
		return errors.New("timetable number is required")
	}

	timetable, err := s.repo.FindStudyTimetableByNumber(timetableNumber)
	if err != nil {
		return errors.New("error finding study timetable: " + err.Error())
	}

	if timetable.Number == "" {
		return errors.New("study timetable not found")
	}

	// Delete child periods and parent timetable in a single transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := s.periodRepo.DeletePeriodsByStudyTimetableNumber(timetableNumber, tx); err != nil {
		tx.Rollback()
		return errors.New("error deleting timetable periods: " + err.Error())
	}

	if err := s.repo.DeleteStudyTimetable(timetable, tx); err != nil {
		tx.Rollback()
		return errors.New("error deleting study timetable: " + err.Error())
	}

	return tx.Commit().Error
}
