package handlers

import (
	// "context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"

	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/api/middlewares"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	jsonUtils "github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	grpcUserClient proto.UserGRPCServiceClient
	grpcConnClient proto.ConnectionGRPCServiceClient
}

func NewUserHandler(client proto.UserGRPCServiceClient, connection proto.ConnectionGRPCServiceClient) *UserHandler {
	return &UserHandler{
		grpcUserClient: client,
		grpcConnClient: connection,
	}
}

func (h *UserHandler) RegisterUserRoutes(user fiber.Router) error {
	user.Get("/", h.getUsers, api.AuthMiddleware)
	user.Post("/", h.createUser) // Used for Signup
	user.Put("/:uid", h.updateUser, api.AuthMiddleware)
	user.Delete("/:uid", h.deleteUser, api.AuthMiddleware)

	return nil
}

func (h *UserHandler) getUsers(c fiber.Ctx) error {
	fmt.Println("Received getUsers request")
	uid := c.Query("uid")
	filter := c.Query("filter")

	if filter == "user" {
		h.getUser(c)
	} else if filter == "connection" {
		uid, err := strconv.Atoi(uid)
		if err != nil {
			return jsonUtils.WriteJson(fiber.StatusBadRequest, nil, &jsonUtils.ApiError{
				Code:    1,
				Message: err.Error(),
				Details: err.Error(),
			}, c)
		}

		getUsersRequest := &proto.GetUsersRequest{
			Id: int64(uid),
		}
		// fmt.Printf("Request received : %v", getConnectionsRequest)
		response, err := h.grpcUserClient.GetUsers(c.Context(), getUsersRequest)
		if err != nil {
			grpcErr := status.Convert(err)
			return jsonUtils.WriteJson(fiber.StatusInternalServerError, nil, &jsonUtils.ApiError{
				Code:    int(grpcErr.Code()),
				Message: grpcErr.Message(),
				Details: grpcErr.Message(),
			}, c)
		}

		if response != nil {
			return jsonUtils.WriteJson(200, response, nil, c)
		}

	} else {
		c.JSON(fiber.Map{
			"message": "Invalid filter",
		})
	}

	return nil

}

func (h *UserHandler) getUser(c fiber.Ctx) error {
	fmt.Println("Received get user request")
	getUserRequest := new(proto.GetUserRequest)
	uid, err := strconv.Atoi(c.Query("uid"))
	if err != nil {
		return jsonUtils.WriteJson(fiber.StatusBadRequest, nil, &jsonUtils.ApiError{
			Code:    1,
			Message: err.Error(),
			Details: err.Error(),
		}, c)
	}
	getUserRequest.Id = int64(uid)
	response, err := h.grpcUserClient.GetUser(c.Context(), getUserRequest)
	if err != nil {
		grpcError := status.Convert(err)
		return jsonUtils.WriteJson(fiber.StatusInternalServerError, nil, &jsonUtils.ApiError{
			Code:    int(grpcError.Code()),
			Message: grpcError.Message(),
			Details: grpcError.Message(),
		}, c)
	}

	fmt.Printf("Response from the server : %v", response)
	return jsonUtils.WriteJson(200, response, nil, c)

}

func (h *UserHandler) createUser(c fiber.Ctx) error {
	fmt.Println("Received create user request")
	createUserRequest := new(proto.UpdateUserRequest)
	e := c.Bind().Body(&createUserRequest)
	if e != nil {
		return jsonUtils.WriteJson(fiber.StatusBadRequest, nil, &jsonUtils.ApiError{
			Code:    1,
			Message: e.Error(),
			Details: e.Error(),
		}, c)
	}

	createUserRequest.IsCreate = true
	fmt.Printf("Request received : %v", createUserRequest)
	fmt.Println("Making request to the gRPC User service")
	response, err := h.grpcUserClient.UpdateUser(c.Context(), createUserRequest)
	fmt.Println("Response from the gRPC User service %v", response)
	fmt.Println("Error from the gRPC User service %v", err)
	if err != nil {
		grpcError := status.Convert(err)
		return jsonUtils.WriteJson(fiber.StatusInternalServerError, nil, &jsonUtils.ApiError{
			Code:    int(grpcError.Code()),
			Message: grpcError.Message(),
			Details: grpcError.Message(),
		}, c)
	}

	return jsonUtils.WriteJson(200, response, nil, c)

}

func (h *UserHandler) updateUser(c fiber.Ctx) error {
	fmt.Println("Received update user request")

	updateUserRequest := new(proto.UpdateUserRequest)
	if err := c.Bind().Body(updateUserRequest); err != nil {
		return jsonUtils.WriteJson(
			fiber.ErrBadRequest.Code,
			nil,
			jsonUtils.BadRequestApiError, c)
	}

	updateUserRequest.IsCreate = false

	// Call gRPC method to update user
	_, err := h.grpcUserClient.UpdateUser(c.Context(), updateUserRequest)
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
	uid, err := strconv.Atoi(c.Params("uid"))
	if err != nil {
		jsonUtils.WriteJson(
			fiber.ErrBadRequest.Code,
			nil,
			jsonUtils.BadRequestApiError,
			c,
		)
	}

	// Call gRPC method to delete user
	_, err = h.grpcUserClient.DeleteUser(c.Context(), &proto.DeleteUserRequest{Id: int64(uid)})
	if err != nil {
		// Handle gRPC errors
		grpcError, ok := status.FromError(err)
		if ok && grpcError.Code() == codes.NotFound {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, grpcError.Message())
	}

	// Handle successful deletion response if needed
	// fmt.Printf("User with ID %s deleted\n", uid)

	// Return success response
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("User with ID %s deleted successfully", uid),
	})
}
