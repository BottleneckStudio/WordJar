package controllers

import (
	"github.com/BottleneckStudio/WordJar/models"
	"github.com/gin-gonic/gin"
)

func IndexController(c *gin.Context) {
	data := gin.H{
		"version": "v1",
	}
	OutputDataAsJSON(c, data, "ok", "Welcome to WordJar API")
}

func WordController(c *gin.Context) {
	word := c.Param("word")
	result := models.CrawlWord(word)
	data := gin.H{
		"result": result,
	}
	OutputDataAsJSON(c, data, "ok", "1 result")
}
