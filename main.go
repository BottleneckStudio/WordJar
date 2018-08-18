package main

import (
	"net/http"
	"os"

	"github.com/BottleneckStudio/WordJar/controllers"
	"github.com/BottleneckStudio/WordJar/middlewares/cors"
	"github.com/BottleneckStudio/WordJar/middlewares/gzip"
	"github.com/gin-gonic/gin"
)

var port = ":" + os.Getenv("PORT")

func main() {
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.Middleware())
	// Before initializing routes,
	// all response should be gzipped
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Recover from panic/errors
	router.Use(gin.Recovery())
	initializeRoutes(router)

	router.Run(port)
	http.Handle("/", router)
	// appengine.Main()
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
