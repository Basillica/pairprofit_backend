package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"pairprofit.com/x/helpers"
)

func RedisMiddleware() gin.HandlerFunc {
	var client *redis.Client
	dsn := helpers.GetEnv("REDIS_DNS", "redis:6379")
	client = redis.NewClient(&redis.Options{
		Addr:       dsn,
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	return func(c *gin.Context) {
		c.Set("redisClient", client)
	}
}
