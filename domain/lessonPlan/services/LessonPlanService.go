package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
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
	NewLessonPlan(lessonPlanInfo requests.LessonPlanInfo, userNumber string) error
	AddLessonPlanReferences(lessonPlanReferenceRequest requests.LessonPlanReferenceRequest, userNumber string) error
	UpdateLessonPlanReference(lessonPlanReferenceRequest requests.Material, referenceNumber, userNumber string) error
	DeleteLessonPlanReference(referenceNumber string) error
	Invalidate(lessonPlanNumber string, user *authModel.User) error
}

type lessonPlanServiceImp struct {
	Db *gorm.DB
}

func (l *lessonPlanServiceImp) Invalidate(lessonPlanNumber string, user *authModel.User) error {
	lessonPlanRepository.UseTx(config.DB)

	lessonPlan, err := lessonPlanRepository.FindLessonPlanByNumber(lessonPlanNumber)
	if err != nil {
		return err
	}

	if lessonPlan.Number == "" {
		return errors.New("Invalid lesson plan number ")
	}

	if lessonPlan.Status == constants.StatusInactive {
		return errors.New("Lesson plan already in invalidated ")
	}

	lessonPlan.Status = constants.StatusInactive
	lessonPlan.UserNumber = user.Number

	err = lessonPlanRepository.SaveLessonPlan(lessonPlan)
	if err != nil {
		return err
	}

	return nil
}

func (l *lessonPlanServiceImp) UpdateLessonPlanReference(material requests.Material, referenceNumber, userNumber string) error {
	lessonPlanRepository.UseTx(config.DB)
	lessonPlanReference, err := lessonPlanRepository.FindLessonPlanReferenceByNumber(referenceNumber)
	if err != nil {
		return errors.New("Checking lesson plan reference failed: " + err.Error())
	}

	if lessonPlanReference.Number == "" {
		return errors.New("Invalid Lesson plan reference number ")
	}

	duplicateByCode, err := lessonPlanRepository.FindLessonPlanReferenceByCode(lessonPlanReference.LessonPlanNumber, material.LearningMaterialNumber, material.Code)
	if err != nil {
		return errors.New("Checking duplicate reference code failed : " + err.Error())
	}

	if duplicateByCode.Number != "" && duplicateByCode.Number != referenceNumber {
		return errors.New("Code already used ")
	}

	duplicateByName, err := lessonPlanRepository.FindLessonPlanReferenceByName(lessonPlanReference.LessonPlanNumber, material.LearningMaterialNumber, material.Name)
	if err != nil {
		return errors.New("Checking duplicate reference name failed : " + err.Error())
	}

	if duplicateByName.Number != "" && duplicateByCode.Number != referenceNumber {
		return errors.New("Name already added ")
	}

	lessonPlanReference.LearningMaterialNumber = material.LearningMaterialNumber
	lessonPlanReference.Name = material.Name
	lessonPlanReference.Code = material.Code
	lessonPlanReference.UserNumber = userNumber

	err = lessonPlanRepository.SaveLessonPlanReference(lessonPlanReference)
	if err != nil {
		return errors.New("Updating lesson plan reference failed : " + err.Error())
	}

	return nil
}

func (l *lessonPlanServiceImp) DeleteLessonPlanReference(referenceNumber string) error {
	lessonPlanRepository.UseTx(config.DB)
	reference, err := lessonPlanRepository.FindLessonPlanReferenceByNumber(referenceNumber)
	if err != nil {
		return errors.New("Checking reference number failed: " + err.Error())
	}

	if reference.Number == "" {
		return errors.New("lesson plan reference not found ")
	}

	err = lessonPlanRepository.DeleteLessonPlanReference(reference)
	if err != nil {
		return errors.New("Removing lesson plan reference error: " + err.Error())
	}

	return nil
}

func (l *lessonPlanServiceImp) AddLessonPlanReferences(lessonPlanReferenceRequest requests.LessonPlanReferenceRequest, userNumber string) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		lessonPlanRepository.UseTx(tx)

		lessonPlan, err := lessonPlanRepository.FindLessonPlanByNumber(lessonPlanReferenceRequest.LessonPlanNumber)
		if err != nil {
			return errors.New("Checking lesson plan existence failed : " + err.Error())
		}

		if lessonPlan.Number == "" {
			return errors.New("Invalid lesson plan number ")
		}

		for _, reference := range lessonPlanReferenceRequest.References {

			learningMaterial, err := repositories4.FindLearningMaterialByNumber(reference.LearningMaterialNumber, tx)
			if err != nil {
				return errors.New("Checking learningMaterial existence failed : " + err.Error())
			}

			if learningMaterial.Number == "" {
				return errors.New("Invalid learningMaterial number ")
			}

			duplicateByCode, err := lessonPlanRepository.FindLessonPlanReferenceByCode(lessonPlan.Number, learningMaterial.Number, reference.Code)
			if err != nil {
				return errors.New("Checking duplicate reference code failed : " + err.Error())
			}

			if duplicateByCode.Number != "" {
				return errors.New("Code already used ")
			}

			duplicateByName, err := lessonPlanRepository.FindLessonPlanReferenceByName(lessonPlan.Number, learningMaterial.Number, reference.Name)
			if err != nil {
				return errors.New("Checking duplicate reference name failed : " + err.Error())
			}

			if duplicateByName.Number != "" {
				return errors.New("Name already added ")
			}

			LessonPlanReference := models.LessonPlanReference{
				LessonPlanNumber:       lessonPlan.Number,
				LearningMaterialNumber: learningMaterial.Number,
				Name:                   reference.Name,
				Code:                   reference.Code,
				UserNumber:             userNumber,
				Number:                 uuid.NewV4().String(),
			}

			err = lessonPlanRepository.SaveLessonPlanReference(LessonPlanReference)
			if err != nil {
				return errors.New("Saving lesson plan reference failed : " + err.Error())
			}

		}

		return nil
	})

	return err
}

func (l *lessonPlanServiceImp) NewLessonPlan(lessonPlanInfo requests.LessonPlanInfo, userNumber string) error {

	schemeOfWork, err := repositories2.GetSchemeOfWorkBySchemeOfWorkSubTopicNumber(lessonPlanInfo.SchemeSubTopicNumber)
	if err != nil {
		return err
	}

	if schemeOfWork.Number == "" {
		return errors.New("Invalid scheme of work sub topic number ")
	}

	if schemeOfWork.Status != constants.StatusActive {
		return errors.New("The scheme is not active ")
	}

	lessonPlanRepository.UseTx(config.DB)
	checkDuplicate, err := lessonPlanRepository.GetLessonPlanBySchemeOfWorkSubTopicNumber(lessonPlanInfo.SchemeSubTopicNumber)
	if err != nil {
		return err
	}

	if checkDuplicate.Number != "" && checkDuplicate.Number != lessonPlanInfo.LessonPlanNumber && checkDuplicate.Status == constants.StatusActive {
		return errors.New("lesson plan already exists")
	}

	lessonPlanNumber := lessonPlanInfo.LessonPlanNumber

	if lessonPlanNumber == "" {
		lessonPlanNumber = uuid.NewV4().String()
	}

	saveLessonPlan := models.LessonPlan{
		Number:               lessonPlanNumber,
		SchemeSubTopicNumber: lessonPlanInfo.SchemeSubTopicNumber,
		Status:               constants.StatusActive,
		UserNumber:           userNumber,
		UpdatedAt:            time.Now(),
		CreatedAt:            time.Now(),
	}

	oldLessonPlan, err := lessonPlanRepository.FindLessonPlanByNumber(lessonPlanNumber)
	if err != nil {
		return errors.New("Checking old lesson plan existence failed : " + err.Error())
	}

	if oldLessonPlan.Number != "" {
		saveLessonPlan.CreatedAt = oldLessonPlan.CreatedAt
		saveLessonPlan.Status = oldLessonPlan.Status
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		lessonPlanRepository.UseTx(tx)
		conceptRepository.UseTx(tx)

		err = lessonPlanRepository.SaveLessonPlan(saveLessonPlan)
		if err != nil {
			return err
		}

		for _, concept := range lessonPlanInfo.LessonPlanConceptInfo {

			conceptNumber := concept.LessonPlanConceptNumber
			if conceptNumber == "" {
				conceptNumber = uuid.NewV4().String()
			}

			saveLessonConcepts := models.LessonPlanConcept{
				Number:           conceptNumber,
				ConceptNumber:    concept.LessonPlanConceptNumber,
				Code:             concept.Code,
				LessonPlanNumber: saveLessonPlan.Number,
				UserNumber:       userNumber,
				UpdatedAt:        time.Now(),
				CreatedAt:        time.Now(),
			}

			//	adds new concepts when Concept number doesnt exists
			err = lessonPlanRepository.SaveLessonPlanConcept(saveLessonConcepts)
			if err != nil {
				return errors.New("Failed to save concept " + err.Error())
			}

			lessonPlanSteps := concept.StepList

			for _, step := range lessonPlanSteps {

				stepNumber := step.Number
				if stepNumber == "" {
					stepNumber = uuid.NewV4().String()
				}

				if step.Step != "" {
					saveLearningSteps := models.LessonPlanStep{
						Number:                  stepNumber,
						Code:                    step.Code,
						Step:                    step.Step,
						LessonPlanConceptNumber: conceptNumber,
						UpdatedAt:               time.Now(),
						CreatedAt:               time.Now(),
					}

					err = lessonPlanRepository.SaveConceptStep(saveLearningSteps)

					if err != nil {
						return errors.New("Failed to save lesson plan step " + err.Error())
					}
				}

			}

		}

		return nil
	})

	return err
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

	if len(concepts) == 0 {
		return nil, "", errors.New(fmt.Sprintf("Sub-topic '%s' has no concepts for lesson planning", subTopic.Name))
	}

	defaultOpenSheet := 0

	for _, concept := range concepts {

		sheetName := fmt.Sprintf("Sheet-%d", concept.Code)

		defaultOpenSheet, err = fileTemplate.NewSheet(sheetName)
		if err != nil {
			return nil, "", errors.New(fmt.Sprintf("Error '%s' creating sheet '%s'", err.Error(), concept.Name))
		}

		sheetsHeaders := [][]string{
			{"Concept name:", concept.Name},
			{"Code", "Step"},
		}

		for i, row := range sheetsHeaders {
			for j, cell := range row {
				cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
				err = fileTemplate.SetCellValue(sheetName, cellRef, cell)
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
			rows, err := excelFile.GetRows(fmt.Sprintf("Sheet-%d", concept.Code))
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

	for rowIndex, row := range conceptSteps[2:] /*Skip the header rwo*/ {
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
