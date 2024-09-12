package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
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

	//TODO Upate student photo in student table

	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		return err
	}

	student.Photo = photo

	err = config.DB.Save(&student).Error

	if err != nil {
		return err
	}

	return nil
}
