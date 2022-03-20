package middleware

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func S3Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		awsCfg, err := config.LoadDefaultConfig(c)
		if err != nil {
			log.Println(err)
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		awsCfg.Region = "eu-central-1"
		c.Set("s3Client", s3.NewFromConfig(awsCfg))
	}
}
