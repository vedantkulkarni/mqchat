package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct { // Depends on DB and ConnectionGRPCServiceClient
	proto.UnimplementedUserGRPCServiceServer
	ConnGrpcClient proto.ConnectionGRPCServiceClient
	// MessageGrpcClient proto.MessageGRPCServiceClient
	DB *sql.DB
}


func (u *UserGRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	fmt.Println("Handled by GetUser")

	var user *models.User
	var err error
	if req.By == "user_id" {

		user, err = models.Users(qm.Where("user_id=?", req.Id)).One(ctx, u.DB)
	} else if req.By == "user_email" {

		user, err = models.Users(qm.Where("user_email=?", req.Email)).One(ctx, u.DB)
	}

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

	if(req.IsCreate) {
		return createUser(ctx, u, req)
	} else {
		return updateUser(ctx, u, req)
	}

	return nil, nil
}


// Helper functions

func createUser(ctx context.Context, u *UserGRPCServer, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	fmt.Println("Handled by CreateUser in UserService")
	user := &models.User{
		UserName:  req.User.Username,
		UserEmail: req.User.Email,
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

	return &proto.UpdateUserResponse{
		User: createUserResponse,
	}, nil
}


func updateUser(ctx context.Context, u *UserGRPCServer, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	fmt.Println("Handled by UpdateUser in UserService")
	user := &models.User{
		UserName:  req.User.Username,
		UserEmail: req.User.Email,
	}

	_, err := models.Users(qm.Where("user_id=?", req.User.Id)).UpdateAll(ctx, u.DB, models.M{
		"user_name":  user.UserName,
		"user_email": user.UserEmail,
	})

	if err != nil {
		fmt.Printf("Error occurred while updating user : %v", err)
		pgErr := database.ParsePGXErrorUser(err)
		return nil, status.Error(codes.Internal, pgErr)
	}

	updateUserResponse := &proto.User{
		Id:       req.User.Id,
		Username: user.UserName,
		Email:    user.UserEmail,
	}

	fmt.Printf("User updated successfully : %v", updateUserResponse)

	return &proto.UpdateUserResponse{
		User: updateUserResponse,
	}, nil
}