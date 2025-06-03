package auth

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/vedantkulkarni/mqchat/api/middleware" // For token generation (to verify it's called)
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"google.golang.org/grpc" // Added for grpc.CallOption
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockUserGRPCServiceClient is a mock implementation of UserGRPCServiceClient for testing.
type MockUserGRPCServiceClient struct {
	proto.UnimplementedUserGRPCServiceServer                                                                 // Recommended for forward compatibility
	GetUserFunc    func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error)
	GetUsersFunc   func(ctx context.Context, req *proto.GetUsersRequest, opts ...grpc.CallOption) (*proto.GetUsersResponse, error)
	UpdateUserFunc func(ctx context.Context, req *proto.UpdateUserRequest, opts ...grpc.CallOption) (*proto.UpdateUserResponse, error)
	DeleteUserFunc func(ctx context.Context, req *proto.DeleteUserRequest, opts ...grpc.CallOption) (*proto.DeleteUserResponse, error)
	// Add other methods if AuthService starts using them
}

// GetUser implements the UserGRPCServiceClient interface for the mock.
func (m *MockUserGRPCServiceClient) GetUser(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, req, opts...)
	}
	return nil, errors.New("GetUserFunc not implemented in mock")
}

// Implement other methods of UserGRPCServiceClient to match the interface.
func (m *MockUserGRPCServiceClient) GetUsers(ctx context.Context, req *proto.GetUsersRequest, opts ...grpc.CallOption) (*proto.GetUsersResponse, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx, req, opts...)
	}
	return nil, errors.New("GetUsersFunc not implemented")
}

func (m *MockUserGRPCServiceClient) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest, opts ...grpc.CallOption) (*proto.UpdateUserResponse, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, req, opts...)
	}
	return nil, errors.New("UpdateUserFunc not implemented")
}

func (m *MockUserGRPCServiceClient) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest, opts ...grpc.CallOption) (*proto.DeleteUserResponse, error) {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, req, opts...)
	}
	return nil, errors.New("DeleteUserFunc not implemented")
}


func TestAuthService_LoginUser(t *testing.T) {
	// Mock the middleware.GenerateToken function if it has external dependencies or to simplify testing
	// For this example, we'll assume middleware.GenerateToken is simple enough.
	// If GenerateToken were complex, you might need to refactor it to be interface-based and mockable,
	// or use a package-level variable to swap its implementation during tests.

	// Pre-generate a token for comparison if needed, or just check that a token is returned.
	// Let's assume user ID 1 exists for successful login.
	expectedTokenForUser1, _ := middleware.GenerateToken(strconv.Itoa(1))


	testCases := []struct {
		name          string
		email         string
		password      string
		mockGetUser   func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error)
		expectedToken string
		expectedError string
	}{
		{
			name:     "Successful Login",
			email:    "test@example.com",
			password: "password123",
		mockGetUser: func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
				if req.Email == "test@example.com" && req.By == "email" {
					return &proto.GetUserResponse{User: &proto.User{Id: 1, Email: "test@example.com", Password: "password123"}}, nil
				}
				return nil, status.Error(codes.NotFound, "user not found")
			},
		expectedToken: expectedTokenForUser1, // This will be checked for non-emptiness
			expectedError: "",
		},
		{
			name:     "Incorrect Password",
			email:    "test@example.com",
			password: "wrongpassword",
		mockGetUser: func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
				if req.Email == "test@example.com" {
					return &proto.GetUserResponse{User: &proto.User{Id: 1, Email: "test@example.com", Password: "password123"}}, nil
				}
				return nil, status.Error(codes.NotFound, "user not found")
			},
		expectedToken: "", // Expect no token
			expectedError: "invalid password",
		},
		{
			name:     "User Not Found",
			email:    "nonexistent@example.com",
			password: "password123",
		mockGetUser: func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
				return nil, status.Error(codes.NotFound, "user not found")
			},
		expectedToken: "", // Expect no token
			expectedError: "user not found", // This error comes from our service logic based on gRPC status
		},
		{
			name:     "Password Too Short",
			email:    "test@example.com",
			password: "pass",
		mockGetUser: func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
				// This won't even be called if password validation fails first
				return &proto.GetUserResponse{User: &proto.User{Id: 1, Email: "test@example.com", Password: "password123"}}, nil
			},
		expectedToken: "", // Expect no token
			expectedError: "password must be at least 8 characters",
		},
		{
			name:     "GetUser returns non-gRPC error",
			email:    "test@example.com",
			password: "password123",
		mockGetUser: func(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
				return nil, errors.New("some internal user service error")
			},
		expectedToken: "", // Expect no token
			expectedError: "failed to get user: some internal user service error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserClient := &MockUserGRPCServiceClient{
				GetUserFunc: tc.mockGetUser,
			}
			authSvc := NewAuthService(mockUserClient)

			response, err := authSvc.LoginUser(context.Background(), tc.email, tc.password)

			if tc.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error '%s', but got nil", tc.expectedError)
				}
				if err.Error() != tc.expectedError {
					t.Fatalf("Expected error message '%s', but got '%s'", tc.expectedError, err.Error())
				}
				if response != nil {
					t.Errorf("Expected nil response on error, but got %+v", response)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, but got '%s'", err.Error())
				}
				if response == nil {
					t.Fatal("Expected a response, but got nil")
				}
				if response.Token == "" { // Check if token is non-empty instead of exact match if generation is tricky
					t.Error("Expected a token, but got an empty string")
				}
				// If you can reliably mock/predict token generation, uncomment:
				// if response.Token != tc.expectedToken {
				// 	t.Errorf("Expected token '%s', but got '%s'", tc.expectedToken, response.Token)
				// }
			}
		})
	}
}
