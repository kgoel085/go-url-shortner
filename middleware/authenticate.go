package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized !",
		})
		return
	}

	tokenUserId, tokenErr := utils.ValidateJWT(token, utils.LoginJwtType)
	if tokenErr != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", tokenErr.Error()),
		})
		return
	}

	context.Set(config.JWT_LOGGED_IN_USER, tokenUserId) // Set it to be available for further requests
	context.Next()
}
