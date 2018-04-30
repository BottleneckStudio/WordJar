package main

import (
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
		router.GET("/", sampleHandler)
	}
}

func sampleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "HELLO",
	})
}
