package routes

import (
	"net/http"

	"example.com/url-shortner/model"
	"github.com/gin-gonic/gin"
)

func AppRoutes(router *gin.RouterGroup) {
	router.GET("/ping", handlePing)
}

// @Summary      Ping
// @Description  Health check endpoint. Returns "pong" if the server is running.
// @Tags         App
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.APIResponse "Success" "Example: {\"message\": \"pong\"}"
// @Router       /ping [get]
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, model.APIResponse{
		Message: "pong",
	})
}
