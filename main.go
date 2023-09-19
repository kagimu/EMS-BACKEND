//package tutorView
package main

import (
	"database/sql"
	"fmt"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/api/routes"
	authModels "github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/campus/models"
	classModels "github.com/hamsoftug/wholesomeEd/domain/class/models"
	conceptModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/concept/models"
	models12 "github.com/hamsoftug/wholesomeEd/domain/curriculum/importCurriculum/controllers/models"
	learningAreaModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/learningarea/models"
	strandModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/strand/models"
	subjectClusterModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/subtopic/models"
	topicModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/topic/models"
	schoolModels "github.com/hamsoftug/wholesomeEd/domain/school/models"
	subjectModels "github.com/hamsoftug/wholesomeEd/domain/subject/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func main() {
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
