package routes

import (
	"example.com/url-shortner/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(server *gin.Engine) {
	// Initialize global rate limiter middleware
	server.Use(middleware.GlobalRateLimit)

	// Initialize all routes
	AppRoutes(server.Group("/app"))
	UserRoutes(server.Group("/user"))
	OtpRoutes(server.Group("/otp"))
	UrlShorterRoutes(server.Group("/"))
}
