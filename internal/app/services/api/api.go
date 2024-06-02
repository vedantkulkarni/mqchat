package api

import (
	"fmt"

	"github.com/vedantkulkarni/mqchat/internal/app/proto"
	"github.com/vedantkulkarni/mqchat/internal/app/services/api/handlers"
	"google.golang.org/grpc"

	"github.com/gofiber/fiber/v3"
)

type API struct {
	addr string
	grpcConn *grpc.ClientConn
}

func NewAPI(connStr string) (*API, error) {
	//Connect to the gRPC server
	var opts []grpc.DialOption
	conn, err := grpc.NewClient(connStr, opts...)
	if err != nil {
		fmt.Println("Error occured while connecting to the gRPC server")
		return nil, err
	}
	defer conn.Close()	


	return &API{
		addr: connStr,
		grpcConn: conn,
	}, nil
}

func (a *API) Start() error {
	app := fiber.New()
	api := app.Group("/api")

	v1 := api.Group("/v1")

	// Group the routes
	userRoutes := v1.Group("/users")
	auth := v1.Group("/auth") 
	chat := v1.Group("/chat")
	session := v1.Group("/session")

	// Register the routes

	// User Routes
	client := proto.NewUserGRPCServiceClient(a.grpcConn)
	userHandler := handlers.NewUserHandler(client)
	err := userHandler.RegisterUserRoutes(userRoutes)
	if err != nil {
		fmt.Println("Error occured while registering the user routes")
		return err
	}
	
	
	
	handlers.RegisterAuthRoutes(auth)
	handlers.RegisterChatRoutes(chat)
	handlers.RegisterSessionRoutes(session)


	return nil
}
