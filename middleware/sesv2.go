package middleware

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gin-gonic/gin"
)

func Sesv2Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		awsCfg, err := config.LoadDefaultConfig(c)
		if err != nil {
			log.Println(err)
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		awsCfg.Region = "eu-central-1"
		c.Set("sesClient", sesv2.NewFromConfig(awsCfg))
	}
}
