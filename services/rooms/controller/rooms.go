package Rooms

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	database "github.com/vedantkulkarni/mqchat/db"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomsGRPCServer struct {
	proto.UnimplementedRoomGRPCServiceServer
	DB *sql.DB
}

func (c *RoomsGRPCServer) CreateRooms(ctx context.Context, req *proto.CreateRoomRequest) (*proto.CreateRoomResponse, error) {
	Room := &models.Room{
		UserID1: int(req.UserId_1),
		UserID2: int(req.UserId_2),
	}

	err := Room.Insert(ctx, c.DB, boil.Infer())
	if err != nil {
		// database error
		fmt.Println(err)
		return nil, err
	}

	response := &proto.CreateRoomResponse{}
	response.ConnId = int64(Room.RoomID)

	return response, nil
}



func (c *RoomsGRPCServer) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.GetRoomsResponse, error) {
	UserID := int(req.UserId)
	fmt.Printf("Handled request for user: %v by Roomss GRPC server", UserID)
	rooms, err := models.Rooms(qm.Where("user_id_1=?", UserID), qm.Or("user_id_2=?", UserID)).All(ctx, c.DB)
	if err != nil {

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Error while getting Rooms: %v", err))
	}

	if rooms == nil {
		return nil, status.Error(codes.NotFound, "Unknown error occurred while getting Rooms")
	}

	fmt.Println(rooms)

	var Rooms []*proto.Room

	for _, conn := range rooms {
		Rooms = append(Rooms, &proto.Room{
			Id:       int64(conn.RoomID),
			UserId_1: int64(conn.UserID1),
			UserId_2: int64(conn.UserID2),
		})
	}

	response := &proto.GetRoomsResponse{
		UserId:      int64(UserID),
		Rooms: Rooms,
	}
	return response, nil
}

func NewRoomsGRPCServer(db *database.DbInterface) (*RoomsGRPCServer, error) {
	return &RoomsGRPCServer{
		DB: db.DB,
	}, nil
}

func (c *RoomsGRPCServer) StartService() error {
	var block chan struct{}
	//Listen to gRPC responses

	port := utils.GetEnvVarInt("ROOMS_SERVICE_GRPC_PORT", 8004)
	host := utils.GetEnvVar("ROOMS_SERVICE_GRPC_HOST", "localhost")

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		fmt.Printf("Error occured while listening to the port %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			println("Error occurred while closing the listener")
		}
	}(listener)

	g := grpc.NewServer()
	fmt.Println("Starting gRPC Rooms server")

	proto.RegisterRoomGRPCServiceServer(g, c)

	if err := g.Serve(listener); err != nil {
		fmt.Println("Error occurred while serving the gRPC server")
		return err
	}

	<-block
	return nil

}
