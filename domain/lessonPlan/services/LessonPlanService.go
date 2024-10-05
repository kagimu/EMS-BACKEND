package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"
	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LessonPlanService interface {
	ImportPlan(lessonPlanFile multipart.File, subTopicNumber, currentUserNumber string) error
	ExportLessonPlanTemplate(schemeSubTopicNumber string) (*excelize.File, string, error)
}

type lessonPlanServiceImp struct {
	Db *gorm.DB
}

func (l *lessonPlanServiceImp) ExportLessonPlanTemplate(schemeSubTopicNumber string) (*excelize.File, string, error) {
	fileTemplate := excelize.NewFile()

	schemeOfWorkSubTopic, err := repositories2.GetSchemeOfWorkSubTopicByNumber(schemeSubTopicNumber)
	if err != nil {
		return nil, "", errors.New("Checking scheme subtopic error: " + err.Error())
	}

	if schemeOfWorkSubTopic.Number == "" {
		return nil, "", errors.New("Invalid scheme subtopic, we can not export anything ")
	}

	subTopic, err := repositories3.GetSubTopicByNumber(schemeOfWorkSubTopic.SubTopicNumber)
	if err != nil {
		return nil, "", errors.New("Checking subtopic error " + err.Error())
	}

	concepts, err := resolveConcepts(subTopic.Number, config.DB)
	if err != nil {
		return nil, "", err
	}

	sheetsHeaders := [][]string{
		{"Code", "Step"},
	}

	defaultOpenSheet := 0

	for _, concept := range concepts {

		defaultOpenSheet, err = fileTemplate.NewSheet(concept.Name)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error '%s' creating sheet '%s'", err.Error(), concept.Name))
			continue
		}

		for i, row := range sheetsHeaders {
			for j, cell := range row {
				cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
				err = fileTemplate.SetCellValue(concept.Name, cellRef, cell)
				if err != nil {
					return nil, "", errors.New(fmt.Sprintf("Error '%s' writing sheet '%s' headers", err.Error(), concept.Name))
				}
			}
		}

	}

	fileTemplate.SetActiveSheet(defaultOpenSheet)

	unixTime := time.Now().Unix()
	return fileTemplate, fmt.Sprintf("%s-%d.xlsx", toFileName(subTopic.Name), unixTime), nil
}

func toFileName(subTopicName string) string {
	lowercase := strings.ToLower(subTopicName)

	re := regexp.MustCompile(`[^a-z0-9]+`)

	result := re.ReplaceAllString(lowercase, "-")

	return strings.Trim(result, "-")
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
	concepts, err := resolveConcepts(schemeOfWorkSubTopic.SubTopicNumber, config.DB)
	if err != nil {
		return err
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

func resolveConcepts(subTopicNumber string, Db *gorm.DB) ([]conceptModels.Concept, error) {
	conceptRepository.UseTx(Db)
	concepts, err := conceptRepository.FindConceptsBySubTopicNumber(subTopicNumber)
	if err != nil {
		return concepts, errors.New("Getting concepts error " + err.Error())
	}

	if len(concepts) == 0 {
		return concepts, errors.New("The selected sub-topic has no concepts ")
	}

	return concepts, nil
}

func processConceptSteps(conceptSteps [][]string, concept conceptModels.Concept, lessonPlanConcept models.LessonPlanConcept) error {

	var lessonPlanSteps = make([]models.ExtractedStep, 0)

	for rowIndex, row := range conceptSteps[1:] /*Skip the header rwo*/ {
		if len(row) < 2 {
			fmt.Println(fmt.Sprintf("No code and step at row %d for concept %s", rowIndex, concept.Name))
			continue
		}

		rowCode := row[0]
		step := row[1]
		code_, err := strconv.Atoi(rowCode)

		if err != nil {
			return errors.New(fmt.Sprintf("Code error %s for step %s on concept %s", err.Error(), step, concept.Name))
		}
		code := int64(code_)

		lessonPlanSteps = append(lessonPlanSteps, models.ExtractedStep{
			Code: code,
			Step: step,
		})

	}

	err := checkStepDuplicates(lessonPlanSteps)
	if err != nil {
		return err
	}

	for _, extractedStep := range lessonPlanSteps {

		currentStep, err := lessonPlanRepository.FindConceptStepByCode(extractedStep.Code, lessonPlanConcept.Number)
		if err != nil {
			return errors.New("Checking code duplicate error " + err.Error())
		}

		if currentStep.Number == "" {

			currentStep, err = lessonPlanRepository.FindConceptStepByStep(extractedStep.Step, lessonPlanConcept.Number)
			if err != nil {
				return errors.New("Checking step duplicate error " + err.Error())
			}

			if currentStep.Number == "" {
				currentStep = models.LessonPlanStep{
					Number:                  uuid.NewV4().String(),
					Step:                    extractedStep.Step,
					Code:                    extractedStep.Code,
					LessonPlanConceptNumber: lessonPlanConcept.Number,
				}
			} else {
				currentStep.Code = extractedStep.Code
			}

		} else {
			currentStep.Step = extractedStep.Step
		}

		err = lessonPlanRepository.SaveConceptStep(currentStep)
		if err != nil {
			return errors.New(fmt.Sprintf("Error %s saving step %s on concept %s", err.Error(), extractedStep.Step, concept.Name))
		}
	}

	return nil
}

func checkStepDuplicates(steps []models.ExtractedStep) error {
	var cleanedSteps = make([]string, len(steps))
	var cleanedCodes = make([]int64, len(steps))

	for _, step := range steps {

		for _, cleanedStep := range cleanedSteps {
			if cleanedStep == step.Step {
				return errors.New("Step '" + cleanedStep + "' already defined")
			}
		}
		cleanedSteps = append(cleanedSteps, step.Step)

		for _, cleanedCode := range cleanedCodes {
			if cleanedCode == step.Code {
				return errors.New(fmt.Sprintf("Code %d already used", cleanedCode))
			}
		}
		cleanedCodes = append(cleanedCodes, step.Code)

	}

	return nil
}

func NewLessonPlanService() LessonPlanService {
	return &lessonPlanServiceImp{
		Db: config.DB,
	}
}
