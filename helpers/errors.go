package helpers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"pairprofit.com/x/types/requests"
)

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "email":
		return "This field must be a valid email address"
	case "oneof":
		return "This field must be one of password, refresh_token, otp_request, or otp_token if grant_type or one of regularAccount, serviceProvider if account_type"
	case "required":
		return "This field is required field and has to be provided"
	case "alpha":
		return "This field should contain only english alphabets"
	case "lte":
		return "This field should be less than " + fe.Param()
	case "gte":
		return "This field should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func ValidatePayload(err error, c *gin.Context) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]requests.ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = requests.ErrorMsg{
				Field:   fe.Field(),
				Message: GetErrorMsg(fe)}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
	}
}
