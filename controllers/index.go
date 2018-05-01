package controllers

import "github.com/gin-gonic/gin"

func IndexController(c *gin.Context) {
	data := gin.H{
		"version": "v1",
	}
	OutputDataAsJSON(c, data, "ok", "Welcome to WordJar API")
}

func WordController(c *gin.Context) {
	word := c.Param("word")
	OutputJSON(c, "ok", word)
}
