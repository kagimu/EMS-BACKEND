package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/curriculum/importCurriculum/controllers/models"
	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func ImportCurriculum(context *gin.Context) {

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found!"})
		return
	}

	//validate maximum file size
	excelFile, fileError := context.FormFile("excel_file")

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

	fileNumber := uuid.NewV4().String() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	saveFileData := models2.ImportCurriculum{
		Number:     fileNumber,
		Name:       excelName,
		FileName:   fileName,
		UserNumber: userNumber,
	}
	config.DB.Save(&saveFileData)

	//read now the excel file
	readExcelFile(path)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Excel File uploaded Successfully!", "data": saveFileData})
}

func readExcelFile(filePath string) {
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
		if idx == 1 {
			for _, colCell := range row {
				fmt.Print(colCell, "\t")
			}
			//fmt.Println()
		}
	}

}

func FetchUploadedCurriculum(context *gin.Context) {

	var curriculumImport []models2.ImportCurriculum
	if err := config.DB.
		Order("created_at desc").
		Find(&curriculumImport).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculumImport})
}
