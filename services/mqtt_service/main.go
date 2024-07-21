package mqttservice

import (
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/services/mqtt_service/controller"
)

type MQTTService struct {
	Server                    *mqtt.Server
	ChatGRPCClient            *proto.ChatServiceClient
	ChatGRPCGetMessagesClient *proto.ChatService_GetMessagesClient
}

func NewMQTTService(chatClient *proto.ChatServiceClient, messageStreamClient *proto.ChatService_GetMessagesClient) *MQTTService {
	mqttServer := mqtt.New(&mqtt.Options{
		InlineClient: false,
	})
	var clientList map[string]*mqtt.Client = make(map[string]*mqtt.Client)
	_ = mqttServer.AddHook(new(auth.AllowHook), nil)
	_ = mqttServer.AddHook(new(controller.ChatMQTTHook), &controller.ChatHookOptions{Server: mqttServer, ChatGRPCClient: chatClient, ChatGRPCGetMessagesClient: messageStreamClient, ClientConns: &clientList})
	return &MQTTService{
		Server:                    mqttServer,
		ChatGRPCClient:            chatClient,
		ChatGRPCGetMessagesClient: messageStreamClient,
	}
}


func (m *MQTTService) Start(port string) {
	// Start the MQTT server
	listner:= listeners.NewTCP(listeners.Config{Type: "tcp", Address: fmt.Sprintf(":%s", port)})
	err:=m.Server.AddListener(listner)
	if err!=nil{
		fmt.Println("Error adding listener to mqtt server")
	}	

	
	go func() {
		err := m.Server.Serve()
		if err != nil {
			fmt.Println("Error occurred while serving mqtt server")
		}
	}()

}
