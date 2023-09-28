//package tutorView
package main

import (
	"database/sql"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	models12 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	strandModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	subjectClusterModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	timetableModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var err error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})

	//config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	//All auto migrations
	//defer config.DB.Close()
	sqlDB, err := config.DB.DB()
	if err != nil {
		return
	}
	//All auto migrations
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {

		}
	}(sqlDB)

	errMigrate := config.DB.AutoMigrate(
		&authModels.User{},
		&authModels.Contact{},
		&authModels.PasswordReset{},
		&authModels.UserPassword{},
		&schoolModels.School{},
		&models.Campus{},
		&models.Level{},
		&classModels.Class{},
		&learningAreaModels.LearningArea{},
		&classModels.SubjectCLass{},
		&subjectClusterModels.SubjectCluster{},
		&subjectModels.Subject{},
		&strandModels.Strand{},
		&topicModels.Topic{},
		&subTopicModels.SubTopic{},
		&conceptModels.Concept{},
		&classModels.ClassLevel{},
		&models12.ImportCurriculum{},
		&classModels.Stream{},
		&timetableModels.Timetable{},
	)
	if errMigrate != nil {
		return
	}

	router := routes.SetupRoute()

	errPort := router.Run(":2054")
	if errPort != nil {
		fmt.Println("Error starting server:", errPort)
	}

}
