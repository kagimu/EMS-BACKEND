package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SchemeOfWorkService struct {
}

var repo = repository.NewLessonPlanRepository()

func (s SchemeOfWorkService) InvalidateSchemeOfWork(schemeOfWorkNumber, userNumber string) error {
	schemeOfWork, err := repositories.GetSchemeOfWorkByNumber(schemeOfWorkNumber)
	if err != nil {
		return err
	}

	if schemeOfWork.Number == "" {
		return errors.New("Invalid schemeOfWork number ")
	}

	if schemeOfWork.Status == constants.StatusInactive {
		return errors.New("This scheme is already invalidated ")
	}

	//Find active lesson plan
	repo.UseTx(config.DB)
	lessonPlane, err := repo.FindLessonPlanBySchemeOfWorkNumber(schemeOfWork.Number)
	if err != nil {
		return err
	}

	if lessonPlane.Number != "" && lessonPlane.Status != constants.StatusInactive {
		return errors.New("This scheme of work already have a lesson plan ")
	}

	schemeOfWork.Status = constants.StatusInactive
	schemeOfWork.UserNumber = userNumber

	err = repositories.SaveSchemeOfWork(schemeOfWork)
	if err != nil {
		return err
	}

	return nil
}

func (s SchemeOfWorkService) ActivateSchemeOfWork(schemeOfWorkNumber, userNumber string) error {
	schemeOfWork, err := repositories.GetSchemeOfWorkByNumber(schemeOfWorkNumber)
	if err != nil {
		return err
	}

	if schemeOfWork.Number == "" {
		return errors.New("Invalid schemeOfWork number ")
	}

	if schemeOfWork.Status == constants.StatusActive {
		return errors.New("This scheme is already active ")
	}

	schemeOfWork.Status = constants.StatusActive
	schemeOfWork.UserNumber = userNumber

	err = repositories.SaveSchemeOfWork(schemeOfWork)
	if err != nil {
		return err
	}

	return nil
}

func (s SchemeOfWorkService) Update(schemeOfWorkInfo requests.SchemeOfWorkInfo, userNumber string) error {

	name := schemeOfWorkInfo.Name
	classNumber := schemeOfWorkInfo.ClassNumber
	subjectNumber := schemeOfWorkInfo.SubjectNumber
	subtopicList := schemeOfWorkInfo.SubtopicList

	scheme, err := repositories.GetSchemeOfWorkByNumber(schemeOfWorkInfo.SchemeNumber)
	if err != nil {
		return err
	}

	if scheme.Number == "" {
		return errors.New("Invalid schemeOfWork number ")
	}

	checkDuplicate, err := repositories.FindSchemeOfWorkByNameAndClassAndSubjectAndStatus(name, classNumber, subjectNumber, constants.StatusActive)
	if err != nil {
		return err
	}

	if checkDuplicate.Number != "" && checkDuplicate.Number != scheme.Number {
		return errors.New("This name is already in use ")
	}

	checkDuplicate.Name = name
	checkDuplicate.CampusClassNumber = classNumber
	checkDuplicate.AcademicYearNumber = schemeOfWorkInfo.AcademicYearNumber
	checkDuplicate.UpdatedAt = time.Now()
	checkDuplicate.CreatedAt = scheme.CreatedAt
	checkDuplicate.UserNumber = userNumber

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err = repositories.SaveSchemeOfWork(checkDuplicate, tx)
		if err != nil {
			return err
		}

		for _, schemeList := range subtopicList {
			codeNo, err := strconv.ParseInt(strconv.Itoa(int(schemeList.Code)), 10, 64)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}

			item, err := repositories.FindSchemeOfWorkSubTopicByNumber(schemeList.Number, tx)
			if err != nil {
				return err
			}

			if item.Number == "" {
				return errors.New("Invalid Item number for scheme of work subtopic ")
			}

			item.Code = codeNo
			item.EstimatedTime = schemeList.EstimatedTime
			item.SubTopicNumber = schemeList.Number
			item.UserNumber = userNumber
			item.UpdatedAt = time.Now()

			err = repositories.SaveSchemeOfWorkSubTopic(item, tx)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func InitSchemeOfWorkService() *SchemeOfWorkService {
	return &SchemeOfWorkService{}
}
