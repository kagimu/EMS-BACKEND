package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	studentDetailsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	strandModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	subjectClusterModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	schemeOfWorkmodel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	learningMaterialsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	previousSchoolsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	timetableModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func main() {
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	errZone := os.Setenv("TZ", "Africa/Nairobi")
	if errZone != nil {
		return
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
		&eosModels.Entity{},
		&authModels.Contact{},
		&authModels.PasswordReset{},
		&authModels.UserPassword{},
		&schoolModels.School{},
		&schoolModels.SchoolCurriculum{},
		&models2.Curriculum{},
		&models.Campus{},
		&LevelModels.CampusLevel{},
		&LevelModels.Level{},
		&LevelModels.CurriculumLevel{},
		&classModel.Class{},
		&learningAreaModels.LearningArea{},
		&classModel.ClassLevelSubject{},
		&subjectClusterModels.SubjectCluster{},
		&subjectModels.Subject{},
		&strandModels.Strand{},
		&topicModels.Topic{},
		&subTopicModels.SubTopic{},
		&conceptModels.Concept{},
		&classModel.CampusClass{},
		&models2.ImportCurriculum{},
		&classModel.Stream{},
		&academicYearModel.AcademicYear{},
		&studyPeriodModel.StudyPeriod{},
		&models2.ImportCurriculum{},
		&classModel.Stream{},
		&timetableModels.StudyTimetable{},
		&timetableModels.StudyTimetablePeriods{},
		&timetableModels.TimetableSlots{},
		&eosModels.GeneralCalendar{},
		&eosModels.CalendarEvent{},
		&academicYearModel.AcademicYear{},
		&schemeOfWorkmodel.SchemeOfWork{},
		&schemeOfWorkmodel.SchemeOfWorkSubTopic{},
		&lessonPlanModel.LessonPlan{},
		&lessonPlanModel.LessonPlanSteps{},
		&studentDetailsModels.StudentDetail{},
		&classModel.CurriculumClass{},
		&sourceOfSchoolInfoModels.SourceOfSchoolInfo{},
		&parentModels.Parent{},
		&studentModels.Student{},
		&applicationModel.Application{},
		&previousSchoolsModel.PreviousSchools{},
		&studentModels.APN{},
		&learningMaterialsModels.LearningMaterial{},
	)

	if errMigrate != nil {
		return
	}

	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to initialize logger")
		}
	}(logger)

	logger.Info("logger construction succeeded")

	router := routes.SetupRoute()
	port := os.Getenv("APP_PORT")
	errPort := router.Run(":" + port)
	if errPort != nil {
		fmt.Println("Error starting server:", errPort)
	}

}
