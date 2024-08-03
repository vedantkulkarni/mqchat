package user

import (
	"context"
	"database/sql"
	"fmt"

	database "github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct { // Depends on DB
	proto.UnimplementedUserGRPCServiceServer
	RoomGRPCClient proto.RoomGRPCServiceClient
	DB             *sql.DB
}

func (u *UserGRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	var user *models.User
	var err error
	if req.By == "user_id" {
		user, err = models.Users(qm.Where("user_id=?", req.Id)).One(ctx, u.DB)
	} else if req.By == "email" {
		user, err = models.Users(qm.Where("email=?", req.Email)).One(ctx, u.DB)
	}

	if err != nil || user == nil {
		return nil, status.Error(codes.NotFound, "User does not exist.")
	}

	fmt.Printf("User found : %v", user)

	newUser := proto.GetUserResponse{
		User: &proto.User{
			Id:       int64(user.UserID),
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		}}

	return &newUser, nil

}

func (u *UserGRPCServer) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {

	uid := req.Id

	response, err := u.RoomGRPCClient.GetRooms(ctx, &proto.GetRoomsRequest{UserId: uid})

	if err != nil {
		return nil, status.Error(codes.NotFound, "Rooms not found")
	}

	var user_ids []int

	for _, conn := range response.Rooms {
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
			Id:       int64(user.UserID),
			Username: user.Username,
			Email:    user.Email,
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

	if req.IsCreate {
		return createUser(ctx, u, req)
	} else {
		return updateUser(ctx, u, req)
	}

}

// Helper functions

func createUser(ctx context.Context, u *UserGRPCServer, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	fmt.Println("Handled by CreateUser in UserService")
	user := &models.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: req.User.Password,
	}

	err := user.Insert(ctx, u.DB, boil.Infer())

	if err != nil {
		fmt.Printf("Error occurred while creating user : %v", err)
		pgErr := database.ParsePGXErrorUser(err)
		return nil, status.Error(codes.Internal, pgErr)
	}

	createUserResponse := &proto.User{
		Id:       int64((user.UserID)),
		Username: user.Username,
		Email:    user.Email,
	}

	fmt.Printf("User created successfully : %v", createUserResponse)

	return &proto.UpdateUserResponse{
		User: createUserResponse,
	}, nil
}

func updateUser(ctx context.Context, u *UserGRPCServer, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	fmt.Println("Handled by UpdateUser in UserService")
	user := &models.User{
		Username: req.User.Username,
		Email:    req.User.Email,
	}

	_, err := models.Users(qm.Where("user_id=?", req.User.Id)).UpdateAll(ctx, u.DB, models.M{
		"user_name":  user.Username,
		"user_email": user.Email,
	})

	if err != nil {
		fmt.Printf("Error occurred while updating user : %v", err)
		pgErr := database.ParsePGXErrorUser(err)
		return nil, status.Error(codes.Internal, pgErr)
	}

	updateUserResponse := &proto.User{
		Id:       req.User.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	fmt.Printf("User updated successfully : %v", updateUserResponse)

	return &proto.UpdateUserResponse{
		User: updateUserResponse,
	}, nil
}
