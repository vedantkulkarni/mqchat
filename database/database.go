package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log"
)

type Database interface {
}

type DbInterface struct {
	DB *sql.DB
}

func NewPostgresDB() (*DbInterface, error) {
	config := NewDatabaseConfig()
	var dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable", config.User, config.Password, config.Host, config.Port)

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to create connection to database: %v\n", err)
	}
	// defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	return &DbInterface{
		DB: sqlDB,
	}, nil
}

func (p *DbInterface) NewSqlBoilerDb() error {
	return nil
}
