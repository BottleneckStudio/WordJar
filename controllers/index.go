package controllers

import "github.com/gin-gonic/gin"

func IndexController(c *gin.Context) {
	OutputJSON(c, "ok", "Welcome to Index")
}

func WordController(c *gin.Context) {
	word := c.Param("word")
	OutputJSON(c, "ok", word)
}
