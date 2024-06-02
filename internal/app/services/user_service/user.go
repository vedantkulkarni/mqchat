package usersservice

import (
	"github.com/vedantkulkarni/mqchat/internal/common/database"
)


type UserGRPCServer struct {
	db *database.PostgresDB	
}

func  NewUserGRPCServer(db *database.PostgresDB) (*UserGRPCServer, error) {
	return &UserGRPCServer{
		db: db,
	}, nil
}
