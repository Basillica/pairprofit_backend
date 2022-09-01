package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/cognito"
	"pairprofit.com/x/types/cookies"
	"pairprofit.com/x/types/requests"
)

func Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	ca := &cognito.CognitoAuth{
		Username:  req.Username,
		GrantType: req.GrantType,
		Password:  req.Password,
	}

	accessToken, refreshTokenOrSession, err := ca.CreateToken(c)
	if err != nil {
		var nae *cognitotypes.NotAuthorizedException
		if errors.As(err, &nae) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "status": "cognito user not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	cognitoUsername, err := getUsername(accessToken, appenv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	expirationTime := helpers.FormatExpTime(time.Now().Local().Add(time.Hour * time.Duration(23)))

	if len(*accessToken) < 1 {
		c.SetCookie("session", *refreshTokenOrSession, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, appenv.COOKIE_HTTPONLY)
	} else {
		c.SetCookie("access_token", *accessToken, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, appenv.COOKIE_HTTPONLY)
		c.SetSameSite(http.SameSiteStrictMode)
	}

	if req.GrantType == "password" {
		c.SetCookie("refresh_token", *refreshTokenOrSession, 25920000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
		c.SetSameSite(http.SameSiteStrictMode)
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("xexptk", expirationTime, 25920000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("uuid", cognitoUsername, 25920000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
	c.SetSameSite(http.SameSiteStrictMode)

	if ca.GrantType == "refresh_token" {
		c.JSON(http.StatusCreated, gin.H{
			"access": *accessToken,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"access":  *accessToken,
			"refresh": *refreshTokenOrSession,
		})
	}
}

func getUsername(accessToken *string, appenv *appenv.AppConfig) (string, error) {
	mac := hmac.New(sha256.New, []byte(appenv.COGNITO_CLIENT_SECRET))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	_, err, cognitoUsername := cookies.VerifyTokenAndGetUserName(*accessToken, secretHash, appenv)
	return cognitoUsername, err
}
