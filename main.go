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

	bookRouter := router.Group("/book")
	{
		bookRouter.POST("/list", controllers.GetBookList)
		bookRouter.POST("/borrow", middlewares.CheckAuth, controllers.BorrowBooks)
	}
	recordRouter := router.Group("/record")
	{
		recordRouter.POST("/list", middlewares.CheckAuth, controllers.GetRecordList)
		recordRouter.POST("/extend", middlewares.CheckAuth, controllers.ExtendRecords)
		recordRouter.POST("/return", middlewares.CheckAuth, controllers.ReturnRecords)
	}
	router.Run()
}
