package logging

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func ParsePGXError(err error) string {
    if pgErr, ok := err.(*pgconn.PgError); ok {
        switch pgErr.Code {
        case "23505": // unique_violation
            fmt.Println("Duplicate key value violates unique constraint:", pgErr.Message)
			return pgErr.Message
        case "23503": // foreign_key_violation
            fmt.Println("Foreign key violation:", pgErr.Message)
			return pgErr.Message
        default:
            fmt.Println("PostgreSQL error:", pgErr.Message)
			return pgErr.Message
        }
    } else {
        fmt.Println("Error is not a pgconn.PgError")
    }

	return "Error is not a pgconn.PgError"; 
}

