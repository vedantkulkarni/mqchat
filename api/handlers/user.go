package handlers

import (
	// "context"
	"fmt"
	"google.golang.org/grpc/codes"

	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/internal/app/protogen/proto"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	grpcClient     proto.UserGRPCServiceClient
	grpcConnClient proto.ConnectionGRPCServiceClient
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
	user.Put("/:uid", h.updateUser)
	user.Delete("/:uid", h.deleteUser)

	return nil
}

func (h *UserHandler) getUsers(c fiber.Ctx) error {
	fmt.Println("Received getUsers request")
	getConnectionsRequest := new(proto.GetConnectionsRequest)
	response, err := h.grpcConnClient.GetConnections(c.Context(), getConnectionsRequest)
	if err != nil {
		grpcErr := status.Convert(err)
		return fiber.NewError(fiber.StatusInternalServerError, grpcErr.Message())
	}

	fmt.Printf("Response from the server : %v", response)
	return c.JSON(response)

}

func (h *UserHandler) getUser(c fiber.Ctx) error {
	fmt.Println("Received get user request")
	getUserRequest := new(proto.GetUserRequest)
	getUserRequest.Id = c.Params("uid")
	response, err := h.grpcClient.GetUser(c.Context(), getUserRequest)
	if err != nil {
		grpcError := status.Convert(err)
		return fiber.NewError(fiber.StatusInternalServerError, grpcError.Message())
	}

	fmt.Printf("Response from the server : %v", response)
	return c.JSON(response)

}

func (h *UserHandler) createUser(c fiber.Ctx) error {
	fmt.Println("Received create user request")
	createUserRequest := new(proto.CreateUserRequest)
	e := c.Bind().Body(&createUserRequest)
	if e != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, e.Error())
	}
	_, err := h.grpcClient.CreateUser(c.Context(), createUserRequest)
	if err != nil {
		grpcError := status.Convert(err)
		return fiber.NewError(fiber.StatusConflict, grpcError.Message())
	}

	return nil

}

func (h *UserHandler) updateUser(c fiber.Ctx) error {
	fmt.Println("Received update user request")

	updateUserRequest := new(proto.UpdateUserRequest)
	if err := c.Bind().Body(updateUserRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Call gRPC method to update user
	_, err := h.grpcClient.UpdateUser(c.Context(), updateUserRequest)
	if err != nil {
		// Handle gRPC errors
		grpcError, ok := status.FromError(err)
		if ok && grpcError.Code() == codes.NotFound {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, grpcError.Message())
	}

	// Handle successful update response if needed
	fmt.Println("User updated successfully")

	// Return success response
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) deleteUser(c fiber.Ctx) error {
	fmt.Println("Received delete user request")

	userID := c.Params("user_id")

	// Call gRPC method to delete user
	_, err := h.grpcClient.DeleteUser(c.Context(), &proto.DeleteUserRequest{Id: userID})
	if err != nil {
		// Handle gRPC errors
		grpcError, ok := status.FromError(err)
		if ok && grpcError.Code() == codes.NotFound {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, grpcError.Message())
	}

	// Handle successful deletion response if needed
	fmt.Printf("User with ID %s deleted\n", userID)

	// Return success response
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("User with ID %s deleted successfully", userID),
	})
}
