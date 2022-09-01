package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
	"pairprofit.com/x/types/email"
)

func EmailEndPoint(c *gin.Context) {
	var recipients []sib_api_v3_sdk.SendSmtpEmailTo
	recipients = append(recipients, sib_api_v3_sdk.SendSmtpEmailTo{
		Email: "eienneceasar@gmail.com",
		Name:  "anthony",
	})

	et := &email.EmailTemplate{
		TemplateFormatMap: map[string]string{"Firstname": "Anthony", "Link": "https://google.com"},
	}
	et.Html = et.ParseHtmlFileToString("Templates/signUp.html")
	ems := &email.EmailSender{
		Recipients:  &recipients,
		Subject:     "You have succesfully created your account",
		HtmlContent: et.FormatTemplateString(),
	}
	ems.SendEmail(c)

	c.JSON(http.StatusAccepted, gin.H{"status": "sent"})
}
