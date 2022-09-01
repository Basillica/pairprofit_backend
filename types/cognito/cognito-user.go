package cognito

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	aws_types "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/types/appenv"
)

type CognitoUser struct {
	Email      string
	Username   string
	GivenName  string
	FamilyName string
	Custom     CognitoUserCustomAttributes
}
type CognitoUserCustomAttributes struct {
	ImageUri    string
	AccountType string
}

func (cu *CognitoUser) Save(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	adminCreateUserResult, err := cognito.AdminCreateUser(c, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   aws.String(cu.Email),
		UserAttributes: []aws_types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(cu.Email)},
			{Name: aws.String("custom:image_uri"), Value: aws.String(cu.Custom.ImageUri)},
			{Name: aws.String("custom:account_type"), Value: aws.String(cu.Custom.AccountType)},
			{Name: aws.String("given_name"), Value: aws.String(cu.GivenName)},
			{Name: aws.String("family_name"), Value: aws.String(cu.FamilyName)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
		MessageAction: "SUPPRESS",
	})
	if err != nil {
		return err
	}
	log.Println("Cognito AdminCreateUser: ", adminCreateUserResult.User.Enabled)
	return nil
}

func (cu *CognitoUser) SetPermanentPassword(c *gin.Context, password string) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	if _, err := cognito.AdminSetUserPassword(c, &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   aws.String(cu.Email),
		Password:   aws.String(password),
		Permanent:  true,
	}); err != nil {
		return err
	}
	log.Println("Cognito AdminSetUserPassword")
	return nil
}

func (cu *CognitoUser) ConfirmUserSignUp(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	if _, err := cognito.AdminConfirmSignUp(c, &cognitoidentityprovider.AdminConfirmSignUpInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   aws.String(cu.Email),
	}); err != nil {
		return err
	}
	return nil
}

func (cu *CognitoUser) UpdateUserAttributes(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	if _, err := cognito.AdminUpdateUserAttributes(c, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   aws.String(cu.Email),
		UserAttributes: []aws_types.AttributeType{
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	}); err != nil {
		return err
	}
	return nil
}

func (cu *CognitoUser) GetUser(c *gin.Context, password string) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	user, err := cognito.AdminGetUser(c, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
		Username:   aws.String(cu.Username),
	})
	if err != nil {
		log.Println("Error encountered", err.Error())
		return err
	}
	if user.UserStatus == "RESET_REQUIRED" {
		if err := cu.SetPermanentPassword(c, password); err != nil {
			panic(err)
		}
		if err := cu.UpdateUserAttributes(c); err != nil {
			panic(err)
		}
		return nil
	}
	return nil
}
