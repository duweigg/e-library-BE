package main

import (
	"library/controllers"
	"library/initializers"
	"library/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()

}

func main() {
	router := gin.Default()

	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", controllers.CreateUser)
		userRouter.POST("/signin", controllers.SignIn)
		userRouter.GET("/info", middlewares.CheckAuth, controllers.GetUserInfo)
	}

	bookRouter := router.Group("/books")
	{
		bookRouter.GET("/list", controllers.GetBookList)
		bookRouter.POST("/borrow", middlewares.CheckAuth, controllers.BorrowBooks)
		bookRouter.POST("/extend", middlewares.CheckAuth, controllers.ExtendBooks)
		bookRouter.POST("/return", middlewares.CheckAuth, controllers.ReturnBooks)
	}
	router.Run()
}