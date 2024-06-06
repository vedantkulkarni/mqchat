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
	grpcAddr string 
}

func NewAPI(addr string, grpcAddr string) (*API, error) {
	

	return &API{
		addr: addr,
		grpcAddr: grpcAddr,
	}, nil
}

func (a *API) Start() error {

	//Connect to the gRPC server
	// conn, err :
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1%s",a.grpcAddr), opts...)
	if err != nil {
		fmt.Println("Error occured while connecting to the gRPC server")
		return  err
	}

	fmt.Println("grpc dial successful")
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
	client := proto.NewUserGRPCServiceClient(conn)
	userHandler := handlers.NewUserHandler(client)
	err = userHandler.RegisterUserRoutes(userRoutes)
	if err != nil {
		fmt.Println("Error occured while registering the user routes")
		return err
	}
	fmt.Println("User routes registered successfully")
	
	handlers.RegisterAuthRoutes(auth)
	handlers.RegisterChatRoutes(chat)
	handlers.RegisterSessionRoutes(session)
	app.Listen(a.addr)	


	return nil
}
