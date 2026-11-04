package service

import (
	"errors"

	campusRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	curriculumModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	schoolRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	levelRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	uuid "github.com/satori/go.uuid"
)

type TimetableSlotService interface {
	CreateTimetableSlot(request requests.TimetableSlotRequest) error
	UpdateTimetableSlot(slotNumber string, request requests.UpdateTimetableSlotRequest) error
	DeleteTimetableSlot(slotNumber string) error
	GetTimetableSlotByNumber(slotNumber string) (models.TimetableSlots, error)
	GetAllTimetableSlots(schoolNumber, levelNumber, streamNumber string) ([]responses.TimetableSlotResponse, error)
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
			Number:                 uuid.NewV4().String(),
			Name:                   request.Name,
			Type:                   request.Type,
			SchoolNumber:           request.SchoolNumber,
			CampusNumber:           request.CampusNumber,
			CampusCurriculumNumber: request.CampusCurriculumNumber,
			LevelNumber:            request.LevelNumber,
			AcademicYearNumber:     request.AcademicYearNumber,
			ClassNumber:            request.ClassNumber,
			StreamNumber:           request.StreamNumber,
			StartTime:              slotItem.StartTime,
			EndTime:                slotItem.EndTime,
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
	if request.CampusNumber != "" {
		slot.CampusNumber = request.CampusNumber
	}
	if request.CampusCurriculumNumber != "" {
		slot.CampusCurriculumNumber = request.CampusCurriculumNumber
	}
	if request.LevelNumber != "" {
		slot.LevelNumber = request.LevelNumber
	}
	if request.AcademicYearNumber != "" {
		slot.AcademicYearNumber = request.AcademicYearNumber
	}
	if request.ClassNumber != "" {
		slot.ClassNumber = request.ClassNumber
	}
	if request.StreamNumber != "" {
		slot.StreamNumber = request.StreamNumber
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

func (s *timetableSlotService) GetAllTimetableSlots(schoolNumber, levelNumber, streamNumber string) ([]responses.TimetableSlotResponse, error) {
	// If streamNumber is provided but levelNumber/schoolNumber are empty, let's resolve them
	if streamNumber != "" && (schoolNumber == "" || levelNumber == "") {
		stream, err := classRepositories.GetStreamByNumber(streamNumber)
		if err == nil && stream.Number != "" {
			campusClass, err := classRepositories.GetCampusClassByNumber(stream.CampusClassNumber)
			if err == nil && campusClass.Number != "" {
				if levelNumber == "" {
					levelNumber = campusClass.CampusLevelNumber
				}
				if schoolNumber == "" {
					schoolNumber = campusClass.Campus.SchoolNumber
				}
			}
		}
	}

	// Try querying with all parameters, including streamNumber
	slots, err := s.repo.FindAllTimetableSlots(schoolNumber, levelNumber, streamNumber)
	if err != nil {
		return nil, err
	}

	// If no slots found and streamNumber was provided, let's fallback to query by levelNumber and schoolNumber to fetch legacy slots
	if len(slots) == 0 && streamNumber != "" {
		slots, err = s.repo.FindAllTimetableSlots(schoolNumber, levelNumber, "")
		if err != nil {
			return nil, err
		}
	}

	var response []responses.TimetableSlotResponse
	for _, slot := range slots {
		res := responses.TimetableSlotResponse{
			ID:                     slot.ID,
			Number:                 slot.Number,
			Name:                   slot.Name,
			SchoolNumber:           slot.SchoolNumber,
			CampusNumber:           slot.CampusNumber,
			CampusCurriculumNumber: slot.CampusCurriculumNumber,
			LevelNumber:            slot.LevelNumber,
			AcademicYearNumber:     slot.AcademicYearNumber,
			ClassNumber:            slot.ClassNumber,
			StreamNumber:           slot.StreamNumber,
			Type:                   slot.Type,
			StartTime:              slot.StartTime,
			EndTime:                slot.EndTime,
		}

		// Fill in stream details if this was a legacy slot but we queried it under a specific stream
		if res.StreamNumber == "" && streamNumber != "" {
			res.StreamNumber = streamNumber
		}

		// Fallback derivation for legacy slots where DB fields are empty
		if res.SchoolNumber == "" || res.CampusNumber == "" || res.LevelNumber == "" {
			// Derivation if stream is available
			if res.StreamNumber != "" {
				stream, err := classRepositories.GetStreamByNumber(res.StreamNumber)
				if err == nil && stream.Number != "" {
					res.Stream = stream.Name
					if res.ClassNumber == "" {
						res.ClassNumber = stream.CampusClassNumber
					}
				}
			}

			// Derivation if class is available
			if res.ClassNumber != "" {
				campusClass, err := classRepositories.GetCampusClassByNumber(res.ClassNumber)
				if err == nil && campusClass.Number != "" {
					res.Class = campusClass.ClassName
					if res.CampusNumber == "" {
						res.CampusNumber = campusClass.CampusNumber
					}
					if res.LevelNumber == "" {
						res.LevelNumber = campusClass.CampusLevelNumber
					}
				}
			}

			// Derivation if campus is available
			if res.CampusNumber != "" {
				campus, err := campusRepositories.GetCampusByNumber(res.CampusNumber)
				if err == nil && campus.Number != "" {
					res.Campus = campus.Name
					if res.SchoolNumber == "" {
						res.SchoolNumber = campus.SchoolNumber
					}
				}
			}
		}

		// Query names/labels
		if res.SchoolNumber != "" {
			school, err := schoolRepositories.GetSchoolByNumber(res.SchoolNumber)
			if err == nil && school.Number != "" {
				res.School = school.Name
			}
		}

		if res.CampusNumber != "" {
			campus, err := campusRepositories.GetCampusByNumber(res.CampusNumber)
			if err == nil && campus.Number != "" {
				res.Campus = campus.Name
			}
		}

		if res.LevelNumber != "" {
			// Check if LevelNumber is a settings Level or a CampusLevel
			level, err := levelRepositories.GetLevelByNumber(res.LevelNumber)
			if err == nil && level.Number != "" {
				res.Level = level.Name
			} else {
				campusLevel, err := levelRepositories.GetCampusLevelByNumber(res.LevelNumber)
				if err == nil && campusLevel.Number != "" {
					curriculumLevel, err := levelRepositories.GetCurriculumLevelByNumber(campusLevel.CurriculumLevelNumber)
					if err == nil && curriculumLevel.Number != "" {
						if res.CampusCurriculumNumber == "" {
							res.CampusCurriculumNumber = curriculumLevel.CurriculumNumber
						}
						lvl, err := levelRepositories.GetLevelByNumber(curriculumLevel.LevelNumber)
						if err == nil && lvl.Number != "" {
							res.Level = lvl.Name
						}
					}
				}
			}
		}

		if res.CampusCurriculumNumber != "" {
			curriculum, err := curriculumModels.GetCurriculumByNumber(res.CampusCurriculumNumber)
			if err == nil && curriculum.Number != "" {
				res.Curriculum = curriculum.Name
			}
		}

		if res.ClassNumber != "" {
			campusClass, err := classRepositories.GetCampusClassByNumber(res.ClassNumber)
			if err == nil && campusClass.Number != "" {
				res.Class = campusClass.ClassName
			}
		}

		if res.StreamNumber != "" {
			stream, err := classRepositories.GetStreamByNumber(res.StreamNumber)
			if err == nil && stream.Number != "" {
				res.Stream = stream.Name
			}
		}

		response = append(response, res)
	}

	return response, nil
}
