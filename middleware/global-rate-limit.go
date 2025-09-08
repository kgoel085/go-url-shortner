package middleware

import (
	"net/http"
	"time"

	"example.com/url-shortner/config"
	"example.com/url-shortner/db"
	"github.com/gin-gonic/gin"
)

func GlobalRateLimit(context *gin.Context) {
	clientKey := prepareClientKey(context)

	// If the request exceeds the rate limit, abort the request
	if ok, _ := db.CheckRateLimitInTimeUnit(context, clientKey, 20, time.Minute); !ok {
		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too Many Requests"})
		return
	}

	context.Next()
}

func prepareClientKey(context *gin.Context) string {
	// You can customize this function to create a more specific key based on your requirements
	// For example, you might want to include user ID if available
	userID := getUserIDFromContext(context) // implement this as needed
	if userID != 0 {
		return "global:user:" + string(rune(userID))
	}
	return "global:" + context.ClientIP()
}

func getUserIDFromContext(context *gin.Context) int64 {
	// Implement logic to extract user ID from context, e.g., from JWT claims or session
	// Return 0 if the user is not authenticated
	if userID, exists := context.Get(config.JWT_LOGGED_IN_USER); exists {
		if id, ok := userID.(int64); ok {
			return id
		}
	}
	return 0
}
