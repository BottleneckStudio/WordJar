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

	if result.Created.IsZero() == true {
		OutputError(c)
		return
	}
	OutputDataAsJSON(c, result, "ok", "1 result")
}
