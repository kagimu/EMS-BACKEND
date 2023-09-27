package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	authController "github.com/hamsoftug/wholesomeEd/domain/auth/controllers"
	"github.com/hamsoftug/wholesomeEd/domain/auth/middleware"
	"github.com/hamsoftug/wholesomeEd/domain/campus/controllers"
	classController "github.com/hamsoftug/wholesomeEd/domain/class/controllers"
	conceptController "github.com/hamsoftug/wholesomeEd/domain/curriculum/concept/controllers"
	curricullumController "github.com/hamsoftug/wholesomeEd/domain/curriculum/importCurriculum/controllers"
	learningAreaController "github.com/hamsoftug/wholesomeEd/domain/curriculum/learningarea/controllers"
	strandController "github.com/hamsoftug/wholesomeEd/domain/curriculum/strand/controllers"
	subjectClusterController "github.com/hamsoftug/wholesomeEd/domain/curriculum/subjectcluster/controllers"
	subTopicController "github.com/hamsoftug/wholesomeEd/domain/curriculum/subtopic/controllers"
	topicController "github.com/hamsoftug/wholesomeEd/domain/curriculum/topic/controllers"
	schoolController "github.com/hamsoftug/wholesomeEd/domain/school/controllers"
	subjectController "github.com/hamsoftug/wholesomeEd/domain/subject/controllers"
	timetableController "github.com/hamsoftug/wholesomeEd/domain/timetable/controllers"
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

		//	manage class
		main.POST("new-class", classController.NewClass)
		main.POST("new-class-level", classController.AddClassToLevel)
		main.GET("classes", classController.FetchClasses)
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

		main.POST("add-timetable-period", timetableController.AddPeriod)

	}
	//router.GET("/profile/:name", profile)
	return router
}
