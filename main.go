package main

import (
	"fmt"

	"example.com/url-shortner/config"
	"example.com/url-shortner/db"
	"example.com/url-shortner/routes"
	"example.com/url-shortner/utils"
	"example.com/url-shortner/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	utils.InitLogger()

	config.LoadConfig()
	db.InitDB()
	validator.LoadCustomBindings()
	routes.SetUpRouter(server)

	appUrl := fmt.Sprintf("%s:%s", config.Config.App.Host, config.Config.App.Port)

	utils.Log.Info("Starting server...")
	server.Run(appUrl)
}
