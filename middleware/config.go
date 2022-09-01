package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_config "pairprofit.com/x/types/appenv"
)

func ConfigMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appenv, err := LoadConfig(".")
		if err != nil {
			fmt.Println(err)
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.Set("appenv", &appenv)
	}
}

func LoadConfig(path string) (config _config.AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
