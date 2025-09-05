package middleware

import (
	"fmt"
	"net/http"

	"example.com/url-shortner/config"
	"example.com/url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized !",
		})
		return
	}

	tokenUserId, tokenErr := utils.ValidateJWT(token)
	if tokenErr != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", tokenErr.Error()),
		})
		return
	}

	context.Set(config.JWT_LOGGED_IN_USER, tokenUserId) // Set it to be available for further requests
	context.Next()
}
