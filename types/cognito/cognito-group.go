package cognito

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/types/appenv"
)

type CognitoGroup struct {
	Name string
}

func (cg *CognitoGroup) Upsert(c *gin.Context) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	if _, err := cognito.GetGroup(c, &cognitoidentityprovider.GetGroupInput{
		GroupName:  aws.String(cg.Name),
		UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
	}); err != nil {
		var nf *types.ResourceNotFoundException
		if errors.As(err, &nf) {
			log.Println("Cognito GetGroup: not found")
			_, err := cognito.CreateGroup(c, &cognitoidentityprovider.CreateGroupInput{
				GroupName:  aws.String(cg.Name),
				UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			})
			if err != nil {
				log.Println("Cognito CreateGroup: error")
				return err
			}
			log.Println("Cognito CreateGroup: success")
			return nil
		}
		log.Println("Cognito GetGroup: error")
		return err
	}

	log.Println("Cognito GetGroup: success ")
	return nil
}

func (cg *CognitoGroup) AddUser(c *gin.Context, cu *CognitoUser) error {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	if _, err := cognito.AdminAddUserToGroup(
		c,
		&cognitoidentityprovider.AdminAddUserToGroupInput{
			GroupName:  aws.String(cg.Name),
			UserPoolId: aws.String(appenv.COGNITO_POOL_ID),
			Username:   aws.String(cu.Username),
		},
	); err != nil {
		return err
	}
	log.Println("Cognito AddUserToGroup")
	return nil
}
