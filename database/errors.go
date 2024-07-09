package database


import (
	"database/sql"
	"strings"
)


// ParsePGXError is a helper function to parse the error returned by pgx driver
func ParsePGXErrorUser(err error) string {
	if err == nil {
		return ""
	}
	if err == sql.ErrNoRows {
		return "Not found"
	}
	if strings.Contains(err.Error(), "violates unique constraint") {
		return "User already exists"
	}
	if strings.Contains(err.Error(), "violates not-null constraint") {
		return "Username or email cannot be empty"
	}
	if strings.Contains(err.Error(), "violates foreign key constraint") {
		return "Invalid user id"
	}
	if strings.Contains(err.Error(), "violates check constraint") {
		return "Invalid email"
	}
	if strings.Contains(err.Error(), "violates check constraint") {
		return "Invalid username"
	}
	return "Error occured while creating user"
}

