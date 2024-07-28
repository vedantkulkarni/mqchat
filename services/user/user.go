package user

import (
	"fmt"
	"log"
	"net"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	env "github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc"
)

func NewUserGRPCServer(db *database.DbInterface) (*UserGRPCServer, error) {
	//This microservice internally connects to the 'connections' microservice
	connectionServicePort := env.GetEnvVarInt("CONNECTION_SERVICE_GRPC_PORT", 2100)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", connectionServicePort), opts...)
	if err != nil {
		log.Fatalf("Error occurred while connecting to the gRPC server : %v", err)
		return nil, err
	}

	return &UserGRPCServer{
		DB:             db.DB,
		ConnGrpcClient: proto.NewConnectionGRPCServiceClient(conn),
	}, nil
}

func (u *UserGRPCServer) StartService(port string) error {
	var block chan struct{}
	
	listener, err := net.Listen("tcp", "localhost:"+port)
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
	fmt.Println("Starting gRPC user server")

	proto.RegisterUserGRPCServiceServer(g, u)

	fmt.Println("gRPC user server registered successfully")
	if err := g.Serve(listener); err != nil {
		log.Fatalf("Error occured while serving the gRPC server : %v", err)
		return err
	}
	fmt.Println("gRPC user server started successfully!")

	<-block

	return nil
}
