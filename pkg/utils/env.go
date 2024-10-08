package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvVarInt(key string, fallback uint16) string {

	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		log.Printf("Environment variable %s not set or is empty, returning fallback", key)
		return strconv.Itoa(int(fallback))
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Environment variable %s is invalid, returning fallback", key)
	}

	// Check for range of port numbers
	if res < 0 || res > 65535 {
		log.Printf("Environment variable %s is out of range, returning fallback", key)
	}

	return strconv.Itoa(res)

}

func GetEnvVar(key string, fallback string) string {

	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		log.Printf("Environment variable %s not set, returning fallback", key)
		return fallback
	}

	return value

}
