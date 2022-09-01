package auth

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	auth "pairprofit.com/x/types/cookies"
)

func TokenRevoke(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	_, ok := helpers.GetUserOutput(c, accessToken)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	if _, err := cognito.RevokeToken(c, &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(appenv.COGNITO_CLIENT_ID),
		Token:        aws.String(refreshToken),
		ClientSecret: aws.String(appenv.COGNITO_CLIENT_SECRET),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	log.Println("Cognito RevokeToken")
	cu := &auth.CookieResp{}
	c, _ = cu.RemoveCookies(c)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
