package controllers

import "github.com/gin-gonic/gin"

func IndexController(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "HELLO",
	})

}

func WordController(c *gin.Context) {
	word := c.Param("word")
	c.JSON(200, gin.H{
		"word":   word,
		"status": "ok",
	})
}
