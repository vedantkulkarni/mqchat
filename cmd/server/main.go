package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vedantkulkarni/mqchat/internal/common/database"
	"github.com/vedantkulkarni/mqchat/internal/routes"
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
	server := gin.Default()

	//Register the routes
	routes.RegisterRoutes(server)

	//Connect to database.
	pgdb, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Error occured")
	}
	fmt.Printf("%v\n", pgdb)

	//test endpoint
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	server.Run()

	fmt.Println("Server Started !")
}
