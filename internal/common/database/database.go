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

	var user string
	err = dbpool.QueryRow(context.Background(), "select last_name from users where first_name='Vedant'").Scan(&user)
	if err != nil {
		fmt.Println("Error while obtaining users")
		fmt.Print(err)
	}
	fmt.Println(user)
	return &PostgresDB{
		db: dbpool,
	}, nil
}
