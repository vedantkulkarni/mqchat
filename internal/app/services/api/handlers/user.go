package handlers

import (
	// "context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/internal/app/proto"
)

type UserHandler struct {
	grpcClient proto.UserGRPCServiceClient	
}

func NewUserHandler(client proto.UserGRPCServiceClient) *UserHandler {
	return &UserHandler{
		grpcClient: client,
	}
}

func (h *UserHandler) RegisterUserRoutes(user fiber.Router) error {
	user.Get("/", h.getUsers)
	user.Post("/", h.createUser)
	user.Get("/:uid", h.getUser)
	user.Put("/:uid", h.updeUser)
	user.Delete("/:uid", h.deleteUser)

	return nil
}



func(h *UserHandler) getUsers(c fiber.Ctx) error {

	return c.SendString("Get Users")

}

func (h *UserHandler) createUser(c fiber.Ctx) error {

	fmt.Println("Create User")
	fmt.Println(c.Request())
	// createUserRequest := proto.CreateUserRequest{
	// 	Username : c.Body()["username"],
	// 	Email : "vedantk60@gmail.com",
	// }

 	// response, err := h.grpcClient.CreateUser(context.Background(), &createUserRequest)
	// if err != nil {
	// 	fmt.Printf("Error occured while creating the user, Error : %v", err)
	// 	return err
	// }
	// fmt.Printf("Response from the server : %v", response)

	return nil 

}

func (h *UserHandler) getUser(c fiber.Ctx) error {
	return c.SendString("Get User")

}

func (h *UserHandler) updeUser(c fiber.Ctx) error {
	return c.SendString("Update User")

}

func (h *UserHandler) deleteUser(c fiber.Ctx) error {
	return c.SendString("Delete User")

}
