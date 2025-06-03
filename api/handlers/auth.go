package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/gen/proto" // For AuthGRPCServiceClient and request/response types
	"github.com/vedantkulkarni/mqchat/pkg/utils" // For WriteJson and error utilities
	"google.golang.org/grpc/status"             // For handling gRPC errors
)

type AuthHandler struct {
	authService proto.AuthGRPCServiceClient // Changed from userService
	// userService *proto.UserGRPCServiceClient // Keep if other methods in AuthHandler might need it. For now, removing.
}

func NewAuthHandler(authClient proto.AuthGRPCServiceClient) *AuthHandler { // Changed parameter
	return &AuthHandler{
		authService: authClient,
	}
}

func (a *AuthHandler) RegisterAuthRoutes(auth fiber.Router) error {
	auth.Post("/login/", a.login)
	return nil
}

func (a *AuthHandler) login(c fiber.Ctx) error {
	requestPayload := new(proto.AuthLoginRequest) // Use the generated proto type for request
	err := c.Bind().Body(requestPayload)
	if err != nil {
		return utils.WriteJson(fiber.StatusBadRequest, nil, &utils.ApiError{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot parse request: " + err.Error(),
			Details: err.Error(),
		}, c)
	}

	// Basic validation before hitting the service (optional, as service also validates)
	if requestPayload.Email == "" || requestPayload.Password == "" {
		return utils.WriteJson(fiber.StatusBadRequest, nil, &utils.ApiError{
			Code:    fiber.StatusBadRequest,
			Message: "Email and password are required.",
		}, c)
	}

	// Call the Login RPC on the AuthGRPCServiceClient
	grpcResponse, err := a.authService.Login(c.Context(), requestPayload)
	if err != nil {
		// Handle gRPC level errors (e.g., service unavailable)
		st, _ := status.FromError(err)
		fiberErr := utils.CheckGRPCError(*st) // Assuming CheckGRPCError converts gRPC status to fiber error info
		return utils.WriteJson(fiberErr.Code, nil, &utils.ApiError{
			Code:    fiberErr.Code,
			Message: "Login failed (gRPC error): " + st.Message(),
			Details: st.Message(),
		}, c)
	}

	// Handle application-level errors returned in the response body
	if grpcResponse.Error != "" {
		// Determine appropriate HTTP status code based on the error message
		// This is a simplification; more structured errors would be better.
		httpStatusCode := fiber.StatusUnauthorized // Default to Unauthorized
		if grpcResponse.Error == "user not found" {
			httpStatusCode = fiber.StatusNotFound
		} else if grpcResponse.Error == "password must be at least 8 characters" {
			httpStatusCode = fiber.StatusBadRequest
		}

		return utils.WriteJson(httpStatusCode, nil, &utils.ApiError{
			Code:    httpStatusCode,
			Message: grpcResponse.Error,
		}, c)
	}

	// Success
	return utils.WriteJson(
		fiber.StatusOK,
		fiber.Map{
			"token": grpcResponse.Token,
		},
		nil,
		c,
	)
}
