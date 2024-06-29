package connection_service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vedantkulkarni/mqchat/database"
	"github.com/vedantkulkarni/mqchat/gen/models"
	"github.com/vedantkulkarni/mqchat/internal/app/protogen/proto"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type ConnectionGRPCServer struct {
	proto.UnimplementedConnectionGRPCServiceServer
	DB *sql.DB
}

func (c *ConnectionGRPCServer) CreateConnection(ctx context.Context, req *proto.CreateConnectionRequest) (*proto.CreateConnectionResponse, error) {
	connection := &models.Connection{
		UserID1: int(req.UserId_1),
		UserID2: int(req.UserId_2),
	}

	err := connection.Insert(ctx, c.DB, boil.Infer())
	if err != nil {
		// database error
		fmt.Println(err)
		return nil, err
	}

	response := &proto.CreateConnectionResponse{}
	response.ConnId = int64(connection.ID)

	return response, nil
}

func (c *ConnectionGRPCServer) GetConnection(ctx context.Context, req *proto.GetConnectionRequest) (*proto.GetConnectionResponse, error) {
	connection := &models.Connection{
		ID:      int(req.ConnId),
		UserID1: int(req.UserId_1),
		UserID2: int(req.UserId_2),
	}

	fmt.Println(connection)

	conn, err := models.Connections(qm.Where("user_id_1=?", connection.UserID1), qm.And("user_id_2=?", connection.UserID2)).One(ctx, c.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Error while getting connection: %v", err))
	}

	if conn == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Unknown error occurred while getting connection"))
	}

	fmt.Println(conn)
	response := &proto.GetConnectionResponse{
		Connection: &proto.Connection{
			UserId_1:   int64(conn.UserID1),
			UserId_2:   int64(conn.UserID2),
			UserName_1: "Vedant",
			UserName_2: "Vedant",
			Id:         int64(conn.ID),
		},
	}
	return response, nil
}

func (c *ConnectionGRPCServer) GetConnections(ctx context.Context, req *proto.GetConnectionsRequest) (*proto.GetConnectionsResponse, error) {
	UserID := int(req.UserId)

	conn, err := models.Connections(qm.Where("user_id_1=?", UserID), qm.Or("user_id_2=?", UserID)).All(ctx, c.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Error while getting connection: %v", err))
	}

	if conn == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Unknown error occurred while getting connection"))
	}

	fmt.Println(conn)

	var connections []*proto.Connection

	for _, conn := range conn {
		connections = append(connections, &proto.Connection{
			Id:         int64(conn.ID),
			UserId_1:   int64(conn.UserID1),
			UserId_2:   int64(conn.UserID2),
			UserName_1: "Vedant",
			UserName_2: "Vedant",
		})
	}

	response := &proto.GetConnectionsResponse{
		UserId:      int64(UserID),
		Connections: connections,
	}
	return response, nil
}

func NewConnectionGRPCServer(db *database.DbInterface) (*ConnectionGRPCServer, error) {
	return &ConnectionGRPCServer{
		DB: db.DB,
	}, nil
}

func (c *ConnectionGRPCServer) StartService(listener net.Listener) error {

	g := grpc.NewServer()
	fmt.Println("Starting gRPC connection server")

	proto.RegisterConnectionGRPCServiceServer(g, c)

	fmt.Println("gRPC connection server registered successfully")
	if err := g.Serve(listener); err != nil {
		fmt.Println("Error occurred while serving the gRPC server")
		return err
	}
	fmt.Println("gRPC connection server started successfully!")
	return nil

}
