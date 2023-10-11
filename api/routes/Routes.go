package routes

import (
	academicYearController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/controllers"
	authController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/middleware"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/controllers"
	classController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/controllers"
	conceptController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/controllers"
	curricullumController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers"
	learningAreaController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/controllers"
	strandController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/controllers"
	subjectClusterController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/controllers"
	subTopicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/controllers"
	topicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/controllers"
	schoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/controllers"
	studyPeriodController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/controllers"
	subjectController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/controllers"
	timetableController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRoute() *gin.Engine {
	router := gin.Default()

	//setting cors in the system
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:2054", "https://wholesomeed.com"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization", "Accept-Encoding", "X-CSRF-Token", "Access-Control-Allow-Headers", "Content-Type", "user-number"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			///return origin == "https://github.com"
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	//this is to validate user toke before accessing any route
	//router.Use(middlewares.JWTMiddleware())
	router.MaxMultipartMemory = 20 << 20 // 20 MiB
	//router.Static("/", "./fileRepo")
	//router.Use(static.Serve("/api/fileRepo", static.LocalFile("/api/fileRepo", true)))
	//router.Use(static.Serve("/fileRepo", static.LocalFile("/fileRepo", true)))
	router.Use(static.Serve("fileRepo", static.LocalFile("fileRepo", true)))

	auth := router.Group("/auth")
	{
		auth.POST("login", authController.Login)
		auth.POST("register", authController.Register)
		auth.POST("set-password", authController.UserPassword)
		auth.POST("reset-password", authController.ResetPassword)
	}

	main := router.Group("/main", middleware.AuthorizeJWT())
	{
		//main.GET("profile/:id", middlewares.JWTMiddleware(), controllers.Profile)
		main.GET("profile", authController.Profile)

		//manage School
		main.POST("upload-curriculum", curricullumController.ImportCurriculum)
		main.GET("uploaded-curriculum", curricullumController.FetchUploadedCurriculum)

		//manage School
		main.POST("new-school", schoolController.NewSchool)
		main.POST("edit-school", schoolController.EditSchool)
		main.POST("delete-school", schoolController.DeleteSchool)
		main.GET("schools", schoolController.FetchSchools)
		main.GET("school-info", schoolController.FetchSchoolInfo)

		//manage campus
		main.POST("new-campus", controllers.NewCampus)
		main.POST("edit-campus", controllers.EditCampus)
		main.POST("delete-campus", controllers.DeleteCampus)
		main.GET("campuses", controllers.FetchCampuses)
		main.GET("campus-info", controllers.CampusInfo)
		main.GET("level", controllers.FetchLevel)
		main.GET("campus-level", controllers.FetchCampusLevel)

		//	manage class
		main.POST("new-class", classController.NewClass)
		main.POST("new-class-level", classController.AddClassToLevel)
		main.GET("classes", classController.FetchClasses)
		main.GET("classes-to-add", classController.FetchClassesToAdd)
		main.GET("class-level", classController.FetchClassLevel)
		main.GET("class-subjects", classController.FetchClassSubjects)
		main.POST("subject-to-class", classController.NewClassSubjects)

		//manage subject
		main.POST("new-subject", subjectController.NewSubject)
		main.GET("subjects", subjectController.FetchSubject)

		//manage Learning Area
		main.POST("new-learning-area", learningAreaController.NewLearningArea)
		main.GET("learning-areas", learningAreaController.FetchLearningArea)

		//subject cluster
		main.POST("new-subject-cluster", subjectClusterController.NewSubjectCluster)
		main.GET("subject-clusters", subjectClusterController.FetchSubjectCluster)

		//strand
		main.POST("new-strand", strandController.NewStrand)
		main.GET("strands", strandController.FetchStrand)

		//topic
		main.POST("new-topic", topicController.NewTopic)
		main.GET("topics", topicController.FetchTopic)

		//sub topic
		main.POST("new-sub-topic", subTopicController.NewSubTopic)
		main.GET("sub-topics", subTopicController.FetchSubTopic)

		//concepts
		main.POST("new-concept", conceptController.NewConcept)
		main.GET("concepts", conceptController.FetchConcept)

		//adds academic year
		main.POST("new-academic-year", academicYearController.NewAcademicYear)
		main.GET("fetch-academic-year", academicYearController.FetchAcademicYear)

		//adds study period
		main.POST("new-study-period", studyPeriodController.NewStudyPeriod)
		main.GET("fetch-study-period", studyPeriodController.FetchStudyPeriod)

		//timetable routes
		timetable := main.Group("/timetable")
		{
			timetable.POST("add-timetable-period", timetableController.AddPeriod)
			timetable.GET("get-timetable-periods", timetableController.FetchTimetablePeriodsGroupedByDay)
			timetable.GET("get-distinct-days", timetableController.GetDistinctDays)
			timetable.GET("get-periods-by-day/:day", timetableController.GetPeriodsByDay)
			timetable.GET("subjects-by-class-level-number/:class_level_number", timetableController.FetchSubjectsByClassLevelNumber)
			timetable.GET("get-classes-by-campus-number-and-level-number", timetableController.FetchClassLevelsByCampusNumberAndLevelNumber)
			timetable.GET("fetch-levels-by-campus-number/:campus_number", timetableController.FetchLevelsByCampusNumber)
		}

		stream := main.Group("/stream")
		{
			stream.POST("new", classController.AddClassStream)
			stream.GET("fetch", classController.FetchStreamsByClassLevelNumber)

		}

	}
	//router.GET("/profile/:name", profile)
	return router
}
