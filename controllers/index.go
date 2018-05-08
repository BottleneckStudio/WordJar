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
	locale := c.Query("tl")
	result := models.CrawlWord(word, locale)
	data := gin.H{
		"result": result,
	}
	OutputDataAsJSON(c, data, "ok", "1 result")
}
