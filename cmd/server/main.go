package main

import (
	"fmt"
	"net/http"
	"github.com/vednatkulkarni/mqchat/internal/app/services/api"
	"github.com/vedantkulkarni/mqchat/internal/common/database"
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

	//Create a new server
	api:= api.NewAPI("localhost:8080")
	err := api.Start()
	if err != nil {
		fmt.Println("Error occured while starting the server")
	}	

}
