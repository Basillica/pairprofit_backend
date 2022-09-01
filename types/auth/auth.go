package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/types/appenv"
)

type PasswordRequestType struct {
	UsernameID string
	Password   string
	Code       string
}

func (pr *PasswordRequestType) ForgotPassword(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	cognitoInput := &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   &pr.UsernameID,
	}
	if _, err := cognito.AdminResetUserPassword(c, cognitoInput); err != nil {
		panic(err)
	}

	return nil
}

func (pr *PasswordRequestType) DisableUser(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	cognitoInput := &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   &pr.UsernameID,
	}
	if _, err := cognito.AdminDisableUser(c, cognitoInput); err != nil {
		panic(err)
	}

	return nil
}

func (pr *PasswordRequestType) EnableUser(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	log.Println(pr.UsernameID)
	cognitoInput := &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   &pr.UsernameID,
	}
	if _, err := cognito.AdminEnableUser(c, cognitoInput); err != nil {
		panic(err)
	}

	return nil
}

func (pr *PasswordRequestType) ConfirmPasswordReset(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	secretHash := secretHash(pr.UsernameID, appenv)
	resetPasswordInput := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(appenv.COGNITO_CLIENT_ID),
		Username:         &pr.UsernameID,
		SecretHash:       &secretHash,
		ConfirmationCode: &pr.Code,
		Password:         &pr.Password,
	}
	if _, err := cognito.ConfirmForgotPassword(c, resetPasswordInput); err != nil {
		panic(err)
	}

	log.Println("Cognito AdminSetUserPassword")
	return nil
}

func secretHash(username string, appenv *appenv.AppConfig) string {
	mac := hmac.New(sha256.New, []byte(appenv.COGNITO_CLIENT_SECRET))
	mac.Write([]byte(username + appenv.COGNITO_CLIENT_ID))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return secretHash
}
