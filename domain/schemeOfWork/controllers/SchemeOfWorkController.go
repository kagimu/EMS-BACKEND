package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/responses"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func NewSchemeOfWork(context *gin.Context) {

	type SubtopicListInfo struct {
		Number        string `json:"number"`
		Name          string `json:"name"`
		Code          int64  `json:"code"`
		SortCode      int64  `json:"sort_code"`
		UserNumber    string `json:"user_number"`
		EstimatedTime int64  `json:"estimated_time"`
		TopicNumber   string `json:"topic_number"`
	}

	type SchemeOfWorkInfo struct {
		Name               string             `json:"name"`
		ClassNumber        string             `json:"class_number"`
		SubjectNumber      string             `json:"subject_number"`
		AcademicYearNumber string             `json:"academic_year_number"`
		SubtopicList       []SubtopicListInfo `json:"subtopicList"`
	}

	var schemeOfWorkInfo SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	name := schemeOfWorkInfo.Name
	classNumber := schemeOfWorkInfo.ClassNumber
	subjectNumber := schemeOfWorkInfo.SubjectNumber
	academicYearNumber := schemeOfWorkInfo.AcademicYearNumber
	subtopicList := schemeOfWorkInfo.SubtopicList
	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	if len(subtopicList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic list is required!"})
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

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Where("academic_year_number = ?", academicYearNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	number := uuid.NewV4().String()
	///save scheme of work data
	saveSchemeOfWork := schemeOfWorkModel.SchemeOfWork{
		Number:             number,
		Name:               name,
		Status:             "Pending",
		CampusClassNumber:  classNumber,
		SubjectNumber:      subjectNumber,
		AcademicYearNumber: academicYearNumber,
		UserNumber:         userNumber,
	}

	if err := tx.Create(&saveSchemeOfWork).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	for _, schemeList := range subtopicList {

		var duplicateSchemeCode schemeOfWorkModel.SchemeOfWorkSubTopic
		if duplicate := tx.
			Where("code = ?", schemeList.SortCode).
			Where("scheme_of_work_number = ?", saveSchemeOfWork.Number).
			First(&duplicateSchemeCode).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(schemeList.SortCode, 10)})
			return
		}

		schemeNo := uuid.NewV4().String()
		saveSchemeOfWorkSubtopic := schemeOfWorkModel.SchemeOfWorkSubTopic{
			Number:             schemeNo,
			Code:               schemeList.SortCode,
			EstimatedTime:      schemeList.EstimatedTime,
			SchemeOfWorkNumber: saveSchemeOfWork.Number,
			SubTopicNumber:     schemeList.Number,
			UserNumber:         userNumber,
		}
		if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
			return
		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme Added Successfully!", "data": saveSchemeOfWork})
}

func UpdateSchemeOfWork(context *gin.Context) {

	type SubtopicListInfo struct {
		Code               int    `json:"code"`
		SortCode           int    `json:"sort_code"`
		EstimatedTime      int    `json:"estimated_time"`
		Name               string `json:"name"`
		Number             string `json:"number"`
		UserNumber         string `json:"user_number"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubTopicNumber     string `json:"sub_topic_number"`
	}

	type SchemeOfWorkInfo struct {
		Name          string             `json:"name"`
		ClassNumber   string             `json:"class_number"`
		SubjectNumber string             `json:"subject_number"`
		SchemeNumber  string             `json:"scheme_number"`
		SubtopicList  []SubtopicListInfo `json:"subtopicList"`
	}

	var schemeOfWorkInfo SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	name := schemeOfWorkInfo.Name
	schemeNumber := schemeOfWorkInfo.SchemeNumber
	classNumber := schemeOfWorkInfo.ClassNumber
	subjectNumber := schemeOfWorkInfo.SubjectNumber
	subtopicList := schemeOfWorkInfo.SubtopicList
	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if schemeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme number is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	if len(subtopicList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic list is required!"})
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

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("number != ?", schemeNumber).
		Where("name = ?", name).
		Where("class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	tx := config.DB.Begin()
	var schemeUpdate schemeOfWorkModel.SchemeOfWork
	updateScheme := tx.Model(&schemeUpdate).
		Where("number = ?", schemeNumber).
		Updates(schemeOfWorkModel.SchemeOfWork{Name: name})
	if updateScheme.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	for _, schemeList := range subtopicList {
		codeNo, err := strconv.ParseInt(strconv.Itoa(schemeList.SortCode), 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var duplicateSchemeCode schemeOfWorkModel.SchemeOfWorkSubTopic
		if duplicate := tx.
			Where("code = ?", codeNo).
			Where("scheme_of_work_number = ?", schemeNumber).
			Where("number != ?", schemeList.Number).
			First(&duplicateSchemeCode).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(codeNo, 10)})
			return
		}

		var schemeSubtopicUpdate schemeOfWorkModel.SchemeOfWorkSubTopic
		updateScheme := tx.Model(&schemeSubtopicUpdate).
			Where("scheme_of_work_number = ?", schemeNumber).
			Where("number = ?", schemeList.Number).
			Updates(schemeOfWorkModel.SchemeOfWorkSubTopic{Code: codeNo})
		if updateScheme.Error != nil {
			tx.Rollback()
			context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to commit!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subtopicList})
}

func FetchSchemeOfWork(context *gin.Context) {

	academicYearNumber := context.Query("academic_year_number")
	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var schemeOfWrkInfo schemeOfWorkModel.SchemeOfWork
	if err := config.DB.
		Where("academic_year_number = ?", academicYearNumber).
		Where("campus_class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Preload("AcademicYear").
		Preload("CampusClass.CurriculumClass.Class").
		Preload("CampusClass.CurriculumClass").
		Preload("CampusClass").
		Preload("Subject").
		Order("created_at desc").
		First(&schemeOfWrkInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	var schemeOfWrkSubTopicInfo []schemeOfWorkModel.SchemeOfWorkSubTopic
	var schemeOfWrkSubTopicInfo_ []responses2.SchemeOfWorkSubtopicInfo
	if err := config.DB.
		Where("scheme_of_work_number = ?", schemeOfWrkInfo.Number).
		Preload("SubTopic.Topic.Strand").
		Preload("SubTopic.Topic").
		Preload("SubTopic").
		Order("code asc").
		Find(&schemeOfWrkSubTopicInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	for _, subTopicList := range schemeOfWrkSubTopicInfo {

		//fetch lesson plan
		/*	var lessonPlanInfo []models2.LessonPlan
			var lessonPlanInfo_ []responses.LessonPlanResponseInfo
			if err := config.DB.
				Where("scheme_sub_topic_number = ?", subTopicList.Number).
				//Preload("Concept").
				Order("created_at desc").
				Find(&lessonPlanInfo).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
				return
			}*/

		/*for _, lessonPlanList := range lessonPlanInfo {

			//fetch lesson plan steps
			var lessonPlanStepsInfo []models2.LessonPlanSteps
			var lessonPlanStepsInfo_ []responses.LessonPlanStepsInfo
			//var lessonPlanInfo_ []responses.LessonPlanResponseInfo
			if err := config.DB.
				Where("lesson_plan_number = ?", lessonPlanList.Number).
				Order("code asc").
				Find(&lessonPlanStepsInfo).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
				return
			}

			for _, lessonPlanSteps := range lessonPlanStepsInfo {
				lessonPlanStepsInfo_ = append(lessonPlanStepsInfo_, responses.LessonPlanStepsInfo{
					Number:           lessonPlanSteps.Number,
					Code:             lessonPlanSteps.Code,
					Step:             lessonPlanSteps.Step,
					LessonPlanNumber: lessonPlanSteps.LessonPlanNumber,
				})
			}

			lessonPlanInfo_ = append(lessonPlanInfo_, responses.LessonPlanResponseInfo{
				Number: lessonPlanList.Number,
				//LessonPlanConceptNumber: lessonPlanList.LessonPlanConceptNumber,
				SchemeSubTopicNumber: lessonPlanList.SchemeSubTopicNumber,
				//ConceptName:             lessonPlanList.Concept.Name,
				//Code:                    lessonPlanList.Code,
				UserNumber:      lessonPlanList.UserNumber,
				LessonPlanSteps: lessonPlanStepsInfo_,
			})
		}*/

		schemeOfWrkSubTopicInfo_ = append(schemeOfWrkSubTopicInfo_, responses2.SchemeOfWorkSubtopicInfo{
			Number:             subTopicList.Number,
			Code:               subTopicList.Code,
			SortCode:           subTopicList.Code,
			Name:               subTopicList.SubTopic.Name,
			TopicName:          subTopicList.SubTopic.Topic.Name,
			StrandName:         subTopicList.SubTopic.Topic.Strand.Name,
			EstimatedTime:      subTopicList.EstimatedTime,
			SchemeOfWorkNumber: subTopicList.SchemeOfWorkNumber,
			SubTopicNumber:     subTopicList.SubTopicNumber,
			UserNumber:         subTopicList.UserNumber,
			//LessonPlan:         lessonPlanInfo_,
		})
	}

	schemeOfWrk_ := responses2.SchemeOfWorkResponseInfo{
		Number:                   schemeOfWrkInfo.Number,
		Name:                     schemeOfWrkInfo.Name,
		Status:                   schemeOfWrkInfo.Status,
		ClassNumber:              schemeOfWrkInfo.CampusClassNumber,
		ClassName:                schemeOfWrkInfo.CampusClass.CurriculumClass.Class.Name,
		SubjectNumber:            schemeOfWrkInfo.SubjectNumber,
		SubjectName:              schemeOfWrkInfo.Subject.Name,
		UserNumber:               schemeOfWrkInfo.UserNumber,
		AcademicYearNumber:       schemeOfWrkInfo.AcademicYearNumber,
		AcademicYearName:         schemeOfWrkInfo.AcademicYear.AcademicYearName,
		SchemeOfWorkSubtopicInfo: schemeOfWrkSubTopicInfo_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schemeOfWrk_})

}

func DownloadSchemeOfWork(context *gin.Context) {

	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var campusClassInfo models3.CampusClass
	if err := config.DB.
		Where("number = ?", classNumber).
		Preload("CurriculumClass.Class").
		Preload("CurriculumClass").
		Preload("Streams").
		Order("code asc").
		First(&campusClassInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	var subtopic []models2.SubTopic
	var subTopicInfo_ []responses.SubTopicList
	if err := config.DB.
		Scopes(classFilter(campusClassInfo.CurriculumClassNumber)).
		Joins("join topics on topics.number = sub_topics.topic_number").
		Preload("Topic.Strand.Subject").
		Preload("Topic.Strand").
		Preload("Topic").
		Order("code asc").
		Find(&subtopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!"})
		return
	}

	for idx, subtopicList := range subtopic {
		if subjectNumber == subtopicList.Topic.Strand.Subject.Number {
			subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
				Number:        subtopicList.Number,
				Code:          subtopicList.Code,
				SortCode:      int64(idx + 1),
				Name:          subtopicList.Name,
				UserNumber:    subtopicList.UserNumber,
				EstimatedTime: subtopicList.EstimatedTime,
				TopicNumber:   subtopicList.TopicNumber,
				TopicName:     subtopicList.Topic.Name,
				StrandNumber:  subtopicList.Topic.Strand.Number,
				StrandName:    subtopicList.Topic.Strand.Name,
			})
		}
	}
	templatePath, err := generateSchemeTemplate(subTopicInfo_)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Error saving Excel file:" + err.Error()})
		return
	}

	//path := "https://" + context.Request.Host + templatePath
	path := "https://" + context.Request.Host + templatePath

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": path})

}

func generateSchemeTemplate(dataInfo []responses.SubTopicList) (string, error) {

	// Create a new Excel file
	xlsx := excelize.NewFile()
	var err error

	// Set the sheet name
	sheetName := "Sheet1"

	err11 := xlsx.SetCellValue(sheetName, columnIndexToLetters(1)+fmt.Sprint(1), "Subtopic")
	if err11 != nil {
		return "", err
	}

	err1 := xlsx.SetCellValue(sheetName, columnIndexToLetters(2)+fmt.Sprint(1), "sorting Code")
	if err1 != nil {
		return "", err1
	}

	// Iterate through the JSON data and write to the Excel file
	row := 2
	col := 1

	for _, lst := range dataInfo {
		// Write the key (JSON field) to the first column
		err := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.Name)
		if err != nil {
			return "", err
		}
		col++

		// Write the value to the second column
		err1 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.SortCode)
		if err1 != nil {
			return "", err1
		}
		col++

		// Move to the next row
		row++
		col = 1
	}

	// Save the Excel file
	//outputFilePath := "output.xlsx"
	path := "/fileRepo/excelFile/" + dataInfo[0].Number + ".xlsx"
	err = xlsx.SaveAs(path)
	if err != nil {
		//fmt.Println("Error saving Excel file:", err)
		return "", err
	}
	return path, err
	//fmt.Println("Excel file generated successfully:", path)

}

func columnIndexToLetters(index int) string {
	// Convert column index to Excel-style column letters
	letters := ""
	for index > 0 {
		index--
		letters = string(rune('A'+index%26)) + letters
		index /= 26
	}
	return letters
}

func classFilter(classNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if classNumber != "" {
			return db.Where("topics.curriculum_class_number = ?", classNumber)
		}
		return db
	}
}

func ImportSchemeOfWork(context *gin.Context) {

	userNumber := context.Request.Header["User-Number"][0]

	excelFile, fileError := context.FormFile("scheme_excel_file")

	schemeName := context.PostForm("name")
	classNumber := context.PostForm("class_number")
	subjectNumber := context.PostForm("subject_number")
	academicYearNumber := context.PostForm("academic_year_number")

	//validate if file is not uploaded
	if fileError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Excel file is missing " + fileError.Error()})
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

	if schemeName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme title is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	//	validate user number
	var user models.User
	if validateUserNumber := config.DB.Where("number = ?", userNumber).First(&user).Error; validateUserNumber != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User Not found!"})
		return
	}

	//first generate scheme of work

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("name = ?", schemeName).
		Where("class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Where("academic_year_number = ?", academicYearNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save scheme of work data
	saveSchemeOfWork := schemeOfWorkModel.SchemeOfWork{
		Number:             uuid.NewV4().String(),
		Name:               schemeName,
		Status:             "Pending",
		CampusClassNumber:  classNumber,
		SubjectNumber:      subjectNumber,
		AcademicYearNumber: academicYearNumber,
		UserNumber:         userNumber,
	}

	if err := tx.Create(&saveSchemeOfWork).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme " + err.Error()})
		return
	}

	//handle file upload
	t := time.Now()
	name := uuid.NewV4().String() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	excelName := excelFile.Filename
	extension := filepath.Ext(excelFile.Filename)
	fileName := name + extension
	path := "fileRepo/schemeOfWork/" + fileName

	err := context.SaveUploadedFile(excelFile, path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + "Failed to save File"})
		return
	}

	schemeNumber := saveSchemeOfWork.Number

	saveFileData := schemeOfWorkModel.SchemeOfWorkImport{
		Number:             uuid.NewV4().String(),
		Name:               excelName,
		FileName:           fileName,
		SchemeOfWorkNumber: schemeNumber,
		UserNumber:         userNumber,
	}

	if err := tx.Create(&saveFileData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save file details"})
		return
	}

	//read now the Excel file
	readSchemeExcelFile(path, userNumber, schemeNumber, classNumber, subjectNumber, tx, context)

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme uploaded Successfully!"})

}

func readSchemeExcelFile(filePath string, userNumber string, schemeNumber string, classNumber string, subjectNumber string, tx *gorm.DB, context *gin.Context) {

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

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx, row := range rows {
		if idx > 0 {
			schemeInfo := requests.ExcelSchemeOfWorkInfo{
				SubTopic:    getValueAtIndex(row, 0),
				SortingCode: getValueAtIndex(row, 1),
			}
			saveSchemeOfWorkData(schemeInfo, idx, userNumber, schemeNumber, classNumber, subjectNumber, tx, context)
		}
	}

}

func saveSchemeOfWorkData(schemeInfo requests.ExcelSchemeOfWorkInfo, idx int, userNumber string, schemeNumber string, classNumber string, subjectNumber string, tx *gorm.DB, context *gin.Context) {

	codeNo, err := strconv.ParseInt(schemeInfo.SortingCode, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var duplicateSchemeCode schemeOfWorkModel.SchemeOfWorkSubTopic
	if duplicate := tx.
		Where("code = ?", codeNo).
		Where("scheme_of_work_number = ?", schemeNumber).
		First(&duplicateSchemeCode).Error; duplicate == nil {
		tx.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(codeNo, 10)})
		return
	}

	//fetch scheme of work
	//----------------
	//first identify which curriculum is on this subject
	var subjectInfo subjectModel.Subject
	if err := tx.
		Where("number = ?", subjectNumber).
		First(&subjectInfo).Error; err != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject!"})
		return
	}

	//then fetch subtopic by curriculum and its  name
	var subtopic models2.SubTopic
	if err := tx.
		Where("name = ?", schemeInfo.SubTopic).
		Where("curriculum_number = ?", subjectInfo.CurriculumNumber).
		First(&subtopic).Error; err != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!!"})
		return
	}

	//-----------------------
	schemeNo := uuid.NewV4().String()
	saveSchemeOfWorkSubtopic := schemeOfWorkModel.SchemeOfWorkSubTopic{
		Number:             schemeNo,
		Code:               codeNo,
		SubTopicNumber:     subtopic.Number,
		EstimatedTime:      subtopic.EstimatedTime,
		SchemeOfWorkNumber: schemeNumber,
		UserNumber:         userNumber,
	}
	if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme of work sub topic"})
		return
	}
	//}

}

func getValueAtIndex(row []string, i int) string {
	if len(row) > i {
		return row[i]
	} else {
		return ""
	}
}
