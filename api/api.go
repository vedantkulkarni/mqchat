package api

import (
	"fmt"
	"log"

	"github.com/vedantkulkarni/mqchat/api/handlers"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"google.golang.org/grpc"

	"github.com/gofiber/fiber/v3"
)

type API struct {
	addr     string
	grpcAddr string
	connGrpcAddr string
	chatGrpcAddr string
}

func NewAPI(addr string, grpcAddr string, connGrpcAddr string) (*API, error) {

	return &API{
		addr:     addr,
		grpcAddr: grpcAddr,
		connGrpcAddr: connGrpcAddr,
	}, nil
}

func (a *API) Start() error {

	//Connect to the gRPC server
	// conn, err :
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	user, err := grpc.NewClient(fmt.Sprintf("localhost:%s", "2000"), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return err
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", a.connGrpcAddr), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return err
	}

	chat, err := grpc.NewClient(fmt.Sprintf("localhost:%s", "2200"), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return err
	}	

	app := fiber.New()
	api := app.Group("/api")

	v1 := api.Group("/v1")

	// Group the routes
	userRoutes := v1.Group("/users")
	auth := v1.Group("/auth")
	chatRoutes := v1.Group("/chat")
	session := v1.Group("/session")

	// User Routes
	userClient := proto.NewUserGRPCServiceClient(user)
	connClient := proto.NewConnectionGRPCServiceClient(conn)
	userHandler := handlers.NewUserHandler(userClient, connClient)
	err = userHandler.RegisterUserRoutes(userRoutes)
	if err != nil {
		fmt.Println("Error occured while registering the user routes")
		return err
	}
	fmt.Println("User routes registered successfully")

	//Chat Routes
	chatClient := proto.NewChatServiceClient(chat)
	chatHandler := handlers.NewChatHandler(chatClient)
	err = chatHandler.RegisterChatRoutes(chatRoutes)
	if err != nil {
		fmt.Println("Error occured while registering the chat routes")
		return err
	}
	fmt.Println("Chat routes registered successfully")

	handlers.RegisterAuthRoutes(auth)
	// handlers.RegisterChatRoutes(chat)
	handlers.RegisterSessionRoutes(session)
	err = app.Listen(fmt.Sprintf(":%s", a.addr))
	if err != nil {
		log.Fatalf("Error occurred while listening to port %v :  %v", a.addr, err)
	}

	return nil
}
