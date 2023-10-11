package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
)

func NewAcademicYear(context *gin.Context) {

	type AcademicYearInfo struct {
		AcademicYearName string   `json:"academic_year_name"`
		NoOfStudyPeriod  string   `json:"no_of_study_period"`
		AcademicYear     []string `json:"academic_year"`
		LevelNumber      string   `json:"level_number"`
		Code             string   `json:"code"`
	}

	var academicYearInfo AcademicYearInfo
	err := context.BindJSON(&academicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearName := academicYearInfo.AcademicYearName
	noOfStudyPeriod := academicYearInfo.NoOfStudyPeriod
	academicYear := academicYearInfo.AcademicYear
	levelNumber := academicYearInfo.LevelNumber
	code := academicYearInfo.Code
	userNumber := context.Request.Header["User-Number"][0]

	if academicYearName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year Name is required!"})
		return
	}

	if noOfStudyPeriod == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Number Of Study Periods is required!"})
		return
	}

	if len(academicYear) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level Year is required!"})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
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

	var duplicateCode academicYearModels.AcademicYear
	if duplicate := config.DB.
		Where("code = ?", code).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateAcademicYear academicYearModels.AcademicYear
	if duplicate := config.DB.
		Where("academic_year_name = ?", academicYearName).
		Where("level_number = ?", levelNumber).
		First(&duplicateAcademicYear).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year already exists"})
		return
	}

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, academicYear[0])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	endDate, errorDate := time.Parse(layout, academicYear[1])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	noOfStudyPeriodData, err := strconv.ParseInt(noOfStudyPeriod, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		//fmt.Println("Error:", err)
		//return
	}

	number := uuid.NewV4().String()
	///save user data
	saveAcademicYear := academicYearModels.AcademicYear{
		Number:                number,
		Code:                  code,
		AcademicYearName:      academicYearName,
		NoOfStudyPeriod:       int32(noOfStudyPeriodData),
		AcademicYearStartDate: startDate,
		AcademicYearEndDate:   endDate,
		UserNumber:            userNumber,
		LevelNumber:           levelNumber,
	}
	config.DB.Save(&saveAcademicYear)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Concept Added Successfully!", "data": saveAcademicYear})
}

func FetchAcademicYear(context *gin.Context) {

	levelNumber := context.Query("level")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var academicYearInfo []academicYearModels.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		//Select("name", "number").
		Where("level_number = ?", levelNumber).
		Order("code asc").
		Find(&academicYearInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
		return
	}

	// Populate the selected columns from the fetched data
	for _, academicList := range academicYearInfo {
		academicYearInfo_ = append(academicYearInfo_, responses.AcademicYearList{
			Number:                academicList.Number,
			Code:                  academicList.Code,
			AcademicYearName:      academicList.AcademicYearName,
			NoOfStudyPeriod:       academicList.NoOfStudyPeriod,
			AcademicYearStartDate: (academicList.AcademicYearStartDate).Format("2-Jan-2006 15:04:05"),
			AcademicYearEndDate:   (academicList.AcademicYearEndDate).Format("2-Jan-2006 15:04:05"),
			UserNumber:            academicList.UserNumber,
			LevelNumber:           academicList.LevelNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})

}
