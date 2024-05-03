package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vedantkulkarni/mqchat/internal"
)

type User struct {
}
type Database interface {
	//User querries
	CreateUser(*User) error
	DeleteUser(string) error
	GetUser(string) (*User, error)

	//Session querries

	//Message querries

}

type PostgresDB struct {
	db *pgxpool.Pool
}

func NewPostgresDB() (*PostgresDB, error) {

	dbpool, err := pgxpool.New(context.Background(), internal.GoDotEnvVariable("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	fmt.Println("Database Success !")
	return &PostgresDB{
		db: dbpool,
	}, nil
}
