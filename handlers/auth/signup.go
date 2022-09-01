package auth

import (
	"errors"
	"fmt"
	"net/http"

	aws_types "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
	"pairprofit.com/x/types/auth"
	"pairprofit.com/x/types/cognito"
	"pairprofit.com/x/types/email"
	"pairprofit.com/x/types/requests"
)

func SignUp(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidatePayload(err, c)
		return
	}

	if _, err := auth.AccountTypes.Parse(req.AccountType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong account type provided"})
		return
	}

	cu := &cognito.CognitoUser{
		Email:      *req.Username,
		GivenName:  req.Firstname,
		FamilyName: req.Lastname,
	}
	if req.ImageUri != nil {
		cu.Custom.ImageUri = *req.ImageUri
	}

	if err := cu.Save(c); err != nil {
		var cogerr *aws_types.UsernameExistsException
		if errors.As(err, &cogerr) {
			c.JSON(http.StatusForbidden, gin.H{"status": "User with the provided email already exists"})
			return
		}
		panic(err)
	}
	if err := cu.SetPermanentPassword(c, req.Password); err != nil {
		panic(err)
	}

	so := &email.SibObject{
		Email:     *req.Username,
		Firstname: &req.Firstname,
		Lastname:  &req.Lastname,
		StatusNo:  2,
	}
	so.UpdateContact(c, so.Firstname, so.Lastname, &so.StatusNo)

	var recipients []sib_api_v3_sdk.SendSmtpEmailTo
	recipients = append(recipients, sib_api_v3_sdk.SendSmtpEmailTo{
		Email: *req.Username,
		Name:  req.Firstname + " " + req.Lastname,
	})

	token := helpers.GetToken(*req.Username)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	cryptoText := helpers.Encrypt([]byte(appenv.CRYPTO_SECRET), req.Password)

	templateData := struct {
		Link string
	}{
		Link: fmt.Sprintf(
			"https://%s.pairprofit.com/confirmation/%s/%s/%s",
			appenv.APP_ENVIRONMENT,
			*req.Username,
			cryptoText,
			token,
		),
	}

	et := &email.EmailTemplate{
		TemplateFormatMap: map[string]string{"Firstname": req.Firstname, "Link": templateData.Link},
	}
	et.Html = et.ParseHtmlFileToString("Templates/signUp.html")

	ems := &email.EmailSender{
		Recipients:  &recipients,
		Subject:     "You have succesfully created your account",
		HtmlContent: et.FormatTemplateString(),
	}
	ems.SendEmail(c)

	c.JSON(http.StatusAccepted, gin.H{"email": req.Username, "status": "Created!", "confirmation_link": templateData})
}
