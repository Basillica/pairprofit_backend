package auth

import (
	"errors"
	"net/http"
	"time"

	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/cognito"
	"pairprofit.com/x/types/requests"
)

func TokenRefresh(c *gin.Context) {
	access, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	_, ok := helpers.GetUserOutput(c, access)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	var req requests.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	ca := &cognito.CognitoAuth{
		GrantType: req.GrantType,
	}

	accessToken, refreshTokenOrSession, err := ca.CreateToken(c)
	if err != nil {
		var nae *cognitotypes.NotAuthorizedException
		if errors.As(err, &nae) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	expirationTime := helpers.FormatExpTime(time.Now().Local().Add(time.Hour * time.Duration(23)))

	if len(*accessToken) < 1 {
		c.SetCookie("session", *refreshTokenOrSession, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, appenv.COOKIE_HTTPONLY)
	} else {
		c.SetCookie("access_token", *accessToken, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, appenv.COOKIE_HTTPONLY)
		c.SetSameSite(http.SameSiteStrictMode)
	}
	c.SetCookie("xexptk", expirationTime, 25920000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
	c.SetSameSite(http.SameSiteStrictMode)

	c.JSON(http.StatusCreated, gin.H{
		"access_token": *accessToken,
	})
}
