package main

import (
	"database/sql"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/pkg/logger"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"github.com/vedantkulkarni/mqchat/services/chat/controller"

)



func main() {

	l:= logger.Get()

	var wg sync.WaitGroup
	err := godotenv.Load(".env")
	if err != nil {
		l.Panic().Err(err).Msg("Error loading .env file")	
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		l.Err(err).Msg("Error occurred while creating the DB connection")
	}

	
	wg.Add(1)
	defer func(DB *sql.DB, wg *sync.WaitGroup) {
		defer wg.Done()
		err := DB.Close()
		if err != nil {
			l.Err(err).Msg("Error occurred while closing the DB connection")	
		}
	}(db.DB, &wg)

	chatServer, err := controller.NewChatGRPCServer(db)
	if err != nil {
		l.Err(err).Msg("Error occurred while creating the chat server")
		return
	}

	chatServicePort := utils.GetEnvVarInt("CHAT_SERVICE_GRPC_PORT", 8002)
	chatServiceHost := utils.GetEnvVar("CHAT_SERVICE_GRPC_HOST", "service")

	


	wg.Add(1)	
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := chatServer.StartService(chatServicePort, chatServiceHost)
		if err != nil {
			l.Err(err).Msg("Error occurred while starting the chat server")
		}
	}(&wg)


	wg.Wait()

	l.Info().Msg("Chat service stopped")

}
