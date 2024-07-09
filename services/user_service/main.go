package user_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/internal/app/protogen/proto"
	env "github.com/vedantkulkarni/mqchat/pkg/utils"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct { // Depends on DB and ConnectionGRPCServiceClient
	proto.UnimplementedUserGRPCServiceServer
	ConnGrpcClient proto.ConnectionGRPCServiceClient
	// MessageGrpcClient proto.MessageGRPCServiceClient
	DB *sql.DB
}

func (u *UserGRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	fmt.Println("Handled by CreateUser in UserService")
	user := &models.User{
		UserName:  req.Username,
		UserEmail: req.Email,
	}

	err := user.Insert(ctx, u.DB, boil.Infer())

	if err != nil {
		fmt.Printf("Error occurred while creating user : %v", err)
		pgErr := database.ParsePGXErrorUser(err)
		return nil, status.Error(codes.Internal, pgErr)
	}

	createUserResponse := &proto.User{
		Id:       int64((user.UserID)),
		Username: user.UserName,
		Email:    user.UserEmail,
	}

	fmt.Printf("User created successfully : %v", createUserResponse)

	return &proto.CreateUserResponse{
		User: createUserResponse,
	}, nil
}

func (u *UserGRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	fmt.Println("Handled by GetUser")

	user, err := models.Users(qm.Where("user_id=?", req.Id)).One(ctx, u.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	newUser := proto.GetUserResponse{
		User: &proto.User{
			Id:       int64(user.UserID),
			Username: user.UserName,
			Email:    user.UserEmail,
		}}

	return &newUser, nil

}

func (u *UserGRPCServer) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {

	uid := req.Id

	response, err := u.ConnGrpcClient.GetConnections(ctx, &proto.GetConnectionsRequest{UserId: uid})

	if err != nil {
		return nil, status.Error(codes.NotFound, "Connections not found")
	}

	var user_ids []int

	for _, conn := range response.Connections {
		if conn.UserId_1 == uid {
			user_ids = append(user_ids, int(conn.UserId_2))
		} else {
			user_ids = append(user_ids, int(conn.UserId_1))
		}
	}

	// Get all users from the database whose ids are in user_ids
	converted_ids := make([]interface{}, len(user_ids))
	for index, uid := range user_ids {
		converted_ids[index] = uid
	}

	users, err := models.Users(qm.Select("user_name", "user_email"), qm.WhereIn("user_id in ?", converted_ids...)).All(ctx, u.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Users not found")
	}

	var newUsers []*proto.User
	for _, user := range users {
		newUsers = append(newUsers, &proto.User{
			// Id:       int64(user.UserID),
			Username: user.UserName,
			Email:    user.UserEmail,
		})
	}

	return &proto.GetUsersResponse{
		Users: newUsers,
	}, nil

}

func (u *UserGRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {

	return nil, nil

}

func (u *UserGRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return nil, nil
}

func NewUserGRPCServer(db *database.DbInterface) (*UserGRPCServer, error) {
	connectionServicePort := env.GetEnvVarInt("CONNECTION_SERVICE_GRPC_PORT", 2100)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", connectionServicePort), opts...)
	if err != nil {
		log.Println("Error occurred while connecting to the gRPC server")
		return nil, err
	}

	return &UserGRPCServer{
		DB:             db.DB,
		ConnGrpcClient: proto.NewConnectionGRPCServiceClient(conn),
	}, nil
}

func (u *UserGRPCServer) StartService(listner net.Listener) error {

	g := grpc.NewServer()
	fmt.Println("Starting gRPC user server")

	proto.RegisterUserGRPCServiceServer(g, u)

	fmt.Println("gRPC user server registered successfully")
	if err := g.Serve(listner); err != nil {
		fmt.Println("Error occured while serving the gRPC server")
		return err
	}
	fmt.Println("gRPC user server started successfully!")
	return nil

}
