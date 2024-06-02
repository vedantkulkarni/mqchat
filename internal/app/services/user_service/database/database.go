package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/vedantkulkarni/mqchat/internal"
)

type Database interface {
	
}

type PostgresDB struct {
	DB *sql.DB
}
func NewPostgresDB() (*PostgresDB, error) {

	pool, err := pgxpool.New(context.Background(), internal.GoDotEnvVariable("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection : %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	pool.Ping(context.Background())

	sqlDB := stdlib.OpenDBFromPool(pool)

	return &PostgresDB{
		DB: sqlDB,
	}, nil
}

func (p *PostgresDB) NewSqlBoilerDb( ) error {
	return nil
}
