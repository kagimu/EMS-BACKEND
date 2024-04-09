package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/responses"
	schoolCurriculumModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func NewSchool(context *gin.Context) {

	var schoolInfo requests.SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := schoolInfo.Name
	curriculum := schoolInfo.Curriculum
	commencementDate := schoolInfo.CommencementDate

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if commencementDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Date of Commencement is required!"})
		return
	}

	if len(curriculum) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateSchool schoolCurriculumModels.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	tx := config.DB.Begin()

	// Parse the string into a time.Time value
	commencementDateData, err := time.Parse("2006-01-02", commencementDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse date" + err.Error()})
		return
	}

	///save user data
	saveSchool := schoolCurriculumModels.School{
		Number:           uuid.NewV4().String(),
		Name:             name,
		CommencementDate: commencementDateData,
		UserNumber:       userNumber,
	}

	if err := tx.Create(&saveSchool).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school", "details": err.Error()})
		return
	}

	for _, item := range curriculum {
		saveSchoolCurriculum := schoolCurriculumModels.SchoolCurriculum{
			Number:           uuid.NewV4().String(),
			CurriculumNumber: item,
			SchoolNumber:     saveSchool.Number,
		}

		if err := tx.Create(&saveSchoolCurriculum).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			//context.JSON(http.StatusBadRequest, gin.H{"message": err})
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save curriculum"})
			return
		}

	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school!!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School Added Successfully!", "data": saveSchool})
}

func EditSchool(context *gin.Context) {

	var schoolInfo requests.SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := schoolInfo.Name
	curriculum := schoolInfo.Curriculum
	commencementDate := schoolInfo.CommencementDate
	schoolNumber := schoolInfo.SchoolNumber

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School number is required!"})
		return
	}

	if commencementDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Commencement date is required!"})
		return
	}

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateSchool schoolCurriculumModels.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		Where("number != ?", schoolNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	// Parse the string into a time.Time value
	commencementDateData, err := time.Parse("2006-01-02", commencementDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	///save user data
	tx := config.DB.Begin()

	var schoolData schoolCurriculumModels.School
	updateSchool := tx.Model(&schoolData).
		Where("number = ?", schoolNumber).
		Updates(schoolCurriculumModels.School{Name: name, UserNumber: userNumber, CommencementDate: commencementDateData})
	if updateSchool.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	//this will add curriculum if it is new
	for _, item := range curriculum {

		var checkCurriculum schoolCurriculumModels.SchoolCurriculum
		if duplicate := tx.
			Where("curriculum_number = ?", item).
			Where("school_number = ?", schoolNumber).
			First(&checkCurriculum).Error; duplicate != nil {

			//		add curriculum if it doesnt exist
			saveSchoolCurriculum := schoolCurriculumModels.SchoolCurriculum{
				Number:           uuid.NewV4().String(),
				CurriculumNumber: item,
				SchoolNumber:     schoolNumber,
			}

			if err := tx.Create(&saveSchoolCurriculum).Error; err != nil {
				// If any error occurs, rollback the transaction
				tx.Rollback()
				// Handle or return the error as needed
				//context.JSON(http.StatusBadRequest, gin.H{"message": err})
				context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save curriculum"})
				return
			}

		}

	}

	//find and delete Curriculum thata are deselected
	var curriculumLst []schoolCurriculumModels.SchoolCurriculum
	if err := config.DB.
		Where("school_number = ?", schoolNumber).
		Find(&curriculumLst).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	var curriculumListToDelete []string
	for _, item := range curriculumLst {
		found := false
		for _, item2 := range curriculum {
			if item.CurriculumNumber == item2 {
				found = true
				break
			}
		}
		if !found {
			curriculumListToDelete = append(curriculumListToDelete, item.CurriculumNumber)
		}

	}

	//continue to delete
	if len(curriculumListToDelete) > 0 {
		for _, curr := range curriculumListToDelete {

			var curriculumData schoolCurriculumModels.SchoolCurriculum
			if err := tx.
				Unscoped().
				Where("curriculum_number = ?", curr).
				Where("school_number = ?", schoolNumber).
				Delete(&curriculumData).Error; err != nil {
				// If any error occurs, rollback the transaction
				tx.Rollback()
				// Handle or return the error as needed
				context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete Curriculum"})
				return
			}

		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update school!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School Edited Successfully!", "data": schoolData})
}

func DeleteSchool(context *gin.Context) {
	type SchoolInfo struct {
		SchoolNumber string `json:"school_number"`
	}

	var schoolInfo SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	schoolNumber := schoolInfo.SchoolNumber

	userNumber := context.Request.Header["User-Number"][0]

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	db := config.DB

	var schoolChecker models2.Campus
	if err := db.
		Where("school_number = ?", schoolNumber).
		First(&schoolChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this School, it has Campus"})
		return
	}

	var schoolData schoolCurriculumModels.School

	classLevelInfo := config.DB.
		Where("number = ?", schoolNumber).
		First(&schoolData)
	if classLevelInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "School not found!"})
		return
	}

	db.Unscoped().Where("number = ?", schoolNumber).Delete(&schoolData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School deleted successfully!"})
}

func FetchSchoolInfo(context *gin.Context) {

	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var schools schoolCurriculumModels.School
	var curriculum_ []responses.SchoolCurriculumResponseInfo
	if err := config.DB.
		Where("number = ?", schoolNumber).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		First(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	for _, schoolCurr := range schools.SchoolCurriculum {
		curriculum_ = append(curriculum_, responses.SchoolCurriculumResponseInfo{
			Name:   schoolCurr.Curriculum.Name,
			Number: schoolCurr.Curriculum.Number,
		})
	}

	schools_ := responses.SchoolResponseInfo{
		Number:           schools.Number,
		Name:             schools.Name,
		CommencementDate: schools.CommencementDate.Format("2006-01-02"),
		SchoolCurriculum: curriculum_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools_})
}

func FetchSchools(context *gin.Context) {

	var schools []schoolCurriculumModels.School
	var schools_ []responses.SchoolResponseInfo
	if err := config.DB.
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Order("created_at desc").
		Find(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	for _, item := range schools {

		var curriculum_ []responses.SchoolCurriculumResponseInfo

		for _, schoolCurr := range item.SchoolCurriculum {
			curriculum_ = append(curriculum_, responses.SchoolCurriculumResponseInfo{
				Name:   schoolCurr.Curriculum.Name,
				Number: schoolCurr.Number,
			})
		}

		schools_ = append(schools_, responses.SchoolResponseInfo{
			Number:           item.Number,
			Name:             item.Name,
			CommencementDate: item.CommencementDate.Format("2-Jan-2006"),
			SchoolCurriculum: curriculum_,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools_})
}
