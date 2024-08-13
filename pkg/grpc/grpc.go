package grpc

import (
	"fmt"

	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/logger"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetChatClientConn(port string) proto.ChatServiceClient {
	l := logger.Get()

	host := utils.GetEnvVar("CHAT_SERVICE_GRPC_HOST", "localhost")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	chat, err := grpc.NewClient(fmt.Sprintf("%v:%s", host, port), opts...)
	if err != nil {
		l.Panic().Err(err).Msg("Error creating gRPC client for chat service")
	}

	if err != nil {
		l.Panic().Err(err).Msg("Error checking health of chat service")
	}

	return proto.NewChatServiceClient(chat)
}

func GetUserClientConn(port string) proto.UserGRPCServiceClient {
	l := logger.Get()
	
	host := utils.GetEnvVar("USER_SERVICE_GRPC_HOST", "localhost")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	user, err := grpc.NewClient(fmt.Sprintf("%v:%s", host, port), opts...)
	if err != nil {
		l.Panic().Err(err).Msg("Error creating gRPC client for user service")
	}

	if err != nil {

		l.Panic().Err(err).Msg("Error checking health of user service")
	}

	return proto.NewUserGRPCServiceClient(user)

}

func GetRoomClientConn(port string) proto.RoomGRPCServiceClient {
	l := logger.Get()

	host := utils.GetEnvVar("ROOM_SERVICE_GRPC_HOST", "localhost")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	room, err := grpc.NewClient(fmt.Sprintf("%v:%s", host, port), opts...)
	if err != nil {
		l.Panic().Err(err).Msg("Error creating gRPC client for room service")
	}

	if err != nil {
		l.Panic().Err(err).Msg("Error checking health of room service")
	}

	return proto.NewRoomGRPCServiceClient(room)
}


