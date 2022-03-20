package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth_types "pairprofit.com/x/types/auth"
	auth_utils "pairprofit.com/x/utils"
)

func CreateTodo(c *gin.Context) {
	var td auth_types.Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	//Extract the access token metadata
	metadata, err := auth_utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userid, err := auth_utils.FetchAuth(c, metadata)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	td.UserID = userid
	//you can proceed to save the Todo to a database
	//but we will just return it to the caller:

	c.JSON(http.StatusCreated, td)
}
