package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/api/requests"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/api/responses"
	models5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	strandModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	schoolCurriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/responses"
	curriculumLevelModels2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func NewCurriculum(context *gin.Context) {

	var curriculumInfo requests.CurriculumCreateRequest
	err := context.BindJSON(&curriculumInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	curriculumInfo.Prepare()

	if err := curriculumInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	name := curriculumInfo.Name

	userNumber := context.Request.Header["User-Number"][0]
	entityNumber, _ := commons.GetUserEntityNumber(context)
	var userData models.User
	//check if user is a valid user
	if err := config.DB.Where("number = ?", userNumber).First(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user data", "error": err.Error()})
		}
		return
	}

	var duplicateSchool models5.Curriculum
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("entity_number = ?", entityNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum with this name  already exists !!!"})
		return
	}

	var dupCode models5.Curriculum
	if duplicate := config.DB.
		Where("code = ?", curriculumInfo.Code).
		Where("entity_number = ?", entityNumber).
		First(&dupCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum with this code  already exists !!!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save user data

	saveCurriculum := models5.Curriculum{
		Number:       uuid.NewV4().String(),
		Name:         name,
		Code:         curriculumInfo.Code,
		UserNumber:   userNumber,
		EntityNumber: entityNumber,
	}

	if err := tx.Create(&saveCurriculum).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Curriculum"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Curriculum"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Curriculum Added Successfully!", "data": saveCurriculum})
}

func ImportCurriculum(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	//validate maximum file size
	excelFile, fileError := context.FormFile("excel_file")

	curriculumNumber := context.PostForm("curriculum_number")

	//validate if file is not uploaded
	if fileError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Excel file is missing"})
		return
	}

	// upload size must not exceed 20 megabytes
	if excelFile.Size > 20971520 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Pdf file must not exceed 20 mbs!"})
		return
	}

	if excelFile.Size < 1 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Pdf"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	//handle file upload
	t := time.Now()
	name := uuid.NewV4().String() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	excelName := excelFile.Filename
	extension := filepath.Ext(excelFile.Filename)
	fileName := name + extension
	path := "fileRepo/excelFile/" + fileName

	err = context.SaveUploadedFile(excelFile, path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + "Failed to save File"})
		return
	}

	// move the assessment grade to another table
	// normalize
	fileNumber := uuid.NewV4().String()
	userNumber := user.Number
	saveFileData := models5.ImportCurriculum{
		Number:           fileNumber,
		Name:             excelName,
		FileName:         fileName,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}
	config.DB.Save(&saveFileData)

	//read now the Excel file
	readExcelFile(path, userNumber, curriculumNumber)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Excel File uploaded Successfully!", "data": saveFileData})
}

func readExcelFile(filePath string, userNumber string, curriculumNumber string) {
	//f, err := excelize.OpenFile("Book1.xlsx")
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//get sheet name dynamically
	sheetName := f.GetSheetList()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheetName[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx, row := range rows {
		if idx > 0 {
			curriculumInfo := models5.ExcelCurriculumInfo{
				LearningAreaCode:     getValueAtIndex(row, 0),
				LearningArea:         getValueAtIndex(row, 1),
				SubjectClusterCode:   getValueAtIndex(row, 2),
				SubjectCluster:       getValueAtIndex(row, 3),
				SubjectCode:          getValueAtIndex(row, 4),
				Subject:              getValueAtIndex(row, 5),
				Class:                getValueAtIndex(row, 6),
				StrandCode:           getValueAtIndex(row, 7),
				Strand:               getValueAtIndex(row, 8),
				TopicCode:            getValueAtIndex(row, 9),
				Topic:                getValueAtIndex(row, 10),
				SubTopicCode:         getValueAtIndex(row, 11),
				SubTopic:             getValueAtIndex(row, 12),
				TeachingSequenceCode: getValueAtIndex(row, 13),
				ConceptCode:          getValueAtIndex(row, 14),
				Concept:              getValueAtIndex(row, 15),
				EstimatedTime:        getValueAtIndex(row, 16),
				RequiredMaterials:    getValueAtIndex(row, 17),
			}
			curriculumInfo.Prepare()
			saveCurriculumData(curriculumInfo, idx, userNumber, curriculumNumber)
		}
	}

}

func saveCurriculumData(info models5.ExcelCurriculumInfo, idx int, userNumber string, curriculumNumber string) {

	//----------------------------------------------
	//	save learning area

	if info.LearningArea == "" {
		return
	}

	learningAreaNo := uuid.NewV4().String()

	learningAreaCode, err := strconv.ParseInt(info.LearningAreaCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	///save user data
	saveLearningArea := learningAreaModels.LearningArea{
		Number:           learningAreaNo,
		Name:             info.LearningArea,
		Code:             learningAreaCode,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}

	//check if learning area exists
	//if it does not exist add it
	var learningAreaInfo learningAreaModels.LearningArea
	learningAreaDup := config.DB.
		Where("name = ?", info.LearningArea).
		Where("curriculum_number = ?", curriculumNumber).
		First(&learningAreaInfo)
	if learningAreaDup.Error != nil {
		err := config.DB.Save(&saveLearningArea).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	learningAreaNumber := ""

	if learningAreaInfo.Number != "" {
		learningAreaNumber = learningAreaInfo.Number
	} else {
		learningAreaNumber = saveLearningArea.Number
	}

	//-------------------------------------------

	//----------------------------------------
	//	adds subject cluster

	if info.SubjectCluster == "" {
		return
	}

	subjectClusterCode, err := strconv.ParseInt(info.SubjectClusterCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	subjectClusterNo := uuid.NewV4().String()
	///save user data
	saveSubjectCluster := models3.SubjectCluster{
		Number:             subjectClusterNo,
		Name:               info.SubjectCluster,
		Code:               subjectClusterCode,
		LearningAreaNumber: learningAreaNumber,
		CurriculumNumber:   curriculumNumber,
		UserNumber:         userNumber,
	}

	var subjectClusterInfo models3.SubjectCluster
	subjectClusterDup := config.DB.
		Where("name = ?", info.SubjectCluster).
		Where("learning_area_number = ?", learningAreaInfo.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&subjectClusterInfo)
	if subjectClusterDup.Error != nil {
		err := config.DB.Save(&saveSubjectCluster).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	//--------------------------------
	//	adds subject

	subjectClusterNumber := ""

	if subjectClusterInfo.Number != "" {
		subjectClusterNumber = subjectClusterInfo.Number
	} else {
		subjectClusterNumber = saveSubjectCluster.Number
	}

	if info.Subject == "" {
		return
	}

	subjectCode, err := strconv.ParseInt(info.SubjectCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	subjectNo := uuid.NewV4().String()
	///save user data
	saveSubject := models4.Subject{
		Number:               subjectNo,
		Name:                 info.Subject,
		Code:                 subjectCode,
		UserNumber:           userNumber,
		CurriculumNumber:     curriculumNumber,
		SubjectClusterNumber: subjectClusterNumber,
	}

	var subjectInfo models4.Subject
	subjectDup := config.DB.
		Where("name = ?", info.Subject).
		Where("subject_cluster_number = ?", subjectClusterInfo.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&subjectInfo)
	if subjectDup.Error != nil {
		err := config.DB.Save(&saveSubject).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}
	//------------------------------

	//	adds Strand

	subjectNumber := ""

	if subjectInfo.Number != "" {
		subjectNumber = subjectInfo.Number
	} else {
		subjectNumber = saveSubject.Number
	}

	if info.Strand == "" {
		return
	}

	strandCode, err := strconv.ParseInt(info.StrandCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	strandNo := uuid.NewV4().String()
	///save user data
	saveStrand := strandModels.Strand{
		Number:           strandNo,
		Name:             info.Strand,
		Code:             strandCode,
		SubjectNumber:    subjectNumber,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}

	var strandInfo strandModels.Strand
	strandDup := config.DB.
		Where("name = ?", info.Strand).
		Where("subject_number = ?", subjectInfo.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&strandInfo)
	if strandDup.Error != nil {
		err := config.DB.Save(&saveStrand).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	//	---------------
	// adds Class

	if info.Class == "" {
		return
	}

	var classes classModels.CurriculumClass
	if err := config.DB.
		Where("classes.abbreviation = ?", info.Class).
		Where("curriculum_levels.curriculum_number = ?", curriculumNumber).
		Joins("join curriculum_levels on curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Joins("join classes on classes.number = curriculum_classes.class_number").
		Preload("CurriculumLevel").
		Preload("Class").
		First(&classes).Error; err != nil {
		//context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	//---------------------

	//Adds a Topic

	strandNumber := ""

	if strandInfo.Number != "" {
		strandNumber = strandInfo.Number
	} else {
		strandNumber = saveStrand.Number
	}

	if info.Topic == "" {
		return
	}

	topicCode, err := strconv.ParseInt(info.TopicCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	topicNo := uuid.NewV4().String()
	///save user data
	saveTopic := topicModels.Topic{
		Number:                topicNo,
		Name:                  info.Topic,
		Code:                  topicCode,
		StrandNumber:          strandNumber,
		UserNumber:            userNumber,
		CurriculumNumber:      curriculumNumber,
		CurriculumClassNumber: classes.Number,
	}

	var topicInfo topicModels.Topic
	topicDup := config.DB.
		Where("name = ?", info.Topic).
		Where("strand_number = ?", strandInfo.Number).
		Where("curriculum_class_number = ?", classes.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&topicInfo)
	if topicDup.Error != nil {
		err := config.DB.Save(&saveTopic).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	//------------------

	//	adds sub topic

	topicNumber := ""
	if topicInfo.Number != "" {
		topicNumber = topicInfo.Number
	} else {
		topicNumber = saveTopic.Number
	}

	if info.SubTopic == "" {
		return
	}
	subtopicNo := uuid.NewV4().String()

	estimatedTimeData := info.EstimatedTime

	if estimatedTimeData == "" {
		estimatedTimeData = "0"
	}

	estTime, _ := strconv.ParseInt(estimatedTimeData, 10, 64)

	subTopicCode, _ := strconv.ParseInt(info.SubTopicCode, 10, 64)

	subTopicTeachingSequenceCode, _ := strconv.ParseInt(info.TeachingSequenceCode, 10, 64)

	fmt.Println("Teaching sequence code " + info.TeachingSequenceCode)

	///save user data
	saveSubTopic := subTopicModels.SubTopic{
		Number:              subtopicNo,
		Name:                info.SubTopic,
		Code:                subTopicCode,
		TopicNumber:         topicNumber,
		CurriculumNumber:    curriculumNumber,
		UserNumber:          userNumber,
		TeacherSequenceCode: subTopicTeachingSequenceCode,
	}

	var subtopicInfo subTopicModels.SubTopic
	subtopicDup := config.DB.
		Where("name = ?", info.SubTopic).
		Where("topic_number = ?", topicInfo.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&subtopicInfo)
	if subtopicDup.Error != nil {
		err := config.DB.Save(&saveSubTopic).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	//-------------
	//	adds concept

	subTopicNumber := ""
	if subtopicInfo.Number != "" {
		subTopicNumber = subtopicInfo.Number
	} else {
		subTopicNumber = saveSubTopic.Number
	}

	if info.Concept == "" {
		return
	}

	conceptCode, err := strconv.ParseInt(info.ConceptCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	conceptNo := uuid.NewV4().String()
	///save user data
	saveConcept := conceptModels.Concept{
		Number:            conceptNo,
		Name:              info.Concept,
		Code:              conceptCode,
		SubTopicNumber:    subTopicNumber,
		CurriculumNumber:  curriculumNumber,
		UserNumber:        userNumber,
		EstimatedTime:     estTime,
		RequiredMaterials: info.RequiredMaterials,
	}

	var conceptInfo conceptModels.Concept
	conceptDup := config.DB.
		Where("name = ?", info.Concept).
		Where("sub_topic_number = ?", subtopicInfo.Number).
		Where("curriculum_number = ?", curriculumNumber).
		First(&conceptInfo)
	if conceptDup.Error != nil {
		err := config.DB.Save(&saveConcept).Error
		if err != nil {
			fmt.Println("Insertion failed:", err)
			return
		}
	}

	//--------------

}

func getValueAtIndex(row []string, i int) string {
	if len(row) > i {
		return row[i]
	} else {
		return ""
	}
}

func FetchCurriculum(context *gin.Context) {

	var curriculum []models5.Curriculum
	if err := config.DB.
		Scopes(filters.ByEntityNumber(context.Query("entity_number"))).
		Order("code asc").
		Find(&curriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculum})
}

func FetchCurriculumInfo(context *gin.Context) {
	curriculumNumber := context.Query("curriculum_number")
	var curriculum models5.Curriculum
	if err := config.DB.
		Where("number = ?", curriculumNumber).
		First(&curriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculum})
}

func FetchUploadedCurriculum(context *gin.Context) {
	curriculumNumber := context.Query("curriculum_number")
	var curriculumImport []models5.ImportCurriculum
	if err := config.DB.
		Where("curriculum_number = ?", curriculumNumber).
		Order("created_at desc").
		Find(&curriculumImport).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculumImport})
}

func FetchSchoolCurriculum(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	var schoolCurriculum []schoolCurriculumModel.SchoolCurriculum
	var schoolCurriculum_ []responses.SchoolCurriculumDataList
	if err := config.DB.
		Where("school_number = ?", schoolNumber).
		Preload("Curriculum").
		Find(&schoolCurriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	for _, item := range schoolCurriculum {
		schoolCurriculum_ = append(schoolCurriculum_, responses.SchoolCurriculumDataList{
			Name:   item.Curriculum.Name,
			Number: item.Number,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schoolCurriculum_})
}

func FetchCampusCurriculumLevel(context *gin.Context) {
	campusCurriculumNumber := context.Query("campus_curriculum_number")

	var campusCurriculum campusModel.CampusCurriculum
	if err := config.DB.
		Where("number = ?", campusCurriculumNumber).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Find(&campusCurriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	var campusCurriculumLevel []curriculumLevelModels2.CurriculumLevel
	var campusCurriculumLevel_ []responses2.CurriculumLevelList
	if err := config.DB.
		Where("curriculum_number = ?", campusCurriculum.SchoolCurriculum.CurriculumNumber).
		Preload("Level").
		Preload("Curriculum").
		Order("code asc").
		Find(&campusCurriculumLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	for _, item := range campusCurriculumLevel {
		campusCurriculumLevel_ = append(campusCurriculumLevel_, responses2.CurriculumLevelList{
			Number:           item.Number,
			Name:             item.Level.Name,
			Code:             item.Code,
			UserNumber:       item.UserNumber,
			LevelNumber:      item.LevelNumber,
			CurriculumName:   item.Curriculum.Name,
			CurriculumNumber: item.CurriculumNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campusCurriculumLevel_})
}

func DeleteCurriculum(context *gin.Context) {

	var curriculumInfo requests.CurriculumDelRequest
	err := context.BindJSON(&curriculumInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	curriculumNumber := curriculumInfo.Number

	var curriculumChecker curriculumLevelModels2.CurriculumLevel
	if err := config.DB.
		Where("curriculum_number = ?", curriculumNumber).
		First(&curriculumChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this curriculum, it has level under it."})
		return
	}

	var curriculumChecker2 learningAreaModels.LearningArea
	if err := config.DB.
		Where("curriculum_number = ?", curriculumNumber).
		First(&curriculumChecker2).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this curriculum, it has Learning Area under it."})
		return
	}

	var curriculumData models5.Curriculum
	curriculumDataInfo := config.DB.
		Where("number = ?", curriculumNumber).
		First(&curriculumData)
	if curriculumDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Curriculum not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", curriculumNumber).
		Delete(&curriculumData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete curriculum" + err.Error()})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete curriculum"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Curriculum deleted successfully!"})
}

func CurriculumInfo(context *gin.Context) {

	curriculumNumber := context.Query("curriculum_number")

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum  is required!"})
		return
	}

	var curriculum models5.Curriculum
	if err := config.DB.
		Where("number = ?", curriculumNumber).
		First(&curriculum).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Curriculum found"})
		return
	}

	curriculum_ := responses.SchoolCurriculumDataList{
		Name:   curriculum.Name,
		Number: curriculum.Number,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculum_})
}

func CurriculumEdit(context *gin.Context) {

	var curriculumInfo responses.CurriculumResponse
	err := context.BindJSON(&curriculumInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	curriculumInfo.Prepare()

	if err := curriculumInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	name := curriculumInfo.Name
	curriculumNumber := curriculumInfo.CurriculumNumber

	userNumber := context.Request.Header["User-Number"][0]

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	var userData models.User
	//check if user is a valid user
	if err := config.DB.Where("number = ?", userNumber).First(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user data", "error": err.Error()})
		}
		return
	}

	var duplicateSchool models5.Curriculum
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		Where("number != ?", curriculumNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum already exists!!!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	var curriculumDataInfo models5.Curriculum
	updateCurriculum := tx.Model(&curriculumDataInfo).
		Where("number = ?", curriculumNumber).
		Updates(models5.Curriculum{Name: name, UserNumber: userNumber, Code: curriculumInfo.Code})
	if updateCurriculum.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Curriculum"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Curriculum Edited Successfully!", "data": curriculumDataInfo})
}

func ExportCurriculum(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")
	strandNumber := context.Query("strand_number")
	topicNumber := context.Query("topic_number")
	subtopicNumber := context.Query("sub_topic_number")
	curriculumNumber := context.Query("curriculum_number")
	classNumber := context.Query("class_number")

	//fetch learning area
	learningAreaResp, msg := curriculum.FetchLearningAreaCurriculum(learningAreaNumber, curriculumNumber)
	if msg != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": msg})
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Set value of a cell.
	_ = f.SetCellValue("Sheet2", "A2", "Learning Area.")
	_ = f.SetCellValue("Sheet2", "A3", "Subject Cluster Code.")
	_ = f.SetCellValue("Sheet2", "A4", "Subject")

	_ = f.SetCellValue("Sheet2", "B2", " : ")
	_ = f.SetCellValue("Sheet2", "B3", " : ")
	_ = f.SetCellValue("Sheet2", "B4", " : ")

	// Add header row for detailed curriculum data starting at A5
	_ = f.SetCellValue("Sheet2", "A5", "Strand Code")
	_ = f.SetCellValue("Sheet2", "B5", "Strand Name")
	_ = f.SetCellValue("Sheet2", "C5", "Topic Code")
	_ = f.SetCellValue("Sheet2", "D5", "Topic Name")
	_ = f.SetCellValue("Sheet2", "E5", "Subtopic Code")
	_ = f.SetCellValue("Sheet2", "F5", "Subtopic Name (Lesson)")
	_ = f.SetCellValue("Sheet2", "G5", "Teaching Sequence Code")
	_ = f.SetCellValue("Sheet2", "H5", "Concept Code")
	_ = f.SetCellValue("Sheet2", "I5", "Concept Name")
	_ = f.SetCellValue("Sheet2", "J5", "Estimated Time Load (Minutes)")
	_ = f.SetCellValue("Sheet2", "K5", "Required Materials")

	lstLA := learningAreaResp[0]

	//fetch subject cluster
	subjectClusterResp, msg := curriculum.FetchSubjectClusterCurriculum("", lstLA.Number, subjectClusterNumber)
	if msg != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	lstSC := subjectClusterResp[0]

	subjectResp, msg := curriculum.FetchSubjectCurriculum("", lstSC.Number, subjectNumber)
	if msg != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}
	lstSB := subjectResp[0]

	_ = f.SetCellValue("Sheet2", "C2", lstLA.Name)
	_ = f.SetCellValue("Sheet2", "C3", lstSC.Name)
	_ = f.SetCellValue("Sheet2", "C4", lstSB.Name)

	strandResp, msg := curriculum.FetchStrandCurriculum("", lstSB.Number, strandNumber)
	if msg != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	// Initialize row counter for data rows (starting from row 6)
	rowCounter := 6

	// Variables to track previous values to avoid duplication
	var prevStrandCode, prevTopicCode, prevSubtopicCode int64

	// Loop through strands and populate Excel data
	for _, lstST := range strandResp {
		// Reset topic and subtopic tracking when starting a new strand
		// This ensures that similar codes across different strands are handled correctly
		if lstST.Code != prevStrandCode {
			prevTopicCode = 0
			prevSubtopicCode = 0
		}

		topicResp, err := curriculum.FetchTopicCurriculum("", lstST.Number, topicNumber, classNumber)
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		for _, list := range topicResp {
			// Reset subtopic tracking when starting a new topic
			// This ensures that similar subtopic codes across different topics are handled correctly
			if list.Code != prevTopicCode {
				prevSubtopicCode = 0
			}

			subtopicResp, err := curriculum.FetchSubtopicCurriculum("", list.Number, subtopicNumber)
			if err != "" {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			}

			for _, subtopicList := range subtopicResp {

				conceptResp, err := curriculum.FetchConceptCurriculum("", subtopicList.Number)
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				// If there are concepts, create a row for each concept
				if len(conceptResp) > 0 {
					for _, concept := range conceptResp {
						// Determine which values to show (only when they change)
						strandCode := ""
						strandName := ""
						topicCode := ""
						topicName := ""
						subtopicCode := ""
						subtopicName := ""

						// Show strand info only if it changed
						if lstST.Code != prevStrandCode {
							strandCode = fmt.Sprintf("%d", lstST.Code)
							strandName = lstST.Name
							prevStrandCode = lstST.Code
						}

						// Show topic info only if it changed
						if list.Code != prevTopicCode {
							topicCode = fmt.Sprintf("%d", list.Code)
							topicName = list.Name
							prevTopicCode = list.Code
						}

						// Show subtopic info only if it changed
						if subtopicList.Code != prevSubtopicCode {
							subtopicCode = fmt.Sprintf("%d", subtopicList.Code)
							subtopicName = subtopicList.Name
							prevSubtopicCode = subtopicList.Code
						}

						// Write data to Excel row
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("A%d", rowCounter), strandCode)                // Strand Code
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("B%d", rowCounter), strandName)                // Strand Name
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("C%d", rowCounter), topicCode)                 // Topic Code
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("D%d", rowCounter), topicName)                 // Topic Name
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("E%d", rowCounter), subtopicCode)              // Subtopic Code
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("F%d", rowCounter), subtopicName)              // Subtopic Name (Lesson)
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("G%d", rowCounter), "")                        // Teaching Sequence Code (empty for now)
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("H%d", rowCounter), concept.Code)              // Concept Code
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("I%d", rowCounter), concept.Name)              // Concept Name
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("J%d", rowCounter), concept.EstimatedTime)     // Estimated Time Load (Minutes)
						_ = f.SetCellValue("Sheet2", fmt.Sprintf("K%d", rowCounter), concept.RequiredMaterials) // Required materials

						rowCounter++
					}
				} else {
					// If no concepts, still create a row for the subtopic
					// Determine which values to show (only when they change)
					strandCode := ""
					strandName := ""
					topicCode := ""
					topicName := ""
					subtopicCode := ""
					subtopicName := ""

					// Show strand info only if it changed
					if lstST.Code != prevStrandCode {
						strandCode = fmt.Sprintf("%d", lstST.Code)
						strandName = lstST.Name
						prevStrandCode = lstST.Code
					}

					// Show topic info only if it changed
					if list.Code != prevTopicCode {
						topicCode = fmt.Sprintf("%d", list.Code)
						topicName = list.Name
						prevTopicCode = list.Code
					}

					// Show subtopic info only if it changed
					if subtopicList.Code != prevSubtopicCode {
						subtopicCode = fmt.Sprintf("%d", subtopicList.Code)
						subtopicName = subtopicList.Name
						prevSubtopicCode = subtopicList.Code
					}

					// Write data to Excel row
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("A%d", rowCounter), strandCode)                 // Strand Code
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("B%d", rowCounter), strandName)                 // Strand Name
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("C%d", rowCounter), topicCode)                  // Topic Code
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("D%d", rowCounter), topicName)                  // Topic Name
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("E%d", rowCounter), subtopicCode)               // Subtopic Code
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("F%d", rowCounter), subtopicName)               // Subtopic Name (Lesson)
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("G%d", rowCounter), subtopicList.EstimatedTime) // Estimated Time Load (Minutes)
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("H%d", rowCounter), "")                         // Teaching Sequence Code (empty for now)
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("I%d", rowCounter), "")                         // Concept Code (empty)
					_ = f.SetCellValue("Sheet2", fmt.Sprintf("J%d", rowCounter), "")                         // Concept Name (empty)

					rowCounter++
				}
			}
		}
	}

	f.SetActiveSheet(index)

	// Add borders to the Excel file
	// Define border style
	borderStyle := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	// Apply borders and styling to header row (row 5)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})

	for col := 'A'; col <= 'K'; col++ {
		cellRef := fmt.Sprintf("%c5", col)
		f.SetCellStyle("Sheet2", cellRef, cellRef, headerStyle)
	}

	// Apply borders to basic info rows (rows 2-4)
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Font: &excelize.Font{
			Bold: true,
		},
	})

	for row := 2; row <= 4; row++ {
		for col := 'A'; col <= 'C'; col++ {
			cellRef := fmt.Sprintf("%c%d", col, row)
			f.SetCellStyle("Sheet2", cellRef, cellRef, infoStyle)
		}
	}

	// Apply borders to data rows (from row 6 to the last row)
	for row := 6; row < rowCounter; row++ {
		for col := 'A'; col <= 'K'; col++ {
			cellRef := fmt.Sprintf("%c%d", col, row)
			style, _ := f.NewStyle(&excelize.Style{
				Border: borderStyle,
			})
			f.SetCellStyle("Sheet2", cellRef, cellRef, style)
		}
	}

	unixTime := time.Now().Unix()
	fileName := fmt.Sprintf("%s-%d.xlsx", lstSB.Name, unixTime)

	context.Header("Content-Disposition", "attachment; filename="+fileName)
	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := f.Write(context.Writer); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

}

func ExportLessonPlanTemplate(ctx *gin.Context) (*excelize.File, string, error) {

	curriculumNumber := ctx.Query("curriculum_number")

	fileTemplate := excelize.NewFile()

	schemeOfWorkSubTopic, err := repositories2.GetSchemeOfWorkSubTopicByNumber(curriculumNumber)
	if err != nil {
		return nil, "", errors.New("Checking scheme subtopic error: " + err.Error())
	}

	if schemeOfWorkSubTopic.Number == "" {
		return nil, "", errors.New("Invalid scheme subtopic, we can not export anything ")
	}

	//defaultOpenSheet := 0

	defaultOpenSheet := -1

	sheetName := fmt.Sprintf("Sheet-%d", "name")

	fmt.Println("Sheet Name: ", sheetName)

	defaultOpenSheet, err = fileTemplate.NewSheet(sheetName)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("Error '%s' creating sheet '%s'", err.Error(), " concept.Name"))
	}

	sheetsHeaders := [][]string{
		{"Concept name:", "Name"},
		{"Code", "Step"},
	}

	for i, row := range sheetsHeaders {
		for j, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
			err = fileTemplate.SetCellValue(sheetName, cellRef, cell)
			if err != nil {
				return nil, "", errors.New(fmt.Sprintf("Error '%s' writing sheet '%s' headers", err.Error(), "concept.Name"))
			}
		}
	}

	if defaultOpenSheet < 0 {
		return nil, "", errors.New(fmt.Sprintf("Invalid default open sheet number: %d", defaultOpenSheet))
	}

	fileTemplate.SetActiveSheet(defaultOpenSheet)

	unixTime := time.Now().Unix()
	return fileTemplate, fmt.Sprintf("%s-%d.xlsx", "subTopic.Name", unixTime), nil
}
