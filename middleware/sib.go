package middleware

import (
	"github.com/gin-gonic/gin"
	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/lib"
	"pairprofit.com/x/helpers"
)

func SIBMiddleware() gin.HandlerFunc {
	cfg := sib_api_v3_sdk.NewConfiguration()
	//Configure API key authorization: api-key
	sib_api_key := helpers.GetEnv("SIB_API_KEY", "")
	cfg.AddDefaultHeader("api-key", sib_api_key)
	//Configure API key authorization: partner-key
	cfg.AddDefaultHeader("partner-key", sib_api_key)

	sibClient := sib_api_v3_sdk.NewAPIClient(cfg)
	return func(c *gin.Context) {
		c.Set("sibClient", sibClient)
	}
}
