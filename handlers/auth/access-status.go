package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/cookies"
)

func VerifyAuth(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	secretHash := helpers.GetHashedString(appenv.COGNITO_CLIENT_SECRET)
	va, err, _ := cookies.VerifyTokenAndGetUserName(accessToken, secretHash, appenv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error(), "status": false})
		return
	}
	if (1655586732 - time.Now().Unix()) < 3000 {
		c.JSON(http.StatusAccepted, gin.H{"status": false, "expiry-time": va})
	}

	c.JSON(http.StatusAccepted, gin.H{"status": true, "expiry-time": va})
}
