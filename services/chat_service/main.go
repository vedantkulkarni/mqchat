package chatservice

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/internal/app/protogen/proto"
	"google.golang.org/grpc"
)

type ChatGRPCServer struct {
	DB *sql.DB
	proto.UnimplementedChatServiceServer
}

func (g *ChatGRPCServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	
	return nil, nil
}
func (g *ChatGRPCServer) GetMessages(req *proto.GetMessagesRequest, stream proto.ChatService_GetMessagesServer) error {
	return nil
}

func NewChatGRPCServer(db *database.DbInterface) (*ChatGRPCServer, error) {
	return &ChatGRPCServer{
		DB: db.DB,
	}, nil
}

func (g *ChatGRPCServer) StartService(listner net.Listener) error {
	server := grpc.NewServer()

	proto.RegisterChatServiceServer(server, g)

	fmt.Println("gRPC user server registered successfully")
	if err := server.Serve(listner); err != nil {
		fmt.Println("Error occured while serving the gRPC server")
		return err
	}
	fmt.Println("gRPC user server started successfully!")
	return nil
}
