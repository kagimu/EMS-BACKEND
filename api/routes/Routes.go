package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	authController "github.com/hamsoftug/wholesomeEd/domain/auth/controllers"
	"github.com/hamsoftug/wholesomeEd/domain/campus/controllers"
	controllers3 "github.com/hamsoftug/wholesomeEd/domain/class/controllers"
	controllers5 "github.com/hamsoftug/wholesomeEd/domain/learningarea/controllers"
	controllers2 "github.com/hamsoftug/wholesomeEd/domain/school/controllers"
	controllers4 "github.com/hamsoftug/wholesomeEd/domain/subject/controllers"
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
	router.Use(static.Serve("/api/fileRepo", static.LocalFile("/api/fileRepo", true)))

	auth := router.Group("/auth")
	{
		auth.POST("login", authController.Login)
		auth.POST("register", authController.Register)
		auth.POST("set-password", authController.UserPassword)
		auth.POST("reset-password", authController.ResetPassword)
	}

	main := router.Group("/main")
	{
		//main.GET("profile/:id", middlewares.JWTMiddleware(), controllers.Profile)
		main.GET("profile", authController.Profile)

		//manage School
		main.POST("new-school", controllers2.NewSchool)
		main.GET("schools", controllers2.FetchSchools)

		//manage campus
		main.POST("new-campus", controllers.NewCampus)
		main.GET("campuses", controllers.FetchCampuses)
		main.GET("level", controllers.FetchLevel)

		//	manage class
		main.POST("new-class", controllers3.NewClass)
		main.GET("classes", controllers3.FetchClasses)

		//manage subject
		main.POST("new-subject", controllers4.NewSubject)
		main.GET("subjects", controllers4.FetchSubject)

		//manage Learning Area
		main.POST("new-learning-area", controllers5.NewLearningArea)
		main.GET("learning-areas", controllers5.FetchLearningArea)
		main.POST("subject-to-learning-area", controllers5.AddSubjectToLearningArea)

	}
	//router.GET("/profile/:name", profile)
	return router
}
