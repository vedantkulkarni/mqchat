package config

import (
	"sync"

	util "github.com/vedantkulkarni/mqchat/pkg/utils"
)

type Config struct {
	HttpPort        string
	UserServicePort string
	RoomServicePort string
	ChatServicePort string
	MQTTServicePort string
}

var config Config
var once sync.Once

func Get() *Config {

	// TODO: Implement for dev and prod environments

	once.Do(func() {
		config.HttpPort = util.GetEnvVarInt("HTTP_PORT", 8080)
		config.UserServicePort = util.GetEnvVarInt("USER_SERVICE_GRPC_PORT", 8003)
		config.RoomServicePort = util.GetEnvVarInt("ROOMS_SERVICE_GRPC_PORT", 8004)
		config.ChatServicePort = util.GetEnvVarInt("CHAT_SERVICE_GRPC_PORT", 8002)
		config.MQTTServicePort = util.GetEnvVarInt("MQTT_PORT", 8001)

	})

	return &config
}
