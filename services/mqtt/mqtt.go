package mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc"
)

type MQTTService struct {
	Server         *mqtt.Server
	ChatGRPCClient *proto.ChatServiceClient
}

func NewMQTTService() *MQTTService {
	chatPort := utils.GetEnvVarInt("CHAT_SERVICE_GRPC_PORT", 2200)
	// Create a new MQTT server
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	chat, err := grpc.NewClient(fmt.Sprintf("localhost:%s", chatPort), opts...)
	if err != nil {
		log.Println("Error while listening to mqtt service")
	}
	chatClient := proto.NewChatServiceClient(chat)

	mqttServer := mqtt.New(&mqtt.Options{
		InlineClient: false,
	})
	_ = mqttServer.AddHook(new(auth.AllowHook), nil)
	_ = mqttServer.AddHook(new(controller.ChatMQTTHook), &controller.ChatHookOptions{Server: mqttServer, ChatGRPCClient: &chatClient})
	return &MQTTService{
		Server:         mqttServer,
		ChatGRPCClient: &chatClient,
	}
}

func (m *MQTTService) Start(port string) {
	var block chan struct{}
	// Start the MQTT server
	listener := listeners.NewTCP(listeners.Config{Type: "tcp", Address: fmt.Sprintf(":%s", port)})

	// defer func(listner *listeners.TCP) {
	// 	err :=
	// 	if err != nil {
	// 		fmt.Println("Error occurred while closing the listener")
	// 	}
	// }(listener)

	err := m.Server.AddListener(listener)
	if err != nil {
		fmt.Println("Error adding listener to mqtt server")
	}

	go func() {
		err := m.Server.Serve()
		if err != nil {
			fmt.Println("Error occurred while serving mqtt server")
		}
	}()

	<-block

}
