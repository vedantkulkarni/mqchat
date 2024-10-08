package user

import (
	"fmt"
	"log"
	"net"

	database "github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	env "github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc"
)

func NewUserGRPCServer(db *database.DbInterface) (*UserGRPCServer, error) {
	//This microservice internally connects to the 'ROOMs' microservice
	host := env.GetEnvVar("ROOMS_SERVICE_GRPC_HOST", "service")
	roomConnPort := env.GetEnvVarInt("ROOMS_SERVICE_GRPC_PORT", 2100)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	roomConn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, roomConnPort), opts...)
	if err != nil {
		log.Fatalf("Error occurred while connecting to the gRPC server : %v", err)
		return nil, err
	}

	return &UserGRPCServer{
		DB:             db.DB,
		RoomGRPCClient: proto.NewRoomGRPCServiceClient(roomConn),
	}, nil
}

func (u *UserGRPCServer) StartService() error {

	// host:= "user-service"
	userServicePort := env.GetEnvVarInt("USER_SERVICE_GRPC_PORT", 8003)
	userServiceHost := env.GetEnvVar("USER_SERVICE_GRPC_HOST", "service")
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", userServiceHost, userServicePort))
	if err != nil {
		log.Panic("user service port err:", err)
		listener.Close()
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		fmt.Println("Closed the listner")
		if err != nil {
			fmt.Println("Error occurred while closing the listener")
		}
	}(listener)

	g := grpc.NewServer()
	fmt.Println("Starting gRPC user server on port : ", listener.Addr().String())

	proto.RegisterUserGRPCServiceServer(g, u)

	if err := g.Serve(listener); err != nil {
		log.Fatalf("Error occured while serving the gRPC server : %v", err)
		return err
	}

	return nil
}
