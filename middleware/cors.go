package middleware

import (
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
)

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigin := helpers.GetEnv("ALLOWED_ORIGIN", "http://localhost:3000")
	alllowedHeaders := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, x-csrftoken, Authorization, accept, origin, Cache-Control, X-Requested-With, sentry-trace, x_bearer_token, x-bearer-token, content-disposition"
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", alllowedHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization, Set-Cookie")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
