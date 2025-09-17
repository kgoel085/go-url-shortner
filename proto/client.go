package proto

import (
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/utils"
)

var ClientManager *Manager

type GRPCClientType string

const (
	EmailServiceClientType GRPCClientType = "email"
)

func InitClients() {
	ClientManager = NewManager()

	if config.Config.GRPC.EmailServiceAddr != "" { // Email GRPC service
		utils.Log.Info("GRPC:: connecting to email service at ", config.Config.GRPC.EmailServiceAddr)
		_, err := ClientManager.Connect(string(EmailServiceClientType), config.Config.GRPC.EmailServiceAddr)
		if err != nil {
			utils.Log.Errorf("GRPC:: failed to connect to email service: %v", err)
		} else {
			utils.Log.Info("GRPC:: connected to email service")
		}
	}
}
