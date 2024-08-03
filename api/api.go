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
	httpAddr     string
	userGrpcAddr string
	roomGrpcAddr string
	chatGrpcAddr string
}

func NewAPI(httpAddr string, userAddr string, roomAddr string, chatAddr string) (*API, error) {

	return &API{
		httpAddr:     httpAddr,
		userGrpcAddr: userAddr,
		roomGrpcAddr: roomAddr,
		chatGrpcAddr: chatAddr,
	}, nil
}

func (a *API) Start() error {

	
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	user, err := grpc.NewClient(fmt.Sprintf("localhost:%s", a.userGrpcAddr), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return err
	}
	room, err := grpc.NewClient(fmt.Sprintf("localhost:%s", a.roomGrpcAddr), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return err
	}

	chat, err := grpc.NewClient(fmt.Sprintf("localhost:%s", a.chatGrpcAddr), opts...)
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
	roomClient := proto.NewRoomGRPCServiceClient(room)
	userHandler := handlers.NewUserHandler(userClient, roomClient)
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

	handlers.NewAuthHandler(&userClient).RegisterAuthRoutes(auth)
	handlers.RegisterSessionRoutes(session)
	err = app.Listen(fmt.Sprintf(":%s", a.httpAddr))
	if err != nil {
		log.Fatalf("Error occurred while listening to port %v :  %v", a.httpAddr, err)
	}

	return nil
}
