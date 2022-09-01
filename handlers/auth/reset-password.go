package auth

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/cookies"
	"pairprofit.com/x/types/requests"
)

func ResetPassword(c *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	// access, err := c.Cookie("access_token")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"success": false})
	// 	return
	// }

	// _, ok := helpers.GetUserOutput(c, access)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	// 	return
	// }

	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	randomString := helpers.Decrypt([]byte(appenv.CRYPTO_SECRET), req.Hash)
	secretHash := helpers.Decrypt([]byte(appenv.CRYPTO_SECRET), req.Token)
	if secretHash[len(secretHash)-8:] != randomString {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials provided"})
		return
	}

	if _, err := cognito.AdminEnableUser(
		c, &cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			Username:   &req.Username,
		},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ca := &cookies.CookieResp{
		Username: secretHash[:len(secretHash)-8],
		Password: secretHash[len(secretHash)-8:],
	}

	c, token, err := ca.PersistCookies(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, token)
}
