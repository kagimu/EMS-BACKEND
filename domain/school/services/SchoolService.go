package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/responses"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
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

	logo, err := services.UploadAttachment(ctx, "logo") // commons.HandleFileUpload(ctx, "logo", "school_logos")
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

func UploadPrincipalSignature(schoolNumber string, ctx *gin.Context) error {
	repo := repositories.NewSchoolRepositoryImpl(config.DB)
	school, err := repo.FindSchoolByNumber(schoolNumber)
	if err != nil {
		return errors.New("checking school existence: " + err.Error())
	}

	if school.Number == "" {
		return errors.New("invalid school number")
	}

	signature, err := services.UploadAttachment(ctx, "signature")
	if err != nil {
		return errors.New("uploading principal signature: " + err.Error())
	}

	school.PrincipalSignature = signature

	err = repo.UpdateSchool(school)
	if err != nil {
		return errors.New("updating principal signature: " + err.Error())
	}

	return nil
}

func GetLatestSchoolLogo(schoolNumber string) (string, error) {
	repo := repositories.NewSchoolRepositoryImpl(config.DB)
	school, err := repo.FindSchoolByNumber(schoolNumber)
	if err != nil {
		return "", errors.New("checking school existence: " + err.Error())
	}

	if school.Number == "" {
		return "", errors.New("invalid school number")
	}

	return school.Logo, nil
}

func BuildLeanSchoolResponse(school *schoolModel.School) responses.LeanSchool {
	return responses.LeanSchool{
		Number: school.Number,
		Name:   school.Name,
		Logo:   commons.GenerateDownloadPath(school.Logo, false),
	}
}

func FetchSchoolsByEntityNumber(entityNumber string, ctx *gin.Context) ([]responses2.SchoolResponse, error) {

	//repo := repositories.NewSchoolRepositoryImpl(config.DB)
	//_, err := repo.FindSchoolsByEntityNumber(entityNumber, ctx)
	//if err != nil {
	//	return nil, errors.New("getting schools by entity number: " + err.Error())
	//}

	var schools []schoolModel.School
	var schools_ []responses2.SchoolResponse
	if err := config.DB.
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Order("created_at desc").
		Scopes(
			//filters.ByScopedItems("school", ctx),
			filters.ByEntityNumber(entityNumber),
		).
		Find(&schools).Error; err != nil {
		return nil, err
	}

	for _, item := range schools {

		var curriculum_ []responses2.SchoolCurriculumResponse

		for _, schoolCurr := range item.SchoolCurriculum {
			curriculum_ = append(curriculum_, responses2.SchoolCurriculumResponse{
				Name:   schoolCurr.Curriculum.Name,
				Number: schoolCurr.Number,
			})
		}

		schools_ = append(schools_, responses2.SchoolResponse{
			Number:             item.Number,
			Name:               item.Name,
			CommencementDate:   item.CommencementDate.Format("2-Jan-2006"),
			SchoolCurriculum:   curriculum_,
			Logo:               commons.GenerateDownloadPath(item.Logo, false),
			PrincipalSignature: commons.GenerateDownloadPath(item.PrincipalSignature, false),
			Motto:              item.Motto,
			PrincipalName:      item.PrincipalName,
			PrincipalTitle:     item.PrincipalTitle,
		})
	}

	//var schools = make([]responses2.SchoolResponse, 0)
	//
	//for _, school := range schls {
	//
	//	var curriculum_ []responses2.SchoolCurriculumResponse
	//	schoolCurriculums, err := repositories.FindSchoolCurriculumBySchoolNumber(school.Number)
	//
	//	for _, schoolCurr := range item.SchoolCurriculum {
	//		curriculum_ = append(curriculum_, responses.SchoolCurriculumResponseInfo{
	//			Name:   schoolCurr.Curriculum.Name,
	//			Number: schoolCurr.Number,
	//		})
	//	}
	//
	//	schools = append(schools, responses2.SchoolResponse{
	//		Number:           school.Number,
	//		Name:             school.Name,
	//		CommencementDate: school.CommencementDate,
	//		SchoolCurriculum: nil,
	//		Logo:             commons.GenerateDownloadPath(school.Logo, false),
	//	})
	//}

	return schools_, nil
}
