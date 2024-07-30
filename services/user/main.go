package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vedantkulkarni/mqchat/database"
	user "github.com/vedantkulkarni/mqchat/services/user/controller"
)

func main() {

	var block chan bool  

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

	userServer, err := user.NewUserGRPCServer(db)
	if err != nil {
		fmt.Println("Error occurred while creating the gRPC server : User")
		return
	}
	go func() {
		err := userServer.StartService()
		if err != nil {
			fmt.Println("Error occurred while starting the gRPC server : User")
		}
	}()

	<-block


}
