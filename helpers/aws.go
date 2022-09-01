package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

func GetUserAttr(userOutput *cognitoidentityprovider.GetUserOutput, val string) (res string) {
	for _, v := range userOutput.UserAttributes {
		if *v.Name == val {
			res = *v.Value
		}
	}
	return
}

func AdminGetUserAttr(userOutput *cognitoidentityprovider.AdminGetUserOutput, val string) (res string) {
	for _, v := range userOutput.UserAttributes {
		if *v.Name == val {
			res = *v.Value
		}
	}
	return
}

func GetUserOutput(c *gin.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, bool) {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	userOutput, err := cognito.GetUser(c, &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return userOutput, true
}
