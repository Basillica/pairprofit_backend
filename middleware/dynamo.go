package middleware

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"pairprofit.com/x/types/appenv"

	"github.com/gin-gonic/gin"
)

func DynamoDBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		awsCfg, err := config.LoadDefaultConfig(c)
		appenv := c.MustGet("appenv").(*appenv.AppConfig)
		if err != nil {
			log.Println(err)
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		awsCfg.Region = "eu-central-1"
		if env := appenv.APP_ENVIRONMENT; env == "local" {
			c.Set("dbClient", dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
				o.EndpointResolver = dynamodb.EndpointResolverFromURL("https://0.0.0.0:4566")
			}))
		} else {
			c.Set("dbClient", dynamodb.NewFromConfig(awsCfg))
		}

	}
}
