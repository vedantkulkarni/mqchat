package database

import util "github.com/vedantkulkarni/mqchat/pkg/utils"
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewDatabaseConfig() *DatabaseConfig {
	// host := util.GetEnvVar("HOST", "localhost")
	port := util.GetEnvVar("DATABASE_PORT", "5432")
	user := util.GetEnvVar("DATABASE_USER", "postgres")
	password := util.GetEnvVar("DATABASE_PASSWORD", "postgres")

	return &DatabaseConfig{
		Host:     "mqchat-db-1",
		Port:     port,
		User:     user,
		Password: password,
	}
}
