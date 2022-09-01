package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/cognito"
	"pairprofit.com/x/types/email"
	"pairprofit.com/x/types/requests"
)

func LoginViaLink(c *gin.Context) {
	var req requests.LoginViaLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	if helpers.GetToken(req.Username) != req.Token {
		fmt.Println(helpers.GetToken(req.Username), req.Token)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token provided"})
		return
	}
	fmt.Println(c.Request.URL.Query(), "the quwery params")
	// localhost:8080/auth/link?_u=ezeabasilianthony@gmail.com&_c=gyajuafhYxntDn7IF7BFlN-oXMsWoyug&_t=6e269fd713a3724c8b18250576b7b0da
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	ca := &cognito.CognitoAuth{
		Username:  req.Username,
		GrantType: "password",
		Password:  helpers.Decrypt([]byte(appenv.CRYPTO_SECRET), req.Code),
		Code:      req.Code,
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

	cognitoUsername, err := getUsername(accessToken, appenv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	var recipients []sib_api_v3_sdk.SendSmtpEmailTo
	recipients = append(recipients, sib_api_v3_sdk.SendSmtpEmailTo{
		Email: req.Username,
		Name:  req.Username,
	})
	et := &email.EmailTemplate{
		TemplateFormatMap: map[string]string{"FirstName": "Anthony", "Transcription": "Some boring Transcription"},
	}
	et.Html = et.ParseHtmlFileToString("Templates/feedback.html")

	ems := &email.EmailSender{
		Recipients:  &recipients,
		Subject:     "You have succesfully created your account",
		HtmlContent: et.FormatTemplateString(),
	}
	ems.SendEmail(c)

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
