package types

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"pairprofit.com/x/helpers"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type CookiePayload struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	User_id int    `json:"user_id"`
}

func (ck *TokenDetails) PersistCookies(c *gin.Context) (*gin.Context, *CookiePayload, error) {
	cookieDomain := helpers.GetEnv("COOKIE_DOMAIN", "")
	secureCookie := helpers.GetEnvAsBool("COOKIE_SECURE_ENABLE", false)
	httpOnly := true
	expirationData := time.Now().Local().Add(time.Hour * time.Duration(23))

	c.SetCookie("access_token", ck.AccessToken, int(ck.AtExpires), "/", cookieDomain, secureCookie, httpOnly)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", ck.RefreshToken, int(ck.RtExpires), "/", cookieDomain, secureCookie, false)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Access_UUID", ck.AccessUuid, 25920000, "/", cookieDomain, secureCookie, false)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Refresh_UUID", ck.RefreshUuid, 25920000, "/", cookieDomain, secureCookie, false)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("xexptk", expirationData.Format(time.RFC3339), int(ck.AtExpires), "/", cookieDomain, secureCookie, false)
	c.SetSameSite(http.SameSiteStrictMode)

	token := &CookiePayload{
		Refresh: ck.RefreshToken,
		Access:  ck.AccessToken,
	}
	return c, token, nil
}

func (ck *TokenDetails) RemoveCookies(c *gin.Context) (*gin.Context, error) {
	cookieDomain := helpers.GetEnv("COOKIE_DOMAIN", "")

	c.SetCookie("access_token", "", 30000, "/", cookieDomain, false, false)
	c.SetCookie("refresh_token", "", 30000, "/", cookieDomain, false, false)
	c.SetCookie("Access_UUID", "", 30000, "/", cookieDomain, false, false)
	c.SetCookie("Refresh_UUID", "", 30000, "/", cookieDomain, false, false)

	return c, nil
}
