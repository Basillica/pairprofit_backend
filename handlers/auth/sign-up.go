package auth

import "github.com/gin-gonic/gin"

func SignUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
