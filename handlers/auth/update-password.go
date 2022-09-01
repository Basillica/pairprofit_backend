package auth

import (
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

func UpdatePassword(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	userOutput, ok := helpers.GetUserOutput(c, accessToken)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	var req requests.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	userEmail := helpers.GetUserAttr(userOutput, "email")
	var recipients []sib_api_v3_sdk.SendSmtpEmailTo
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	if _, err := cognito.AdminSetUserPassword(
		c, &cognitoidentityprovider.AdminSetUserPasswordInput{
			Password:   &req.Password,
			UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			Username:   &userEmail,
			Permanent:  true,
		},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipients = append(recipients, sib_api_v3_sdk.SendSmtpEmailTo{
		Email: req.Username,
		Name:  helpers.GetUserAttr(userOutput, "given_name"),
	})

	et := &email.EmailTemplate{
		TemplateFormatMap: map[string]string{
			"FirstName": helpers.GetUserAttr(userOutput, "given_name"),
		},
	}
	et.Html = et.ParseHtmlFileToString("Templates/passwordUpdate.html")
	ems := &email.EmailSender{
		Recipients:  &recipients,
		Subject:     "Password Update!",
		HtmlContent: et.FormatTemplateString(),
	}
	ems.SendEmail(c)

	c.JSON(http.StatusAccepted, gin.H{"status": "created"})
}
