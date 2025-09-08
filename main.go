package main

import (
	"fmt"
	"strings"

	"example.com/url-shortner/config"
	"example.com/url-shortner/db"
	"example.com/url-shortner/routes"
	"example.com/url-shortner/utils"
	"example.com/url-shortner/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	utils.InitLogger()             // Initialize logger
	config.LoadConfig()            // Load ENV variables
	db.InitRedis()                 // Initialize Redis client
	db.InitDB()                    // Initialize Postgres client
	validator.LoadCustomBindings() // Load custom validators
	routes.SetUpRouter(server)     // Setup all routes

	appUrl := fmt.Sprintf("%s:%s", config.Config.APP.Host, config.Config.APP.Port)
	trustedProxies := strings.Split(config.Config.APP.TrustedProxies, ",")
	if len(trustedProxies) > 0 {

		utils.Log.Info("Trusted proxies list: ", trustedProxies)
		server.SetTrustedProxies(trustedProxies)
	}

	utils.Log.Info("Starting server...")
	server.Run(appUrl)
}
