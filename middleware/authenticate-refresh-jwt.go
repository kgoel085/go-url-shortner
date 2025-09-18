package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/model"
	"kgoel085.com/url-shortner/utils"
)

func AuthenticateRefreshToken(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized !",
		})
		return
	}

	decryptedToken, decryptErr := utils.Decrypt(token)
	if decryptErr != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", decryptErr.Error()),
		})
		return
	}

	tokenUserId, tokenErr := utils.ValidateJWT(decryptedToken, utils.RefreshJwtType)
	if tokenErr != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", tokenErr.Error()),
		})
		return
	}

	userRefreshToken, userRefreshTokenErr := model.GetRefreshTokenByToken(token)
	if userRefreshTokenErr != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", userRefreshTokenErr.Error()),
		})
		return
	}

	if time.Now().After(userRefreshToken.ExpiresAt) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Unauthorized - %s", errors.New("Refresh token has expired").Error()),
		})
		return
	}

	context.Set(config.JWT_LOGGED_IN_USER, tokenUserId) // Set it to be available for further requests
	context.Next()
}
