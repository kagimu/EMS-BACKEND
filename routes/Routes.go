package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/controllers"
	"github.com/hamsoftug/wholesomeEd/controllers/profile"
	"github.com/hamsoftug/wholesomeEd/controllers/userManagement"
	"time"
)

func SetupRoute() *gin.Engine {
	router := gin.Default()

	//setting cors in the system
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:2053", "https://api.tutorview.net"},
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
	router.Use(static.Serve("/fileRepo", static.LocalFile("/fileRepo", true)))

	auth := router.Group("/auth")
	{
		auth.POST("login", controllers.Login)
		auth.POST("register", controllers.Register)
	}

	main := router.Group("/main")
	{
		//main.GET("profile/:id", middlewares.JWTMiddleware(), controllers.Profile)
		main.GET("profile", profile.Profile)
		main.GET("admin-user-list", userManagement.Users)

	}
	//router.GET("/profile/:name", profile)
	return router
}
