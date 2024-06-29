package user_service

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/internal/app/protogen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	proto.UnimplementedUserGRPCServiceServer
	DB *sql.DB
}

func (u *UserGRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user := &models.User{
		UserName:  req.Username,
		UserEmail: req.Email,
	}

	err := user.Insert(ctx, u.DB, boil.Infer())
	if err != nil {
		// pgErr:= database.ParsePGXError(err)
		// return nil, status.Error(codes.Internal, pgErr)
	}

	createUserResponse := &proto.User{
		Id:       strconv.Itoa(user.UserID),
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
	userId, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user id")
	}

	user, err := models.Users(qm.Where("user_id=?", userId)).One(ctx, u.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	newUser := proto.GetUserResponse{
		User: &proto.User{
			Id:       strconv.Itoa(user.UserID),
			Username: user.UserName,
			Email:    user.UserEmail,
		}}

	return &newUser, nil

}

func (u *UserGRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {

	return nil, nil

}

func (u *UserGRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return nil, nil
}

func NewUserGRPCServer(db *database.DbInterface) (*UserGRPCServer, error) {
	return &UserGRPCServer{
		DB: db.DB,
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
