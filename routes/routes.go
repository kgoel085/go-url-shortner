package routes

import "github.com/gin-gonic/gin"

func SetUpRouter(server *gin.Engine) {
	AppRoutes(server.Group("/app"))
	UserRoutes(server.Group("/user"))
	OtpRoutes(server.Group("/otp"))
	UrlShorterRoutes(server.Group("/"))
}
