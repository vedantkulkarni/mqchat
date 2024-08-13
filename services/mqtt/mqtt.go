package mqtt

import (
	"fmt"
	"sync"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/config"
	grpcUtils "github.com/vedantkulkarni/mqchat/pkg/grpc"
	"github.com/vedantkulkarni/mqchat/pkg/logger"
)

type MQTTService struct {
	Server                 *mqtt.Server
	ChatGRPCClient         *proto.ChatServiceClient
	ChatGRPCMessagesClient *proto.ChatService_GetMessagesClient
}

func NewMQTTService() *MQTTService {

	config := config.Get()

	mqttServer := mqtt.New(&mqtt.Options{
		InlineClient: false,
	})

	chatClient := grpcUtils.GetChatClientConn(config.ChatServicePort)

	_ = mqttServer.AddHook(new(auth.AllowHook), nil)
	_ = mqttServer.AddHook(new(ChatMQTTHook), &ChatHookOptions{Server: mqttServer, ChatGRPCClient: &chatClient})

	return &MQTTService{
		Server:         mqttServer,
		ChatGRPCClient: &chatClient}
}

func (m *MQTTService) Start(port string, wg *sync.WaitGroup) {
	l := logger.Get()

	defer wg.Done()

	listener := listeners.NewTCP(listeners.Config{Type: "tcp", Address: fmt.Sprintf(":%s", port)})

	err := m.Server.AddListener(listener)
	if err != nil {
		l.Panic().Err(err).Msg("Error adding listener to MQTT server")
	}

	err = m.Server.Serve()
	if err != nil {
		l.Panic().Err(err).Msg("Error starting MQTT server")
	}

}
