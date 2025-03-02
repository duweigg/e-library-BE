package main

import (
	"library/controllers"
	"library/initializers"
	"library/middlewares"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()
}

func main() {
	router := gin.Default()

	// Allow CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Change to frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", controllers.CreateUser)
		userRouter.POST("/signin", controllers.SignIn)
		userRouter.GET("/info", middlewares.CheckAuth, controllers.GetUserInfo)
	}

	bookController := controllers.NewBookController(initializers.DB)
	bookRouter := router.Group("/book")
	{
		bookRouter.POST("/list", bookController.GetBookList)
		bookRouter.POST("/borrow", middlewares.CheckAuth, bookController.BorrowBooks)
	}

	recordController := controllers.NewRecordController(initializers.DB)
	recordRouter := router.Group("/record")
	{
		recordRouter.POST("/list", middlewares.CheckAuth, recordController.GetRecordList)
		recordRouter.POST("/extend", middlewares.CheckAuth, recordController.ExtendRecords)
		recordRouter.POST("/return", middlewares.CheckAuth, recordController.ReturnRecords)
	}
	router.Run()
}
