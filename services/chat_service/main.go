package chatservice

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
)

type ChatGRPCServer struct {
	DB *sql.DB
	proto.UnimplementedChatServiceServer
}

func (g *ChatGRPCServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	//Store chat message in db
	response := &models.Chat{
		ToUserID:   int(req.Message.UserId_2),
		FromUserID: int(req.Message.UserId_1),
		Message:    req.Message.Content,
	}

	err := response.Insert(ctx, g.DB, boil.Infer())
	if err != nil {
		fmt.Println("Error occured while inserting chat message in db")
		return nil, err
	}

	fmt.Println("Chat message inserted successfully in db")

	//Return response
	return &proto.SendMessageResponse{
		Message: req.Message,
	}, nil

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

	fmt.Println("gRPC chat server registered successfully")
	if err := server.Serve(listner); err != nil {
		fmt.Println("Error occured while serving the chat gRPC server")
		return err
	}
	fmt.Println("gRPC chat server started successfully!")
	return nil
}
