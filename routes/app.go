package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/model"
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
// @Router       /app/ping [get]
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, model.APIResponse{
		Message: "pong",
	})
}
