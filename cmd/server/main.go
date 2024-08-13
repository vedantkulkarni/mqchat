package main

import (
	"sync"

	"github.com/joho/godotenv"
	api "github.com/vedantkulkarni/mqchat/api"
	"github.com/vedantkulkarni/mqchat/pkg/config"
	"github.com/vedantkulkarni/mqchat/services/mqtt"

	"github.com/vedantkulkarni/mqchat/pkg/logger"
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

	var wg sync.WaitGroup

	// Set up logger
	l := logger.Get()

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		l.Panic().Err(err).Msg("Error loading .env file")
	}

	//Get configurations
	config := config.Get()

	// MQTT Server
	wg.Add(1)
	mqttServer := mqtt.NewMQTTService()
	go mqttServer.Start(config.MQTTServicePort, &wg)

	// REST API Server
	apiServer, err := api.NewAPI(config.HttpPort)
	if err != nil {
		l.Panic().Err(err).Msg("Error creating API server")
	}
	wg.Add(1)
	go apiServer.Start(&wg)

	l.Info().Msg("API Service started successfully")

	wg.Wait()
}
