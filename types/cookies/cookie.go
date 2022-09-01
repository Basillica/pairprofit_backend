package cookies

import (
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/types/appenv"
)

type CookieResp struct {
	Password string
	Username string
}

type CookiePayload struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	User_id int    `json:"user_id"`
}

func (ca *CookieResp) PersistCookies(c *gin.Context) (*gin.Context, *CookiePayload, error) {

	cognito := c.MustGet("cognito").(*cognitoidentityprovider.Client)
	appenv := c.MustGet("appenv").(*appenv.AppConfig)
	// calculate secret hash
	mac := hmac.New(sha256.New, []byte(appenv.COGNITO_CLIENT_SECRET))

	var input *cognitoidentityprovider.InitiateAuthInput
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

	createTokenOutput, err := cognito.InitiateAuth(c, input)
	expirationTime := helpers.FormatExpTime(time.Now().Local().Add(time.Hour * time.Duration(23)))

	if err != nil {
		log.Println("Cognito InitiateAuth: error", err)
		return nil, &CookiePayload{}, err
	}

	c.SetCookie("access_token", *createTokenOutput.AuthenticationResult.AccessToken, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, appenv.COOKIE_HTTPONLY)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", *createTokenOutput.AuthenticationResult.RefreshToken, 86000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
	c.SetSameSite(http.SameSiteStrictMode)
	// c.SetCookie("UUID", authUser.UUID, 25920000, "/", cookieDomain, secureCookie, false)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("xexptk", expirationTime, 25920000, "/", appenv.COOKIE_DOMAIN, appenv.COOKIE_SECURE_ENABLE, false)
	c.SetSameSite(http.SameSiteStrictMode)

	token := &CookiePayload{
		Refresh: *createTokenOutput.AuthenticationResult.RefreshToken,
		Access:  *createTokenOutput.AuthenticationResult.AccessToken,
		// User_id: authUser.ID,
	}
	return c, token, nil
}

func (ca *CookieResp) RemoveCookies(c *gin.Context) (*gin.Context, error) {
	appenv := c.MustGet("appenv").(*appenv.AppConfig)

	c.SetCookie("access_token", "", 30000, "/", appenv.COOKIE_DOMAIN, false, false)
	c.SetCookie("refresh_token", "", 30000, "/", appenv.COOKIE_DOMAIN, false, false)
	c.SetCookie("UUID", "", 30000, "/", appenv.COOKIE_DOMAIN, false, false)
	c.SetCookie("xexptk", "", 30000, "/", appenv.COOKIE_DOMAIN, false, false)

	return c, nil
}

func VerifyTokenAndGetUserName(tokenString, hmacSampleSecret string, appenv *appenv.AppConfig) (float64, error, string) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	var exp float64 = 0

	config := &Config{
		CognitoRegion:     appenv.AWS_REGION,
		CognitoUserPoolID: appenv.COGNITO_POOL_ID,
	}
	na := NewAuth(config)
	na.SaveJWK()
	for _, val := range na.jwk.Keys {
		claims := token.Claims.(jwt.MapClaims)
		if val.Kid == token.Header["kid"] && claims["iss"] == "https://cognito-idp."+config.CognitoRegion+".amazonaws.com/"+config.CognitoUserPoolID {
			return claims["exp"].(float64), nil, claims["username"].(string)
		}
	}
	return exp, errors.New("The provided token could not be verified"), ""
}

type Auth struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
}

type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
}

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func NewAuth(config *Config) *Auth {
	a := &Auth{
		cognitoRegion:     config.CognitoRegion,
		cognitoUserPoolID: config.CognitoUserPoolID,
	}

	a.jwkURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.cognitoRegion, a.cognitoUserPoolID)
	err := a.SaveJWK()
	if err != nil {
		log.Fatal(err)
	}

	return a
}

func (a *Auth) SaveJWK() error {
	req, err := http.NewRequest("GET", a.jwkURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwk := new(JWK)
	err = json.Unmarshal(body, jwk)
	if err != nil {
		return err
	}

	a.jwk = jwk
	return nil
}

func (a *Auth) JWK() *JWK {
	return a.jwk
}

func (a *Auth) JWKURL() string {
	return a.jwkURL
}

func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}
