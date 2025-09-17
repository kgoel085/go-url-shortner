package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/db"
	"kgoel085.com/url-shortner/proto"
	"kgoel085.com/url-shortner/routes"
	"kgoel085.com/url-shortner/utils"
	"kgoel085.com/url-shortner/validator"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server := gin.Default()

	utils.InitLogger()             // Initialize logger
	config.LoadConfig()            // Load ENV variables
	db.InitRedis()                 // Initialize Redis client
	db.InitDB()                    // Initialize Postgres client
	validator.LoadCustomBindings() // Load custom validators
	proto.InitClients()            // Initialize gRPC clients
	routes.SetUpRouter(server)     // Setup all routes

	appUrl := fmt.Sprintf("%s:%s", config.Config.APP.Host, config.Config.APP.Port)
	trustedProxies := strings.Split(config.Config.APP.TrustedProxies, ",")
	if len(trustedProxies) > 0 {

		utils.Log.Info("Trusted proxies list: ", trustedProxies)
		server.SetTrustedProxies(trustedProxies)
	}

	utils.Log.Info("Starting server...")
	serverErr := server.Run(appUrl)

	if serverErr != nil {
		utils.Log.Fatal("Failed to start server: ", serverErr)
	}
	utils.Log.Info("Server started at ", appUrl)
}
