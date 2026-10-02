package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()

	db := config.DB.Debug()
	var admissions []admissionModel.Admission

	query := db.
		Joins("JOIN students ON students.number = admissions.student_number").
		Where("admissions.school_number = ?", "NON_EXISTENT_SCHOOL")

	// Simulate paginateFunc
	dbClone := query.Session(&gorm.Session{}).Select("COUNT(*)")
	var total int64
	err = dbClone.Count(&total).Error
	fmt.Printf("Count err: %v, total: %d\n", err, total)

	err = query.Find(&admissions).Error
	fmt.Printf("Find err: %v, len: %d\n", err, len(admissions))
}
