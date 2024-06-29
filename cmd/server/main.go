package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	api "github.com/vedantkulkarni/mqchat/api"
	"github.com/vedantkulkarni/mqchat/database"
	connService "github.com/vedantkulkarni/mqchat/internal/app/services/connection_service"
	userService "github.com/vedantkulkarni/mqchat/internal/app/services/user_service"
	util "github.com/vedantkulkarni/mqchat/internal/common"
	"net"
	"os"
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
	HttpPort string
	GrpcPort string
}

func getServerConfig(s *ServerConfig) *ServerConfig {
	s.HttpPort = util.GetEnvVarInt("HTTP_PORT", 8080)
	s.GrpcPort = util.GetEnvVarInt("GRPC_PORT", 2000)

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
	listener, err := net.Listen("tcp", "localhost:"+config.GrpcPort)
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			println("Error occurred while closing the listener")
		}
	}(listener)

	//Initialize User Service
	userServer, err := userService.NewUserGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : User")
		return
	}
	go userServer.StartService(listener)

	//Initialize Connections Service
	connServer, err := connService.NewConnectionGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : Connection")
		return
	}

	go connServer.StartService(listener)

	//Initialize the REST API server
	apiServer, err := api.NewAPI(config.HttpPort, config.GrpcPort)
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
