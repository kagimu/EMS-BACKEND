package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/requests"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func NewCampus(context *gin.Context) {

	var campusInfo requests.CampusInfo
	err := context.BindJSON(&campusInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := campusInfo.Name
	phone := campusInfo.Phone
	email := campusInfo.Email
	location := campusInfo.Location
	campusCurriculum := campusInfo.CampusCurriculum
	schoolNumber := campusInfo.SchoolNumber
	commencementDate := campusInfo.CommencementDate

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if phone == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Phone is required!"})
		return
	}

	if commencementDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Commencement date is required!"})
		return
	}

	if email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
		return
	}

	if location == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Location is required!"})
		return
	}

	if len(campusCurriculum) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School number is required!"})
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

	var duplicateCampus campusModel.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	// Parse the string into a time.Time value
	commencementDateData, err := time.Parse("2006-01-02", commencementDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	saveCampus := campusModel.Campus{
		Number:           uuid.NewV4().String(),
		UserNumber:       userNumber,
		Name:             name,
		Phone:            phone,
		Email:            email,
		Location:         location,
		CommencementDate: commencementDateData,
		SchoolNumber:     schoolNumber,
	}

	if err := tx.Create(&saveCampus).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create campus",
			"err":     err.Error(),
		})
		return
	}

	//save Curriculum
	for _, item := range campusCurriculum {

		saveCurriculum := campusModel.CampusCurriculum{
			Number:                 uuid.NewV4().String(),
			SchoolCurriculumNumber: item,
			CampusNumber:           saveCampus.Number,
		}

		if err := tx.Create(&saveCurriculum).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{
				"err":     err.Error(),
				"message": "Failed to save Curriculum"})
			return
		}

	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add campus!!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
}

func NewCampusLevel(context *gin.Context) {

	var levelInfo requests.CampusLevelInfo
	err := context.BindJSON(&levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	levelNumber := levelInfo.LevelNumber
	campusCurriculumNumber := levelInfo.CampusCurriculumNumber

	tx := config.DB.Begin()
	var curriculumLevel schoolLevelModels.CurriculumLevel
	if err := tx.
		Where("number = ?", levelNumber).
		First(&curriculumLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	saveLevel := schoolLevelModels.CampusLevel{
		Number:                 uuid.NewV4().String(),
		Code:                   curriculumLevel.Code,
		CurriculumLevelNumber:  levelNumber,
		CampusCurriculumNumber: campusCurriculumNumber,
	}

	var checkLevel schoolLevelModels.CampusLevel
	if err := tx.
		Where("curriculum_level_number = ?", levelNumber).
		Where("campus_curriculum_number = ?", campusCurriculumNumber).
		First(&checkLevel).Error; err == nil {
		//delete

		if err := tx.
			Unscoped().
			Where("curriculum_level_number = ?", levelNumber).
			Where("campus_curriculum_number = ?", campusCurriculumNumber).
			Delete(&checkLevel).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete level"})
			return
		}

	} else {

		if err := tx.Create(&saveLevel).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			//context.JSON(http.StatusBadRequest, gin.H{"message": err})
			context.JSON(http.StatusBadRequest, gin.H{
				"err":     err.Error(),
				"message": "Failed to create level"})
			return
		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add level!!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level Added Successfully!", "data": saveLevel})

}

func EditCampus(context *gin.Context) {

	var campusInfo requests.CampusInfo
	err := context.BindJSON(&campusInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := campusInfo.Name
	phone := campusInfo.Phone
	email := campusInfo.Email
	location := campusInfo.Location
	schoolNumber := campusInfo.SchoolNumber
	campusCurriculum := campusInfo.CampusCurriculum
	campusNumber := campusInfo.CampusNumber

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if phone == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Phone is required!"})
		return
	}

	if email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
		return
	}

	if location == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Location is required!"})
		return
	}

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School number is required!"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required!"})
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

	var duplicateCampus campusModel.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		Where("number != ?", campusNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	tx := config.DB.Begin()

	var campusDataInfo campusModel.Campus
	updateLevel := tx.Model(&campusDataInfo).
		Where("number = ?", campusNumber).
		Updates(campusModel.Campus{Name: name, Phone: phone, UserNumber: userNumber, Email: email, Location: location})
	if updateLevel.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	for _, item := range campusCurriculum {
		var checkCurriculum campusModel.CampusCurriculum
		if duplicate := tx.
			Where("school_curriculum_number = ?", item).
			Where("campus_number = ?", campusNumber).
			First(&checkCurriculum).Error; duplicate != nil {

			saveCurriculum := campusModel.CampusCurriculum{
				Number:                 uuid.NewV4().String(),
				SchoolCurriculumNumber: item,
				CampusNumber:           campusNumber,
			}
			if err := tx.Create(&saveCurriculum).Error; err != nil {
				// If any error occurs, rollback the transaction
				tx.Rollback()
				// Handle or return the error as needed
				context.JSON(http.StatusBadRequest, gin.H{
					"err":     err.Error(),
					"message": "Failed to save Curriculum"})
				return
			}
		}
	}

	var curriculumLst []campusModel.CampusCurriculum
	if err := config.DB.
		Where("campus_number = ?", campusNumber).
		Find(&curriculumLst).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	var curriculumListToDelete []string
	for _, item := range curriculumLst {
		found := false
		for _, item2 := range campusCurriculum {
			if item.SchoolCurriculumNumber == item2 {
				found = true
				break
			}
		}
		if !found {
			curriculumListToDelete = append(curriculumListToDelete, item.SchoolCurriculumNumber)
		}

	}

	if len(curriculumListToDelete) > 0 {
		for _, curr := range curriculumListToDelete {

			var curriculumData campusModel.CampusCurriculum
			if err := tx.
				Unscoped().
				Where("school_curriculum_number = ?", curr).
				Where("campus_number = ?", campusNumber).
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

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Edited Successfully!", "data": campusDataInfo})
}

func FetchCampuses(context *gin.Context) {

	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var campuses []campusModel.Campus
	var campuses_ []responses.CampusDataList
	if err := config.DB.
		Where("school_number = ?", schoolNumber).
		Order("created_at desc").
		Find(&campuses).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus found"})
		return
	}

	for _, item := range campuses {

		var campusCurriculum []campusModel.CampusCurriculum
		var curriculum_ []responses.CurriculumDataList
		if err := config.DB.
			Where("campus_number = ?", item.Number).
			Preload("SchoolCurriculum.Curriculum").
			Preload("SchoolCurriculum").
			Find(&campusCurriculum).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
			return
		}

		for _, items := range campusCurriculum {
			curriculum_ = append(curriculum_, responses.CurriculumDataList{
				Number: items.SchoolCurriculumNumber,
				Name:   items.SchoolCurriculum.Curriculum.Name,
			})
		}

		campuses_ = append(campuses_, responses.CampusDataList{
			Number:           item.Number,
			Name:             item.Name,
			Phone:            item.Phone,
			Email:            item.Email,
			Location:         item.Location,
			UserNumber:       item.UserNumber,
			SchoolNumber:     item.SchoolNumber,
			CommencementDate: item.CommencementDate.Format("2-Jan-2006"),
			Curriculum:       curriculum_,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campuses_})
}

func CampusInfo(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var campus campusModel.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Order("created_at desc").
		First(&campus).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus found"})
		return
	}

	//fetch curriculum
	var campusCurriculum []campusModel.CampusCurriculum
	var curriculum_ []responses.CurriculumDataList
	if err := config.DB.
		Where("campus_number = ?", campus.Number).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Find(&campusCurriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	for _, items := range campusCurriculum {
		curriculum_ = append(curriculum_, responses.CurriculumDataList{
			Number: items.SchoolCurriculumNumber,
			Name:   items.SchoolCurriculum.Curriculum.Name,
		})
	}

	campus_ := responses.CampusDataList{
		Number:           campus.Number,
		Name:             campus.Name,
		Phone:            campus.Phone,
		Email:            campus.Email,
		Location:         campus.Location,
		UserNumber:       campus.UserNumber,
		SchoolNumber:     campus.SchoolNumber,
		CommencementDate: campus.CommencementDate.Format("2006-01-02"),
		Curriculum:       curriculum_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campus_})
}

func CampusCurriculum(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	//fetch curriculum
	var campusCurriculum []campusModel.CampusCurriculum
	var curriculum_ []responses.CurriculumDataList
	if err := config.DB.
		Where("campus_number = ?", campusNumber).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Find(&campusCurriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	for _, items := range campusCurriculum {
		curriculum_ = append(curriculum_, responses.CurriculumDataList{
			Number: items.Number,
			Name:   items.SchoolCurriculum.Curriculum.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculum_})
}

func FetchLevel(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus Number is required!"})
		return
	}

	var campus campusModel.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Preload("Level.LearningLevel").
		Preload("Level").
		Order("created_at desc").
		First(&campus).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Campus found"})
		return
	}

	campus_ := responses.CampusDataList{
		Number:       campus.Number,
		Name:         campus.Name,
		Phone:        campus.Phone,
		Email:        campus.Email,
		Location:     campus.Location,
		UserNumber:   campus.UserNumber,
		SchoolNumber: campus.SchoolNumber,
		//Level:        level_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campus_})

}

func FetchCampusLevel(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus Number is required!"})
		return
	}

	var levelInfo []schoolLevelModels.Level
	var levelInfo_ []schoolLevelModels.LevelDataInfo
	if err := config.DB.
		Where("campus_number = ?", campusNumber).
		Preload("LearningLevel").
		Order("created_at desc").
		Find(&levelInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level found"})
		return
	}

	// Populate the selected columns from the fetched data
	/*for _, level := range levelInfo {
		levelInfo_ = append(levelInfo_, campusModel.LevelDataInfo{
			Name:   level.LearningLevel.Name,
			Code:   level.LearningLevel.Code,
			Number: level.Number,
		})
	}*/

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levelInfo_})

}

func FetchCampusLevelList(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus Number is required!"})
		return
	}

	var campusLevelInfo []schoolLevelModels.CampusLevel
	var levelInfo_ []schoolLevelModels.LevelDataInfo
	if err := config.DB.
		Where("campus_number = ?", campusNumber).
		Preload("CurriculumLevel.Level").
		Preload("CurriculumLevel").
		Order("code asc").
		Find(&campusLevelInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level found"})
		return
	}

	// Populate the selected columns from the fetched data
	for _, level := range campusLevelInfo {
		levelInfo_ = append(levelInfo_, schoolLevelModels.LevelDataInfo{
			Name:   level.CurriculumLevel.Level.Name,
			Code:   level.Code,
			Number: level.Number,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levelInfo_})
}

func DeleteCampus(context *gin.Context) {

	type CampusInfo struct {
		CampusNumber string `json:"campus_number"`
	}

	var campusInfo CampusInfo
	err := context.BindJSON(&campusInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	campusNumber := campusInfo.CampusNumber

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

	var campusChecker models3.CampusClass
	if err := db.
		Where("campus_number = ?", campusNumber).
		First(&campusChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Campus, it has Children"})
		return
	}

	var campusData campusModel.Campus

	campusDataInfo := config.DB.
		Where("number = ?", campusNumber).
		First(&campusData)
	if campusDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "School not found!"})
		return
	}

	//delete Campus
	var levelData schoolLevelModels.Level
	db.Unscoped().Where("campus_number = ?", campusNumber).Delete(&levelData)

	//delete Campus
	db.Where("number = ?", campusNumber).Delete(&campusData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus deleted successfully!"})

}
