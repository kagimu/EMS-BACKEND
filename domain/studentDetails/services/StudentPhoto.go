package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/repository"
	"github.com/gin-gonic/gin"
)

func UploadStudentPhoto(studentNumber string, userNumber string, ctx *gin.Context) error {
	repo := repository.NewStudentPhotoRepo(config.DB)

	photo, err := services.UploadAttachment(ctx, "photo")
	if err != nil {
		return err
	}

	studentPhoto := studentDetails.StudentPhoto{
		StudentNumber: studentNumber,
		UserNumber:    userNumber,
		Photo:         photo,
	}

	err = repo.CreateStudentPhoto(&studentPhoto)
	if err != nil {
		return err
	}

	return nil
}
