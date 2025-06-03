package controller

import (
	"context"
	"database/sql" // Required if NewAuthService needs it indirectly or for other setup
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"

	database "github.com/vedantkulkarni/mqchat/db" // For DbInterface
	"github.com/vedantkulkarni/mqchat/gen/proto"
	authService "github.com/vedantkulkarni/mqchat/services/auth" // Alias to avoid conflict
	"github.com/vedantkulkarni/mqchat/pkg/utils" // For env var functions
)

// AuthGRPCServer implements the gRPC server for authentication.
type AuthGRPCServer struct {
	proto.UnimplementedAuthGRPCServiceServer
	service *authService.AuthService
	DB      *sql.DB // Or database.DbInterface if NewAuthService needs the full interface
}

// NewAuthGRPCServer creates a new AuthGRPCServer.
// It initializes the underlying AuthService.
func NewAuthGRPCServer(db *database.DbInterface, userClient proto.UserGRPCServiceClient) (*AuthGRPCServer, error) {
	service := authService.NewAuthService(userClient)
	return &AuthGRPCServer{
		service: service,
		DB:      db.DB, // Assuming AuthService might need db access indirectly or for future features
	}, nil
}

// Login is the gRPC method that calls the AuthService's LoginUser method.
func (s *AuthGRPCServer) Login(ctx context.Context, req *proto.AuthLoginRequest) (*proto.AuthLoginResponse, error) {
	if req == nil {
		return &proto.AuthLoginResponse{Error: "request is nil"}, fmt.Errorf("request is nil")
	}

	response, err := s.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		// Convert domain error to gRPC response
		return &proto.AuthLoginResponse{Error: err.Error()}, nil // Sending error in response body
	}

	return &proto.AuthLoginResponse{Token: response.Token}, nil
}

// StartService initializes and starts the gRPC server for authentication.
func (s *AuthGRPCServer) StartService(host string, port string, healthServer *health.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to listen on %s:%s: %v", host, port, err)
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthGRPCServiceServer(grpcServer, s)

	// Register health check service
	if healthServer != nil {
		healthPb.RegisterHealthServer(grpcServer, healthServer)
		healthServer.SetServingStatus("AuthGRPCService", healthPb.HealthCheckResponse_SERVING)
		log.Println("Auth gRPC health server registered successfully.")
	}

	log.Printf("Auth gRPC server listening on %s:%s", host, port)
	if err := grpcServer.Serve(lis); err != nil {
		if healthServer != nil {
			healthServer.SetServingStatus("AuthGRPCService", healthPb.HealthCheckResponse_NOT_SERVING)
		}
		log.Fatalf("Failed to serve Auth gRPC server: %v", err)
		return err
	}
	return nil
}
