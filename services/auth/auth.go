package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/vedantkulkarni/mqchat/api/middleware" // For GenerateToken
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginResponse contains the result of a login attempt.
// This will eventually be replaced by a gRPC message type.
type LoginResponse struct {
	Token string
	Error string // Consider a more structured error type for gRPC
}

// AuthService provides authentication operations.
type AuthService struct {
	UserClient proto.UserGRPCServiceClient
	// Add any other dependencies like JWT secret key if needed by GenerateToken
}

// NewAuthService creates a new AuthService.
func NewAuthService(userClient proto.UserGRPCServiceClient) *AuthService {
	return &AuthService{
		UserClient: userClient,
	}
}

// validatePassword checks if the password meets basic criteria.
func validatePassword(password string) bool {
	return len(password) >= 8
}

// LoginUser handles the core logic for authenticating a user.
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*LoginResponse, error) {
	if !validatePassword(password) {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if user exists
	userResponse, err := s.UserClient.GetUser(ctx, &proto.GetUserRequest{
		By:    "email",
		Email: email,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			// More specific error handling based on gRPC status code
			if st.Code() == codes.NotFound {
				return nil, errors.New("user not found")
			}
			return nil, fmt.Errorf("failed to get user: %w", err) // Wrap original error
		}
		return nil, fmt.Errorf("failed to get user: %w", err) // Non-gRPC error
	}

	if userResponse.User == nil {
		// This case should ideally be covered by codes.NotFound from GetUser
		return nil, errors.New("user not found (nil user object)")
	}

	// Compare passwords
	if userResponse.User.Password != password {
		return nil, errors.New("invalid password")
	}

	// Generate token
	token, err := middleware.GenerateToken(strconv.Itoa(int(userResponse.User.Id)))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{Token: token}, nil
}
