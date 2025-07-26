package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/responses"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/services"
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

var schemeOfWorkService = services.InitSchemeOfWorkService()

func NewSchemeOfWork(context *gin.Context) {

	var schemeOfWorkInfo requests.SchemeOfWorkInfo
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
		Status:             constants.StatusActive,
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
			Where("code = ?", schemeList.Code).
			Where("scheme_of_work_number = ?", saveSchemeOfWork.Number).
			Where("scheme_of_works.status = ?", constants.StatusActive).
			Joins("JOIN scheme_of_works ON scheme_of_work_number = scheme_of_works.number").
			First(&duplicateSchemeCode).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(schemeList.Code, 10)})
			return
		}

		schemeNo := uuid.NewV4().String()
		saveSchemeOfWorkSubtopic := schemeOfWorkModel.SchemeOfWorkSubTopic{
			Number:             schemeNo,
			Code:               schemeList.Code,
			EstimatedTime:      schemeList.EstimatedTime,
			SchemeOfWorkNumber: saveSchemeOfWork.Number,
			SubTopicNumber:     schemeList.SubTopicNumber,
			UserNumber:         userNumber,
		}
		if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme " + "sub topic " + err.Error()})
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

	var schemeOfWorkInfo requests.SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	err = schemeOfWorkInfo.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userNumber := context.Request.Header["User-Number"][0]
	if userNumber == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Can not determine user"})
		return
	}

	err = schemeOfWorkService.Update(schemeOfWorkInfo, userNumber)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme Updated Successfully!"})
}

func FetchSchemeOfWork(context *gin.Context) {

	academicYearNumber := context.Query("academic_year_number")
	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic year is required!"})
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

	var schemeOfWrkInfo schemeOfWorkModel.SchemeOfWork
	if err := config.DB.
		Where("academic_year_number = ? AND status = ?", academicYearNumber, constants.StatusActive).
		Where("campus_class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Preload("AcademicYear").
		Preload("CampusClass.CurriculumClass.Class").
		Preload("CampusClass.CurriculumClass").
		Preload("CampusClass").
		Preload("CampusClass.Campus").
		Preload("CampusClass.Campus.School").
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

		var conceptCount int64
		var conceptModel models4.Concept
		config.DB.Model(&conceptModel).
			Where("sub_topic_number = ?", subTopicList.SubTopicNumber).
			Count(&conceptCount)

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
			ConceptCount:       conceptCount,
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
		SchoolBadge:              schemeOfWrkInfo.CampusClass.Campus.School.Logo,
		SchoolName:               schemeOfWrkInfo.CampusClass.Campus.School.Name,
		CampusName:               schemeOfWrkInfo.CampusClass.Campus.Name,
		CampusEmail:              schemeOfWrkInfo.CampusClass.Campus.Email,
		CampusPhone:              schemeOfWrkInfo.CampusClass.Campus.Phone,
		CampusLocation:           schemeOfWrkInfo.CampusClass.Campus.Location,
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
				SubjectNumber: subtopicList.Topic.Strand.Subject.Number,
				SubjectName:   subtopicList.Topic.Strand.Subject.Name,
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

	//this function will set all headers for excel
	excelHeaders(sheetName, xlsx)

	// Iterate through the JSON data and write to the Excel file
	row := 2
	col := 1

	for _, lst := range dataInfo {

		//this function will add data to excel document
		addExcelInfoData(sheetName, col, row, lst, xlsx)

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

func addExcelInfoData(sheetName string, col int, row int, lst responses.SubTopicList, xlsx *excelize.File) {
	//Subject
	err := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.SubjectName)
	if err != nil {
		//return "", err
	}
	col++
	//------------------

	//Strand
	err1 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.StrandName)
	if err1 != nil {
		//return "", err1
	}
	col++
	//------------------

	//topic
	err2 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.TopicName)
	if err2 != nil {
		//return "", err2
	}
	col++
	//------------------

	//subtopic
	err3 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.Name)
	if err3 != nil {
		//return "", err3
	}
	col++
	//------------------

	//Estimated Time
	err4 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.EstimatedTime)
	if err4 != nil {
		//return "", err4
	}
	col++
	//------------------------

	//sorting code
	err5 := xlsx.SetCellValue(sheetName, columnIndexToLetters(col)+fmt.Sprint(row), lst.SortCode)
	if err5 != nil {
		//return "", err5
	}
	col++
	//---------------------------
}

func excelHeaders(sheetName string, xlsx *excelize.File) {

	err1 := xlsx.SetCellValue(sheetName, columnIndexToLetters(1)+fmt.Sprint(1), "Subject")
	if err1 != nil {
		//return "", err
	}

	err2 := xlsx.SetCellValue(sheetName, columnIndexToLetters(2)+fmt.Sprint(1), "Strand")
	if err2 != nil {
		//return "", err
	}

	err3 := xlsx.SetCellValue(sheetName, columnIndexToLetters(3)+fmt.Sprint(1), "Topic")
	if err3 != nil {
		//return "", err
	}

	err4 := xlsx.SetCellValue(sheetName, columnIndexToLetters(4)+fmt.Sprint(1), "Subtopic")
	if err4 != nil {
		//return "", err
	}

	err5 := xlsx.SetCellValue(sheetName, columnIndexToLetters(5)+fmt.Sprint(1), "Estimated Time")
	if err5 != nil {
		//return "", err1
	}

	err6 := xlsx.SetCellValue(sheetName, columnIndexToLetters(6)+fmt.Sprint(1), "sorting Code")
	if err6 != nil {
		//return "", err1
	}

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
		Status:             constants.StatusActive,
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
	readSchemeExcelFile(path, userNumber, schemeNumber, subjectNumber, tx, context)

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme uploaded Successfully!"})

}

func readSchemeExcelFile(filePath string, userNumber string, schemeNumber string, subjectNumber string, tx *gorm.DB, context *gin.Context) {

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
			schemeInfo := requests.ExcelSchemeOfWorkInfo{
				SubTopic:      getValueAtIndex(row, 3),
				EstimatedTime: getValueAtIndex(row, 4),
				SortingCode:   getValueAtIndex(row, 5),
			}
			saveSchemeOfWorkData(schemeInfo, userNumber, schemeNumber, subjectNumber, tx, context)
		}
	}

}

func saveSchemeOfWorkData(schemeInfo requests.ExcelSchemeOfWorkInfo, userNumber string, schemeNumber string, subjectNumber string, tx *gorm.DB, context *gin.Context) {

	codeNo, err := strconv.ParseInt(schemeInfo.SortingCode, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	estimatedTime, err := strconv.ParseInt(schemeInfo.EstimatedTime, 10, 64)
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
		EstimatedTime:      estimatedTime,
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

func InvalidateSchemeOfWork(context *gin.Context) {
	schemeOfWorkNumber := context.Param("scheme_of_work_number")
	userNumber := context.Request.Header["User-Number"][0]

	if userNumber == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not determine user "})
		return
	}

	err := schemeOfWorkService.InvalidateSchemeOfWork(schemeOfWorkNumber, userNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme of work invalidated!"})
}

func ActivateSchemeOfWork(context *gin.Context) {
	schemeOfWorkNumber := context.Param("scheme_of_work_number")
	userNumber := context.Request.Header["User-Number"][0]

	if userNumber == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not determine user "})
		return
	}

	err := schemeOfWorkService.ActivateSchemeOfWork(schemeOfWorkNumber, userNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme of work activated!"})
}
