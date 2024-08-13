package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"

	database "github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/protobuf/types/known/timestamppb"

	healthPb "google.golang.org/grpc/health/grpc_health_v1"
)

type ChatGRPCServer struct {
	DB *sql.DB
	proto.UnimplementedChatServiceServer
}

func (g *ChatGRPCServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	//Store chat message in db

	response := &models.Chat{
		UserID1: int(req.Message.UserId_2),
		UserID2: int(req.Message.UserId_1),
		Message: req.Message.Content,
		ChatID:  int(req.Message.RoomId),
	}

	err := response.Insert(ctx, g.DB, boil.Infer())
	if err != nil {
		fmt.Printf("Error occured while inserting chat message in db : %v \n", err)
		return nil, err
	}

	fmt.Println("Chat message inserted successfully in db")

	//Return response
	return &proto.SendMessageResponse{
		Message: req.Message,
	}, nil

}
func (g *ChatGRPCServer) GetMessages(req *proto.GetMessagesRequest, stream proto.ChatService_GetMessagesServer) error {

	//Get chat messages from db
	//TODO: Implement pagination
	chatId, err := strconv.Atoi(req.ChatId)
	chatMessages, err := models.Chats(qm.Where("chat_id=?", chatId), qm.Limit(50), qm.Offset(0)).All(context.Background(), g.DB)
	if err != nil {
		fmt.Println("Error occured while fetching chat messages from db", err)
		return err
	}

	//Send chat messages to client

	for _, chatMessage := range chatMessages {
		message := &proto.Message{
			UserId_1:  int64(chatMessage.UserID1),
			UserId_2:  int64(chatMessage.UserID2),
			RoomId:    int64(chatMessage.ChatID),
			Content:   chatMessage.Message,
			CreatedAt: timestamppb.New(chatMessage.CreatedAt.Time),
		}

		if err := stream.Send(message); err != nil {
			fmt.Println("Error occured while sending chat message to client", err)
			return err
		}
	}

	return nil
}

func NewChatGRPCServer(db *database.DbInterface) (*ChatGRPCServer, error) {
	return &ChatGRPCServer{
		DB: db.DB,
	}, nil
}

func (g *ChatGRPCServer) StartService(port string, host string) error {
	//Listen to gRPC responses
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			println("Error occurred while closing the listener")
		}
	}(listener)
	server := grpc.NewServer()
	healthCheck := health.NewServer()

	healthPb.RegisterHealthServer(server, healthCheck)
	proto.RegisterChatServiceServer(server, g)

	healthCheck.SetServingStatus("", healthPb.HealthCheckResponse_SERVING)

	fmt.Println("gRPC chat server registered successfully")
	if err := server.Serve(listener); err != nil {
		fmt.Println("Error occured while serving the chat gRPC server")
		healthCheck.SetServingStatus("", healthPb.HealthCheckResponse_NOT_SERVING)
		return err
	}
	fmt.Println("gRPC chat server started successfully!")

	return nil
}
