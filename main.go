package main

import (
	"github.com/BottleneckStudio/WordJar/controllers"
	"github.com/BottleneckStudio/WordJar/middlewares/gzip"
	"github.com/gin-gonic/gin"
)

const port = ":3000"

func main() {

	router := gin.Default()
	// Before initializing routes,
	// all response should be gzipped
	router.Use(gzip.Gzip(gzip.DefaultCompression))
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
		v1.GET("/words/:word", controllers.WordController)
	}
}
