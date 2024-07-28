package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	api "github.com/vedantkulkarni/mqchat/api"

	util "github.com/vedantkulkarni/mqchat/pkg/utils"
)

var (
	version   string
	buildDate string
)

const (
	appName         = "mqchat"
	friendlyAppName = "MQTT-GRPC Chat Application"
)

type ServerConfig struct {
	HttpPort        string
	UserServicePort string
	ConnServicePort string
	ChatServicePort string
	MQTTServicePort string
}

func getServerConfig(s *ServerConfig) *ServerConfig {
	s.HttpPort = util.GetEnvVarInt("HTTP_PORT", 8080)
	s.UserServicePort = util.GetEnvVarInt("USER_SERVICE_GRPC_PORT", 2000)
	s.ConnServicePort = util.GetEnvVarInt("CONNECTION_SERVICE_GRPC_PORT", 2100)
	s.ChatServicePort = util.GetEnvVarInt("CHAT_SERVICE_GRPC_PORT", 2200)
	s.MQTTServicePort = util.GetEnvVarInt("CHAT_SERVICE_MQTT_PORT", 2300)
	return s
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}

	//Get configurations
	config := &ServerConfig{}
	config = getServerConfig(config)
	// connServer, err := connection.NewConnectionGRPCServer(db)
	// if err != nil {
	// 	fmt.Println("Error occurred while creating the gRPC server : Connection")
	// 	return
	// }

	// go func() {
	// 	err := connServer.StartService(config.ConnServicePort)
	// 	if err != nil {
	// 		fmt.Println("Error occurred while starting the gRPC server : Connection")
	// 	}
	// }()

	// mqttServer := mqtt.NewMQTTService()
	// mqttServer.Start(config.MQTTServicePort)

	// REST API Server
	apiServer, err := api.NewAPI(config.HttpPort, config.UserServicePort, config.ConnServicePort)
	if err != nil {
		fmt.Println("Error occured while creating the server")
	}

	err = apiServer.Start()
	if err != nil {
		fmt.Println("Error occured while starting the server")
		return
	}

	fmt.Println("First Level of Server Started Successfully!")

}
