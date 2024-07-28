package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/services/chat/controller"
)

func main() {

	var blocker chan bool

	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Error occurred while connecting to the database")
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("Error occurred while closing DB")
		}
	}(db.DB)

	chatServer, err := controller.NewChatGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : Chat")
		return
	}

	go func() {
		err := chatServer.StartService("2200")
		if err != nil {
			fmt.Println("Error occurred while starting the gRPC server : Chat")
		}
	}()

	fmt.Println("Chat server started successfully!")
	fmt.Println("Blocking the server")

	<-blocker

}