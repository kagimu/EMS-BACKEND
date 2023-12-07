package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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

	type CurriculumInfo struct {
		Name string `json:"name"`
	}

	var curriculumInfo CurriculumInfo
	err := context.BindJSON(&curriculumInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := curriculumInfo.Name

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
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
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum already exists!!!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save user data
	saveCurriculum := models5.Curriculum{
		Number:     uuid.NewV4().String(),
		Name:       name,
		UserNumber: userNumber,
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

	userNumber := context.Request.Header["User-Number"][0]

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

	if userNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "user number is required!"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	//	validate user number
	var user models.User
	if validateUserNumber := config.DB.Where("number = ?", userNumber).First(&user).Error; validateUserNumber != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User Not found!"})
		return
	}

	//handle file upload
	t := time.Now()
	name := uuid.NewV4().String() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	excelName := excelFile.Filename
	extension := filepath.Ext(excelFile.Filename)
	fileName := name + extension
	//jpgName := name + ".png"
	//path := "/fileRepo/excelFile/" + fileName
	//path := "/fileRepo/excelFile/" + fileName
	path := "fileRepo/excelFile/" + fileName
	//jpgPath := "/fileRepo/jpgFile/" + jpgName

	/*if extension != ".pdf" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Only pdf files allowed!!"})
		return
	}*/

	/*if excelFile.Header.Get("Content-Type") != "application/pdf" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type"})
		return
	}*/

	err := context.SaveUploadedFile(excelFile, path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + "Failed to save File"})
		return
	}

	//check mime type of file
	/*	mType, err2 := mimetype.DetectFile(path)
		if err2 != nil {
			fmt.Println("Error detecting file path:", err2)
		}

		if mType.String() != "application/pdf" {

			//if the file  type is invalid the pdf file should be removed immediately
			errRemove := os.Remove(path)
			if errRemove != nil {
				context.JSON(http.StatusBadRequest, gin.H{"message": errRemove})
				return
			}
			/////////////////////////////////////////////////////////////////////////

			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type " + mType.String()})
			return
		}*/

	fileNumber := uuid.NewV4().String()
	saveFileData := models5.ImportCurriculum{
		Number:           fileNumber,
		Name:             excelName,
		FileName:         fileName,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}
	config.DB.Save(&saveFileData)

	//read now the excel file
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
	// Get value from cell by given worksheet name and cell reference.
	//cell, err := f.GetCellValue("Sheet1", "B2")
	/*cell, err := f.GetCellValue("Sheet1", "A1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)*/
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx, row := range rows {
		if idx > 0 {
			curriculumInfo := models5.ExcelCurriculumInfo{
				LearningAreaCode:   getValueAtIndex(row, 0),
				LearningArea:       getValueAtIndex(row, 1),
				SubjectClusterCode: getValueAtIndex(row, 2),
				SubjectCluster:     getValueAtIndex(row, 3),
				SubjectCode:        getValueAtIndex(row, 4),
				Subject:            getValueAtIndex(row, 5),
				Class:              getValueAtIndex(row, 6),
				StrandCode:         getValueAtIndex(row, 7),
				Strand:             getValueAtIndex(row, 8),
				TopicCode:          getValueAtIndex(row, 9),
				Topic:              getValueAtIndex(row, 10),
				SubTopicCode:       getValueAtIndex(row, 11),
				SubTopic:           getValueAtIndex(row, 12),
				EstimatedTime:      getValueAtIndex(row, 13),
				ConceptCode:        getValueAtIndex(row, 14),
				Concept:            getValueAtIndex(row, 15),
			}
			saveCurriculumData(curriculumInfo, idx, userNumber, curriculumNumber)
		}
	}

	//fmt.Println(curriculumInfo)

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
		UserNumber:         userNumber,
	}

	var subjectClusterInfo models3.SubjectCluster
	subjectClusterDup := config.DB.
		Where("name = ?", info.SubjectCluster).
		Where("learning_area_number = ?", learningAreaInfo.Number).
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
		SubjectClusterNumber: subjectClusterNumber,
	}

	var subjectInfo models4.Subject
	subjectDup := config.DB.
		Where("name = ?", info.Subject).
		Where("subject_cluster_number = ?", subjectClusterInfo.Number).
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
		Number:        strandNo,
		Name:          info.Strand,
		Code:          strandCode,
		SubjectNumber: subjectNumber,
		UserNumber:    userNumber,
	}

	var strandInfo strandModels.Strand
	strandDup := config.DB.
		Where("name = ?", info.Strand).
		Where("subject_number = ?", subjectInfo.Number).
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
		CurriculumClassNumber: classes.Number,
	}

	var topicInfo topicModels.Topic
	topicDup := config.DB.
		Where("name = ?", info.Topic).
		Where("strand_number = ?", strandInfo.Number).
		Where("curriculum_class_number = ?", classes.Number).
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

	estTime, err := strconv.ParseInt(estimatedTimeData, 10, 64)

	if err != nil {
		// Handle the error if the string cannot be converted to int64
		//return err
	}

	subTopicCode, err := strconv.ParseInt(info.SubTopicCode, 10, 64)
	if err != nil {
		//fmt.Println("Error:", err)
	}

	///save user data
	saveSubTopic := subTopicModels.SubTopic{
		Number:        subtopicNo,
		Name:          info.SubTopic,
		Code:          subTopicCode,
		EstimatedTime: estTime,
		TopicNumber:   topicNumber,
		UserNumber:    userNumber,
	}

	var subtopicInfo subTopicModels.SubTopic
	subtopicDup := config.DB.
		Where("name = ?", info.SubTopic).
		Where("topic_number = ?", topicInfo.Number).
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
		Number:         conceptNo,
		Name:           info.Concept,
		Code:           conceptCode,
		SubTopicNumber: subTopicNumber,
		UserNumber:     userNumber,
	}

	var conceptInfo conceptModels.Concept
	conceptDup := config.DB.
		Where("name = ?", info.Concept).
		Where("sub_topic_number = ?", subtopicInfo.Number).
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
		Order("created_at desc").
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
		Order("created_at desc").
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

func FetchSchoolCurriculumLevel(context *gin.Context) {
	schoolCurriculumNumber := context.Query("school_curriculum_number")

	var schoolCurriculumInfo schoolCurriculumModel.SchoolCurriculum
	if err := config.DB.
		Where("number = ?", schoolCurriculumNumber).
		Order("created_at desc").
		First(&schoolCurriculumInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No School curriculum found!"})
		return
	}

	var schoolCurriculumLevel []curriculumLevelModels2.CurriculumLevel
	var schoolCurriculumLevel_ []responses2.CurriculumLevelList
	if err := config.DB.
		Where("curriculum_number = ?", schoolCurriculumInfo.CurriculumNumber).
		Preload("Level").
		Preload("Curriculum").
		Order("code asc").
		Find(&schoolCurriculumLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	for _, item := range schoolCurriculumLevel {
		schoolCurriculumLevel_ = append(schoolCurriculumLevel_, responses2.CurriculumLevelList{
			Number:           item.Number,
			Name:             item.Level.Name,
			Code:             item.Code,
			UserNumber:       item.UserNumber,
			LevelNumber:      item.LevelNumber,
			CurriculumName:   item.Curriculum.Name,
			CurriculumNumber: item.CurriculumNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schoolCurriculumLevel_})
}
