package handlers

import (
	// "context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"

	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	jsonUtils "github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	grpcUserClient proto.UserGRPCServiceClient
	grpcRoomClient proto.RoomGRPCServiceClient
}

func NewUserHandler(userClient proto.UserGRPCServiceClient, roomClient proto.RoomGRPCServiceClient) *UserHandler {
	return &UserHandler{
		grpcUserClient: userClient,
		grpcRoomClient: roomClient,
	}
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	fmt.Println("Received getUsers request")
	uid := c.Query("uid")
	filter := c.Query("filter")

	if filter == "user" {
		h.GetUser(c)
	} else if filter == "room" {
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

func (h *UserHandler) GetUser(c fiber.Ctx) error {
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
	getUserRequest.By = "user_id"
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

func (h *UserHandler) CreateUser(c fiber.Ctx) error {

	type CreateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	createUserRequest := new(CreateUserRequest)
	e := c.Bind().Body(createUserRequest)
	if e != nil {
		return jsonUtils.WriteJson(fiber.StatusBadRequest, nil, &jsonUtils.ApiError{
			Code:    1,
			Message: e.Error(),
			Details: e.Error(),
		}, c)
	}

	// Call gRPC method to create user
	var req = &proto.UpdateUserRequest{
		User: &proto.User{
			Username: createUserRequest.Username,
			Email:    createUserRequest.Email,
			Password: createUserRequest.Password,
		},
		IsCreate: true,
	}

	response, err := h.grpcUserClient.UpdateUser(c.Context(), req)

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

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
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

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
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
