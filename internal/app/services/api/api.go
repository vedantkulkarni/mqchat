package api

import (
	"github.com/vedantkulkarni/mqchat/internal/app/services/api/handlers"

	"github.com/gofiber/fiber/v3"
)

type API struct {
	addr string

	// Add gRPC client here
}

func NewAPI(connStr string) (*API, error) {
	return &API{
		addr: connStr,
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
	handlers.RegisterUserRoutes(userRoutes)
	handlers.RegisterAuthRoutes(auth)
	handlers.RegisterChatRoutes(chat)
	handlers.RegisterSessionRoutes(session)


	return nil
}
