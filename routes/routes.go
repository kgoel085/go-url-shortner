package routes

import (
	"example.com/url-shortner/config"
	"example.com/url-shortner/docs"
	"example.com/url-shortner/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRouter(server *gin.Engine) {
	setUpSwagger(server)

	// Initialize global rate limiter middleware
	server.Use(middleware.GlobalRateLimit)

	// Initialize all routes
	AppRoutes(server.Group("/app"))
	UserRoutes(server.Group("/user"))
	OtpRoutes(server.Group("/otp"))
	UrlShorterRoutes(server.Group("/"))
}

func setUpSwagger(server *gin.Engine) {
	docs.SwaggerInfo.Title = config.Config.APP.Name
	docs.SwaggerInfo.Description = "A robust URL shortener service built with Go, Gin, PostgreSQL, and Redis. This project provides a scalable API for shortening URLs, tracking usage, and managing links securely."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.Config.APP.Host + ":" + config.Config.APP.Port
	docs.SwaggerInfo.BasePath = "/"
	// Dynamically set schemes based on config or environment
	if config.Config.APP.EnableHTTPS {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	// Swagger docs route
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
