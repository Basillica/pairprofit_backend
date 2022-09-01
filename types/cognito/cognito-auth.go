package cognito

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"pairprofit.com/x/types/appenv"
)

type CognitoAuth struct {
	Username  string
	Password  string
	GrantType string
	Code      string
}

func (ca *CognitoAuth) CreateToken(c *gin.Context) (*string, *string, error) {
	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	// calculate secret hash
	clientSecret := appenv.COGNITO_CLIENT_SECRET
	mac := hmac.New(sha256.New, []byte(clientSecret))

	var input *cognitoidentityprovider.InitiateAuthInput
	var inputOtp *cognitoidentityprovider.AdminRespondToAuthChallengeInput
	if ca.GrantType == "password" {
		// take username for secret hash
		mac.Write([]byte(ca.Username + appenv.COGNITO_CLIENT_ID))
		secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		input = &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			ClientId: aws.String(appenv.COGNITO_CLIENT_ID),
			AuthParameters: map[string]string{
				"USERNAME":    ca.Username,
				"PASSWORD":    ca.Password,
				"SECRET_HASH": secretHash,
			},
		}
	} else if ca.GrantType == "refresh_token" {
		cognitoUsername, err := c.Cookie("uuid")
		if err != nil {
			panic(err.Error())
		}
		mac.Write([]byte(cognitoUsername + appenv.COGNITO_CLIENT_ID))
		secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		refreshToken, _ := c.Cookie("refresh_token")
		input = &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
			ClientId: aws.String(appenv.COGNITO_CLIENT_ID),
			AuthParameters: map[string]string{
				"REFRESH_TOKEN": refreshToken,
				"SECRET_HASH":   secretHash,
			},
		}
	} else if ca.GrantType == "otp_request" {
		// take username_id for secret hash
		mac.Write([]byte(ca.Username + appenv.COGNITO_CLIENT_ID))
		secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		input = &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeCustomAuth,
			ClientId: aws.String(appenv.COGNITO_CLIENT_ID),
			AuthParameters: map[string]string{
				"USERNAME":    ca.Username,
				"SECRET_HASH": secretHash,
			},
		}

		output, err := cognito.InitiateAuth(c, input)
		if err != nil {
			log.Println("Cognito InitiateAuth: error", err)
			return nil, nil, err
		}
		return aws.String(""), output.Session, nil
	} else if ca.GrantType == "otp_token" {
		// take username_id for secret hash
		session, _ := c.Cookie("session")
		mac.Write([]byte(ca.Username + appenv.COGNITO_CLIENT_ID))
		secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		inputOtp = &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
			ChallengeName: types.ChallengeNameTypeCustomChallenge,
			ClientId:      aws.String(appenv.COGNITO_CLIENT_ID),
			UserPoolId:    aws.String(appenv.COGNITO_POOL_ID),
			ChallengeResponses: map[string]string{
				"USERNAME":    ca.Username,
				"ANSWER":      ca.Code,
				"SECRET_HASH": secretHash,
			},
			Session: aws.String(session),
		}

		respondOutput, err := cognito.AdminRespondToAuthChallenge(c, inputOtp)
		if err != nil {
			log.Println("Cognito InitiateAuth: error", err)
			return nil, nil, err
		}

		return respondOutput.AuthenticationResult.AccessToken, respondOutput.AuthenticationResult.RefreshToken, nil
	}

	createTokenOutput, err := cognito.InitiateAuth(c, input)
	if err != nil {
		log.Println("Cognito InitiateAuth: error", err)
		return nil, nil, err
	}
	return createTokenOutput.AuthenticationResult.AccessToken, createTokenOutput.AuthenticationResult.RefreshToken, nil
}
