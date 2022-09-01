package auth

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/email"
	"pairprofit.com/x/types/requests"
)

func ForgotPassword(c *gin.Context) {
	var req requests.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	var recipients []sib_api_v3_sdk.SendSmtpEmailTo
	randomString := helpers.GetUniqueHash(8)
	var userName string

	secretHash := helpers.Encrypt([]byte(appenv.CRYPTO_SECRET), randomString)
	secretToken := helpers.Encrypt([]byte(appenv.CRYPTO_SECRET), req.Username+randomString)
	if _, err := cognito.AdminSetUserPassword(
		c, &cognitoidentityprovider.AdminSetUserPasswordInput{
			Password:   &randomString,
			UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			Username:   &req.Username,
			Permanent:  true,
		},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := cognito.AdminDisableUser(
		c, &cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			Username:   &req.Username,
		},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateData := struct {
		Link string
	}{
		Link: fmt.Sprintf(
			"https://%s.pairprofit.com/%s/%s",
			appenv.APP_ENVIRONMENT,
			secretHash,
			secretToken,
		),
	}

	if res, err := cognito.AdminGetUser(c, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   &req.Username,
	}); err != nil {
		panic(err)
	} else {
		userName = helpers.AdminGetUserAttr(res, "given_name")
	}

	recipients = append(recipients, sib_api_v3_sdk.SendSmtpEmailTo{
		Email: req.Username,
		Name:  userName,
	})

	et := &email.EmailTemplate{
		TemplateFormatMap: map[string]string{"FirstName": userName, "Link": templateData.Link},
	}
	et.Html = et.ParseHtmlFileToString("Templates/passwordResetRequest.html")
	ems := &email.EmailSender{
		Recipients:  &recipients,
		Subject:     "Password Reset!",
		HtmlContent: et.FormatTemplateString(),
	}
	ems.SendEmail(c)

	c.JSON(http.StatusOK, gin.H{"success": true, "link": templateData.Link})
}
