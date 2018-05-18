package controllers

import (
	"github.com/gin-gonic/gin"
)

// OutputDataAsJSON outputs data as JSON format
func OutputDataAsJSON(c *gin.Context, data interface{}, status, msg string) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": msg,
		"data":    data,
	})
}

// OutputJSON outputs message as JSON format
func OutputJSON(c *gin.Context, status, msg string) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": msg,
	})
}

// OutputError outputs error message as JSON format
func OutputError(c *gin.Context) {
	c.JSON(500, gin.H{
		"status":  "error",
		"message": "something went wrong..",
	})
}
