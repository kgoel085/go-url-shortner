package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppRoutes(router *gin.RouterGroup) {
	router.GET("/ping", handlePing)
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
