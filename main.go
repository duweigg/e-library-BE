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
		bookRouter.POST("/borrow", middlewares.CheckAuth, controllers.CreateUser)
		bookRouter.POST("/extend", middlewares.CheckAuth, controllers.SignIn)
		bookRouter.POST("/return", middlewares.CheckAuth, controllers.GetUserInfo)
	}
	router.Run()
}