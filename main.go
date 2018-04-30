package main

import (
	"github.com/BottleneckStudio/WordJar/controllers"
	"github.com/gin-gonic/gin"
)

const port = ":3000"

func main() {

	router := gin.Default()

	initializeRoutes(router)

	router.Run(port)
}

func initializeRoutes(origRouter *gin.Engine) {

	router := origRouter.Group("")
	{
		router.GET("/", controllers.IndexController)
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/words", controllers.IndexController)
	}
}
