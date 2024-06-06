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
	defer db.DB.Close()
	fmt.Println("Database connected successfully")

	//Listen to gRPC requests
	listner, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
		return
	}

	//Initialize the gRPC servers

	//Initialize User Service
	userServer, err := usersservice.NewUserGRPCServer(db)
	if err != nil {
		fmt.Println("Error occured while creating the gRPC server")
		return
	}
	go userServer.StartService(listner)
	// if err != nil {
	// 	fmt.Println("Error occured while starting the gRPC server")
	// 	return
	// }
	// fmt.Println("gRPC User microservice started successfully")

	//Initialize the REST API server
	apiServer, err := api.NewAPI(":8080", ":2000")
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
