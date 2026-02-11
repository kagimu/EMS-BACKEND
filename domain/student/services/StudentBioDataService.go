package services

import (
	"errors"
	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// UploadStudentBioData handles file upload for student bio data
func UploadStudentBioData(ctx *gin.Context, req request.StudentBioDataRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	fileName, err := services.UploadAttachment(ctx, "file")
	if err != nil {
		return err
	}

	db := config.DB.Begin()
	bioDataRepository := repositories.NewStudentBioDataRepository(db)

	// Create new bio data record
	bioData := models.StudentBioData{
		Number:        uuid.NewV4().String(),
		StudentNumber: req.StudentNumber,
		FileName:      fileName,
		FileType:      req.File.Header.Get("Content-Type"),
		FileSize:      req.File.Size,
	}

	if err = bioDataRepository.CreateStudentBioData(bioData); err != nil {
		db.Rollback()
		// Clean up the saved file if database operation fails
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		// Clean up the saved file if commit fails
		//_ = os.Remove(filePath)
		return err
	}

	return nil
}

// GetStudentBioDataByStudentNumber retrieves bio data for a student
func GetStudentBioDataByStudentNumber(studentNumber string) ([]models.StudentBioData, error) {
	if studentNumber == "" {
		return nil, errors.New("student_number is required")
	}

	bioDataRepository := repositories.NewStudentBioDataRepository(config.DB)
	bioData, err := bioDataRepository.GetStudentBioDataByStudentNumber(studentNumber)
	if err != nil {
		return nil, err
	}

	return bioData, nil
}

// DeleteStudentBioData deletes bio data and the associated file
func DeleteStudentBioData(studentNumber string) error {
	if studentNumber == "" {
		return errors.New("student_number is required")
	}

	db := config.DB.Begin()
	bioDataRepository := repositories.NewStudentBioDataRepository(db)

	// Get bio data to retrieve file path
	bioData, err := bioDataRepository.GetStudentBioDataByNumber(studentNumber)
	if err != nil {
		db.Rollback()
		return err
	}

	// Delete from database
	if err := bioDataRepository.DeleteStudentBioDataByNumber(studentNumber); err != nil {
		db.Rollback()
		return err
	}

	// Delete file from disk
	if bioData.FilePath != "" {
		_ = os.Remove(bioData.FilePath)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return err
	}

	return nil
}
