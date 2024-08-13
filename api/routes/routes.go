package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/api/handlers"
	"github.com/vedantkulkarni/mqchat/api/middleware"

	"github.com/vedantkulkarni/mqchat/pkg/config"
	grpcUtils "github.com/vedantkulkarni/mqchat/pkg/grpc"
)

func Init(baseRouter fiber.Router) {
	config := config.Get()

	// Group the routes
	userRoutes := baseRouter.Group("/users")
	auth := baseRouter.Group("/auth")
	chatRoutes := baseRouter.Group("/chat")

	// User Routes
	userHandler := handlers.NewUserHandler(grpcUtils.GetUserClientConn(config.HttpPort), grpcUtils.GetRoomClientConn(config.RoomServicePort))
	RegisterUserRoutes(userRoutes, userHandler)

	// Chat Routes
	chatHandler := handlers.NewChatHandler(grpcUtils.GetChatClientConn(config.ChatServicePort))
	chatHandler.RegisterChatRoutes(chatRoutes)

	userService := grpcUtils.GetUserClientConn(config.UserServicePort)
	handlers.NewAuthHandler(&userService).RegisterAuthRoutes(auth)
}

func RegisterUserRoutes(user fiber.Router, h *handlers.UserHandler) {
	user.Get("/", h.GetUsers, middleware.AuthMiddleware)
	user.Post("/", h.CreateUser) // Used for Signup
	user.Put("/:uid", h.UpdateUser, middleware.AuthMiddleware)
	user.Delete("/:uid", h.DeleteUser, middleware.AuthMiddleware)
}

func RegisterChatRoutes(chat fiber.Router, h *handlers.ChatHandler) {
	chat.Get("/", h.GetMessages)
	chat.Post("/", h.SendMessage)
}
