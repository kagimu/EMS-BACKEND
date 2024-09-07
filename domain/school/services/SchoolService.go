package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RemoveSchool(schoolNumber string) error {

	school, err := repositories.GetSchoolByNumber(schoolNumber)
	if err != nil {
		return errors.New("checking school existence: " + err.Error())
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	campuses, err := repositories2.FindCampuses(schoolNumber)
	if err != nil {
		return errors.New("checking campuses error: " + err.Error())
	}

	if len(campuses) > 0 {
		return errors.New("this school already have campuses, you may need to first remove them")
	}

	schoolCurriculums, err := repositories.FindSchoolCurriculumBySchoolNumber(schoolNumber)
	if err != nil {
		return errors.New("getting school curriculums error: " + err.Error())
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, schoolCurriculum := range schoolCurriculums {
			err := repositories.RemoveSchoolCurriculum(schoolCurriculum, tx)
			if err != nil {
				return errors.New("removing school curriculum error: " + err.Error())
			}
		}

		err := repositories.RemoveSchool(school, tx)
		if err != nil {
			return errors.New("removing school error: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func UploadSchoolLogo(schoolNumber string, ctx *gin.Context) error {
	repo := repositories.NewSchoolRepositoryImpl(config.DB)
	school, err := repo.FindSchoolByNumber(schoolNumber)
	if err != nil {
		return errors.New("checking school existence: " + err.Error())
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	logo, err := commons.HandleFileUpload(ctx, "logo", "school_logo")
	if err != nil {
		return errors.New("uploading school logo: " + err.Error())
	}

	school.Logo = logo

	err = repo.UpdateSchool(school)
	if err != nil {
		return errors.New("updating school logo: " + err.Error())
	}

	return nil

}
