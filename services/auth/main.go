package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"

	"github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/config"
	"github.com/vedantkulkarni/mqchat/pkg/logger"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	authController "github.com/vedantkulkarni/mqchat/services/auth/controller"
)

func main() {
	// Initialize logger
	logManager := logger.NewLogger()
	defer func(Logfile *os.File) {
		err := Logfile.Close()
		if err != nil {
			log.Fatalf("Failed to close log file: %v", err)
		}
	}(logManager.Logfile)
	logManager.Init()

	// Load configuration
	cfg, err := config.LoadConfig(".") // Assuming config is in the root or accessible path
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbInterface, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbInterface.Close()

	// --- gRPC Client for User Service ---
	// The Auth service needs to call the User service.
	userSvcHost := utils.GetEnvVar("USER_SERVICE_GRPC_HOST", "localhost") // Default to localhost for local dev
	userSvcPort := utils.GetEnvVar("USER_SERVICE_GRPC_PORT", "8003")      // Default user service port
	userConnAddr := fmt.Sprintf("%s:%s", userSvcHost, userSvcPort)

	userConn, err := grpc.Dial(userConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User gRPC service at %s: %v", userConnAddr, err)
	}
	defer userConn.Close()
	userGRPCClient := proto.NewUserGRPCServiceClient(userConn)
	log.Printf("Successfully connected to User gRPC service at %s", userConnAddr)
	// --- End gRPC Client for User Service ---


	// Initialize AuthGRPCServer
	authServer, err := authController.NewAuthGRPCServer(dbInterface, userGRPCClient)
	if err != nil {
		log.Fatalf("Failed to create Auth gRPC server: %v", err)
	}

	// Start Auth gRPC service
	authServiceHost := utils.GetEnvVar("AUTH_SERVICE_GRPC_HOST", "0.0.0.0")
	authServicePort := utils.GetEnvVar("AUTH_SERVICE_GRPC_PORT", "8004") // Assign a new port for Auth service

	healthServer := health.NewServer() // Create a new health server

	log.Printf("Starting Auth gRPC service on %s:%s...", authServiceHost, authServicePort)
	if err := authServer.StartService(authServiceHost, authServicePort, healthServer); err != nil {
		log.Fatalf("Auth gRPC service failed to start: %v", err)
	}
}
