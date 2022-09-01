package email

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
)

type Email struct {
	TemplateName string
	TemplateData string
	Subject      string
	ToAddress    string
}

func (e *Email) UpdateTemplate(c *gin.Context) error {
	//version := c.MustGet("version").(string)
	sesClient := c.MustGet("sesClient").(*sesv2.Client)
	//templateName := e.TemplateName + "-" + version
	templateName := e.TemplateName

	//uncomment for updating/overwriting the template

	templateBytes, err := ioutil.ReadFile("./emails/" + e.TemplateName + ".html")
	if err != nil {
		panic(err)
	}
	template := string(templateBytes)
	if _, err := sesClient.UpdateEmailTemplate(c, &sesv2.UpdateEmailTemplateInput{
		TemplateContent: &sesv2types.EmailTemplateContent{
			Html:    aws.String(template),
			Subject: aws.String(e.Subject),
		},
		TemplateName: aws.String(templateName),
	}); err != nil {
		panic(err)
	}

	if _, err := sesClient.GetEmailTemplate(c, &sesv2.GetEmailTemplateInput{
		TemplateName: aws.String(templateName),
	}); err != nil {
		var nf *sesv2types.NotFoundException
		if errors.As(err, &nf) {
			log.Println("SES template not found")
			log.Println("error:", nf)

			templateBytes, err := ioutil.ReadFile("./emails/" + e.TemplateName + ".html")
			if err != nil {
				panic(err)
			}
			template := string(templateBytes)

			if _, err := sesClient.CreateEmailTemplate(c, &sesv2.CreateEmailTemplateInput{
				TemplateContent: &sesv2types.EmailTemplateContent{
					Html:    aws.String(template),
					Subject: aws.String(e.Subject),
				},
				TemplateName: aws.String(templateName),
			}); err != nil {
				return err
			}
			log.Println("SES template created: ", templateName)
			return nil
		}
		return err
	}
	log.Println("SES template exists: ", templateName)
	return nil
}

func (e *Email) Send(c *gin.Context) error {
	//version := c.MustGet("version").(string)
	sesClient := c.MustGet("sesClient").(*sesv2.Client)
	//templateName := e.TemplateName + "-" + version
	templateName := e.TemplateName

	sendOutput, err := sesClient.SendEmail(c, &sesv2.SendEmailInput{
		Content: &sesv2types.EmailContent{
			Template: &sesv2types.Template{
				TemplateData: aws.String(e.TemplateData),
				TemplateName: aws.String(templateName),
			},
		},
		Destination: &sesv2types.Destination{
			ToAddresses: []string{e.ToAddress},
		},
		FromEmailAddress: aws.String(helpers.GetEnv("SES_SENDER_EMAIL", "")),
	})

	log.Println(&sendOutput.ResultMetadata)
	if err != nil {
		return err
	}
	log.Println("Email template: ", templateName)
	log.Println("Email sent: ", *sendOutput.MessageId)

	return nil
}

type Notification struct {
	Subject           string  `json:"subject"`
	Receiver          string  `json:"receiver"`
	ReceiverFirstname *string `json:"receiver_firstname"`
	ReceiverLastname  *string `json:"receiver_lastname"`
	CreatorFirstname  *string `json:"creator_firstname"`
	CreatorLastname   *string `json:"creator_lastname"`
	Template          string  `json:"template"`
	TopicID           int     `json:"topic_id"`
	RecordingID       int     `json:"recording_id"`
	Transcription     string  `json:"transcription"`
	Duration          float64 `json:"duration"`
	CreationTime      string  `json:"creation_time"`
	GroupID           int     `json:"group_id"`
	ImageURI          string  `json:"imageUri"`
}

func (n *Notification) Notify(c *gin.Context) {
	e := &Email{
		TemplateName: n.Template,
		TemplateData: "{\"receiver_firstname\":\"" + *n.ReceiverFirstname + "\"," +
			"\"topic_id\":\"" + strconv.Itoa(n.TopicID) + "\"," +
			"\"transcription\":\"" + n.Transcription + "\"," +
			"\"creator_firstname\":\"" + *n.CreatorFirstname + "\"," +
			"\"creator_lastname\":\"" + *n.CreatorLastname + "\"," +
			"\"recording_id\":\"" + strconv.Itoa(n.RecordingID) + "\"," +
			"\"duration\":\"" + fmt.Sprint(n.Duration) + "\"," +
			"\"creationTime\":\"" + n.CreationTime + "\"," +
			"\"imageUri\":\"" + n.ImageURI + "\"," +
			"\"group_id\":\"" + strconv.Itoa(n.GroupID) + "\"}",
		Subject:   n.Subject,
		ToAddress: n.Receiver,
	}

	if err := e.Send(c); err != nil {
		panic(err)
	}
}

type FollowUpEmail struct {
	Subject         string `json:"subject"`
	Template        string `json:"template"`
	SenderFirstname string `json:"sender_firstname"`
	SenderEmail     string `json:"sender_email"`
	Receiver        string `json:"emails"`
	CreationTime    string `json:"creationTime"`
	Duration        string `json:"duration"`
	Transcription   string `json:"transcription"`
	ImageURI        string `json:"imageUri"`
	Hash            string `json:"_hash"`
}

func (f *FollowUpEmail) SendFollowUpEmail(c *gin.Context) {
	env := helpers.GetEnv("APP_ENVIRONMENT", "staging")
	link := "https://" + env + ".getpairprofit.com/reply/" + f.Hash
	e := &Email{
		TemplateName: f.Template,
		TemplateData: "{\"creationTime\":\"" + f.CreationTime + "\"," +
			"\"transcription\":\"" + f.Transcription + "\"," +
			"\"duration\":\"" + fmt.Sprint(f.Duration) + "\"," +
			"\"imageUri\":\"" + f.ImageURI + "\"," +
			"\"email\":\"" + f.SenderEmail + "\"," +
			"\"link\":\"" + link + "\"," +
			"\"first_name\":\"" + f.SenderFirstname + "\"}",
		Subject:   f.Subject,
		ToAddress: f.Receiver,
	}

	if err := e.Send(c); err != nil {
		panic(err)
	}
}

type ProfileInquiry struct {
	Receivers     []string
	Templates     []string
	CreationTime  string
	Duration      string
	Transcription string
	Sender        string
	ImageUri      string
	FirstName     string
	Profile       bool
}

type pairprofitInquiry struct {
	Receivers     []string
	Templates     []string
	CreationTime  string
	Duration      string
	Transcription string
	Sender        string
	ImageUri      string
	FirstName     string
	Profile       bool
	Hash          string
}

func (f *ProfileInquiry) SendProfileVLInquiry(c *gin.Context) {
	ee := &Email{
		TemplateName: f.Templates[0],
		TemplateData: "{\"creationTime\":\"" + f.CreationTime + "\"," +
			"\"transcription\":\"" + f.Transcription + "\"," +
			"\"duration\":\"" + fmt.Sprint(f.Duration) + "\"," +
			"\"imageUri\":\"" + f.ImageUri + "\"," +
			"\"email\":\"" + f.Sender + "\"}",
		Subject:   "Check your profile inbox",
		ToAddress: f.Receivers[0],
	}
	//
	//if err := ee.UpdateTemplate(c); err != nil {
	//	panic(err)
	//}

	er := &Email{
		TemplateName: f.Templates[1],
		TemplateData: "{\"first_name\":\"" + f.FirstName + "\"," +
			"\"imageUri\":\"" + f.ImageUri + "\"," +
			"\"transcription\":\"" + f.Transcription + "\"}",
		//Subject:   "Thanks for connecting with " + f.FirstName,
		Subject:   "Thanks for connecting over pairprofit",
		ToAddress: f.Receivers[1],
	}

	//if err := er.UpdateTemplate(c); err != nil {
	//	panic(err)
	//}

	if err := ee.Send(c); err != nil {
		panic(err)
	}
	if err := er.Send(c); err != nil {
		panic(err)
	}
}

func (f *pairprofitInquiry) pairprofitInquiry(c *gin.Context) {
	ee := &Email{
		TemplateName: f.Templates[0],
		TemplateData: "{\"creationTime\":\"" + f.CreationTime + "\"," +
			"\"transcription\":\"" + f.Transcription + "\"," +
			"\"duration\":\"" + fmt.Sprint(f.Duration) + "\"," +
			"\"imageUri\":\"" + f.ImageUri + "\"," +
			"\"email\":\"" + f.Sender + "\"}",
		Subject:   "Check your profile inbox",
		ToAddress: f.Receivers[0],
	}

	//if err := ee.UpdateTemplate(c); err != nil {
	//	panic(err)
	//}

	er := &Email{
		TemplateName: f.Templates[1],
		TemplateData: "{\"first_name\":\"" + f.FirstName + "\"," +
			"\"imageUri\":\"" + f.ImageUri + "\"," +
			"\"hash\":\"" + f.Hash + "\"," +
			"\"transcription\":\"" + f.Transcription + "\"}",
		//Subject:   "Thanks for connecting with " + f.FirstName,
		Subject:   "Thanks for connecting over pairprofit",
		ToAddress: f.Receivers[1],
	}

	//if err := er.UpdateTemplate(c); err != nil {
	//	panic(err)
	//}

	if err := ee.Send(c); err != nil {
		panic(err)
	}
	if err := er.Send(c); err != nil {
		panic(err)
	}
}
