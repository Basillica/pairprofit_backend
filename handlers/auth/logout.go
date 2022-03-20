package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cookie_types "pairprofit.com/x/types/cookies"
	auth_utils "pairprofit.com/x/utils"
)

func Logout(c *gin.Context) {
	metadata, err := auth_utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := auth_utils.DeleteTokens(c, metadata)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	ts := &cookie_types.TokenDetails{}
	ts.RemoveCookies(c)
	c.JSON(http.StatusOK, "Successfully logged out")
}
