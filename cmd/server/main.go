package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	api "github.com/vedantkulkarni/mqchat/api"
	"github.com/vedantkulkarni/mqchat/services/mqtt"

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
	RoomServicePort string
	ChatServicePort string
	MQTTServicePort string
}

func getServerConfig(s *ServerConfig) *ServerConfig {
	s.HttpPort = util.GetEnvVarInt("HTTP_PORT", 8080)
	// s.UserServicePort = util.GetEnvVarInt("USER_SERVICE_GRPC_PORT", 2000)
	// s.RoomServicePort = util.GetEnvVarInt("ROOM_SERVICE_GRPC_PORT", 2100)
	// s.ChatServicePort = util.GetEnvVarInt("CHAT_SERVICE_GRPC_PORT", 2200)
	// s.MQTTServicePort = util.GetEnvVarInt("CHAT_SERVICE_MQTT_PORT", 2300)

	// Hardcoded for now
	s.UserServicePort = "2000"
	s.RoomServicePort = "2100"
	s.ChatServicePort = "2200"
	s.MQTTServicePort = "2300"
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

	// MQTT Server
	mqttServer := mqtt.NewMQTTService()
	go func() {
		mqttServer.Start(config.MQTTServicePort)
	}()

	// REST API Server
	apiServer, err := api.NewAPI(config.HttpPort, config.UserServicePort, config.RoomServicePort, config.ChatServicePort)
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
