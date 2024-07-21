package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	api "github.com/vedantkulkarni/mqchat/api"
	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"google.golang.org/grpc"

	chatService "github.com/vedantkulkarni/mqchat/services/chat_service"
	connService "github.com/vedantkulkarni/mqchat/services/connection_service"
	mqttservice "github.com/vedantkulkarni/mqchat/services/mqtt_service"
	userService "github.com/vedantkulkarni/mqchat/services/user_service"

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

	//Connect to the database
	db, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Error occurred while connecting to the database")
	}

	fmt.Println("Database connected successfully")

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("Error occurred while closing DB")
		}
	}(db.DB)

	//Listen to gRPC responses
	listener, err := net.Listen("tcp", "localhost:"+config.UserServicePort)
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
		listener.Close()
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		fmt.Println("Closed the listner")
		if err != nil {
			fmt.Println("Error occurred while closing the listener")
		}
	}(listener)

	//Listen to gRPC responses
	listenerConn, err := net.Listen("tcp", "localhost:"+config.ConnServicePort)
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
		return
	}

	defer func(listener net.Listener) {
		err := listenerConn.Close()
		if err != nil {
			println("Error occurred while closing the listener")
		}
	}(listenerConn)

	//Listen to gRPC responses
	listenerChat, err := net.Listen("tcp", "localhost:"+config.ChatServicePort)
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
		return
	}

	defer func(listener net.Listener) {
		err := listenerChat.Close()
		if err != nil {
			println("Error occurred while closing the listener")
		}
	}(listenerChat)

	//Initialize User Service
	userServer, err := userService.NewUserGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : User")
		return
	}
	go func() {
		err := userServer.StartService(listener)
		if err != nil {
			fmt.Println("Error occurred while starting the gRPC server : User")
		}
	}()

	//Initialize Connections Service
	connServer, err := connService.NewConnectionGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : Connection")
		return
	}

	go func() {
		err := connServer.StartService(listenerConn)
		if err != nil {
			fmt.Println("Error occurred while starting the gRPC server : Connection")
		}
	}()

	// Initialize Chat Service
	chatServer, err := chatService.NewChatGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : Chat")
		return
	}

	go func() {
		err := chatServer.StartService(listenerChat)
		if err != nil {
			fmt.Println("Error occurred while starting the gRPC server : Chat")
		}
	}()

	//Listen to Mqtt connections
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	chat, err := grpc.NewClient(fmt.Sprintf("localhost:%s", config.ChatServicePort), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
	}	
	chatClient := proto.NewChatServiceClient(chat)
	//TODO: Add Streaming client
	
	mqttServer := mqttservice.NewMQTTService(&chatClient, nil)
	mqttServer.Start(config.MQTTServicePort)

	//Initialize the REST API server
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
