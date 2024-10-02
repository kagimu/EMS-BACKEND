package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"
	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
)

type LessonPlanService interface {
	ImportPlan(lessonPlanFile multipart.File, subTopicNumber, currentUserNumber string) error
}

type lessonPlanServiceImp struct {
	Db *gorm.DB
}

var lessonPlanRepository = repository.NewLessonPlanRepository()
var conceptRepository = repositories.NewConceptRepository()

func (l *lessonPlanServiceImp) ImportPlan(lessonPlanFile multipart.File, schemeSubTopicNumber, currentUserNumber string) error {
	excelFile, err := excelize.OpenReader(lessonPlanFile)
	if err != nil {
		return err
	}

	schemeOfWorkSubTopic, err := repositories2.GetSchemeOfWorkSubTopicByNumber(schemeSubTopicNumber)
	if err != nil {
		return errors.New("Checking scheme subtopic error: " + err.Error())
	}

	if schemeOfWorkSubTopic.Number == "" {
		return errors.New("Invalid scheme subtopic ")
	}

	conceptRepository.UseTx(config.DB)
	concepts, err := conceptRepository.FindConceptsBySubTopicNumber(schemeOfWorkSubTopic.SubTopicNumber)
	if err != nil {
		return errors.New("Getting concepts error " + err.Error())
	}

	if len(concepts) == 0 {
		return errors.New("The selected sub-topic has no concepts ")
	}

	lessonPlanRepository.UseTx(config.DB)
	lessonPlan, err := lessonPlanRepository.GetLessonPlanBySchemeOfWorkSubTopicNumber(schemeSubTopicNumber)
	if err != nil {
		return errors.New("Checking lesson plan error " + err.Error())
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		lessonPlanRepository.UseTx(tx)

		if lessonPlan.Number == "" {
			lessonPlan = models.LessonPlan{
				Status:               constants.StatusDraft,
				Number:               uuid.NewV4().String(),
				UserNumber:           currentUserNumber,
				SchemeSubTopicNumber: schemeSubTopicNumber,
			}
			err = lessonPlanRepository.SaveLessonPlan(lessonPlan)
			if err != nil {
				return errors.New("Creating lesson plan error " + err.Error())
			}
		}

		for _, concept := range concepts {
			rows, err := excelFile.GetRows(concept.Name)
			if err != nil {
				return errors.New(fmt.Sprintf("Error %s extracting for %s", err.Error(), concept.Name))
			}

			lessonPlanConcept, err := lessonPlanRepository.GetLessonPlanConcept(lessonPlan.Number, concept.Number)
			if err != nil {
				return errors.New("Checking lesson plan concept error " + err.Error())
			}

			if lessonPlanConcept.Number == "" {
				lessonPlanConcept = models.LessonPlanConcept{
					Number:           uuid.NewV4().String(),
					LessonPlanNumber: lessonPlan.Number,
					ConceptNumber:    concept.Number,
					Code:             concept.Code,
					UserNumber:       currentUserNumber,
				}
				err = lessonPlanRepository.SaveLessonPlanConcept(lessonPlanConcept)
				if err != nil {
					return errors.New("Creating lesson plan concept error " + err.Error())
				}
			}

			err = processConceptSteps(rows, concept, lessonPlanConcept)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func processConceptSteps(conceptSteps [][]string, concept conceptModels.Concept, lessonPlanConcept models.LessonPlanConcept) error {

	for rowIndex, row := range conceptSteps[1:] /*Skip the header rwo*/ {
		if len(row) < 2 {
			fmt.Println(fmt.Sprintf("No code and step at row %d for concept %s", rowIndex, concept.Name))
			continue
		}

		code_ := row[0]
		step := row[1]
		code, err := strconv.Atoi(code_)

		if err != nil {
			return errors.New(fmt.Sprintf("Code error %s for step %s on concept %s", err.Error(), step, concept.Name))
		}

		err = lessonPlanRepository.SaveConceptStep(models.LessonPlanStep{
			Number:                  uuid.NewV4().String(),
			Step:                    step,
			Code:                    int64(code),
			LessonPlanConceptNumber: lessonPlanConcept.Number,
		})

		if err != nil {
			return errors.New(fmt.Sprintf("Error %s saving step %s on concept %s", err.Error(), step, concept.Name))
		}
	}

	return nil
}

func NewLessonPlanService() LessonPlanService {
	return &lessonPlanServiceImp{
		Db: config.DB,
	}
}
