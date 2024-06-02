package main

import (
	"fmt"
	"net"
	"github.com/vedantkulkarni/mqchat/internal/app/services/api"
	usersservice "github.com/vedantkulkarni/mqchat/internal/app/services/user_service"
	"github.com/vedantkulkarni/mqchat/internal/app/services/user_service/database"
)

var (
	version   string
	buildDate string
)

const (
	appName         = "mqchat"
	friendlyAppName = "MQTT-GRPC Chat Application"
)

func main() {

	//Connect to the database
	db, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Error occured while connecting to the database")
	}


	
	//Initialize the REST API server
	apiServer, err := api.NewAPI(":8080") 
	if err != nil {
		fmt.Println("Error occured while creating the server")
	}

	err = apiServer.Start()
	if err != nil {
		fmt.Println("Error occured while starting the server")
	}	


	//Listen to gRPC and REST API requests
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error occured while listening to the port")	
	}

	//Initialize the gRPC servers
	
	//Initialize User Service
	userServer, err := usersservice.NewUserGRPCServer(db)
	if err != nil {
		fmt.Println("Error occured while creating the gRPC server")
	}
	err = userServer.StartService(listner)
	if err != nil {
		fmt.Println("Error occured while starting the gRPC server")
	}
	
}
