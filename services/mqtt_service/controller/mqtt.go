package controller

import (
	"bytes"
	"context"
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

type ChatHookOptions struct {
	Server                    *mqtt.Server
	ClientConns               *map[string]*mqtt.Client
	ChatGRPCClient            *proto.ChatServiceClient
	ChatGRPCGetMessagesClient *proto.ChatService_GetMessagesClient
}

type ChatMQTTHook struct {
	mqtt.HookBase
	config *ChatHookOptions
}

func (h *ChatMQTTHook) ID() string {
	return "events-example"
}

func (h *ChatMQTTHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

func (h *ChatMQTTHook) Init(config any) error {
	h.Log.Info("initialised")
	if _, ok := config.(*ChatHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	h.config = config.(*ChatHookOptions)
	if h.config.Server == nil {
		return mqtt.ErrInvalidConfigType
	}
	return nil
}

// subscribeCallback handles messages for subscribed topics
func (h *ChatMQTTHook) subscribeCallback(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
}

func (h *ChatMQTTHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	//Add to map
	fmt.Printf("Client connected : with userID : %v ", pk.Connect.Username)
	message := string(pk.Connect.Username)
	(*h.config.ClientConns)[message] = cl
	// pk.
	//TODO: Basic checks if the the userID is authentic, if it exisits or if its already connected

	// Add to map
	(*h.config.ClientConns)[message] = cl

	fmt.Printf("Client connected %v \n", cl.ID)
	fmt.Printf("User ID %v \n", message)
	fmt.Println((*h.config.ClientConns))

	// Example demonstrating how to subscribe to a topic within the hook.
	h.config.Server.Subscribe(utils.ClientMessageTopic+cl.ID, 0, h.subscribeCallback)

	// Example demonstrating how to publish a message within the hook

	err := h.config.Server.Publish(utils.ClientChatTopic+cl.ID, []byte("Connect!"), false, 0)
	if err != nil {
		h.Log.Error("hook.publish", "error", err)
	}

	return nil
}

func (h *ChatMQTTHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	//Remove from map
	delete(*h.config.ClientConns, cl.ID)

	if err != nil {
		h.Log.Info("client disconnected", "client", cl.ID, "expire", expire, "error", err)
	} else {
		h.Log.Info("client disconnected", "client", cl.ID, "expire", expire)
	}

}

func (h *ChatMQTTHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
}

func (h *ChatMQTTHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
}

func (h *ChatMQTTHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	// h.Log.Info("received from client", "client", cl.ID, "payload", string(pk.Payload))
	message := pk.Payload
	fmt.Printf("Message received from client %v : %v \n", cl.ID, message)
	chatMessage := &proto.Message{}
	protojson.Unmarshal(message, chatMessage)

	// Send chat message to grpc server
	sendMessageRequest := &proto.SendMessageRequest{
		Message: chatMessage,
	}

	fmt.Println("Sending message to grpc server")
	response, err := (*h.config.ChatGRPCClient).SendMessage(context.Background(), sendMessageRequest)

	if err != nil || response == nil {
		h.Log.Error("Error occurred while sending message to grpc server", "error", err)
		// _ = h.config.Server.Publish(utils.ClientConnSub+cl.ID, []byte("Failed to send message, please try again!"), false, 0)
	}

	fmt.Printf("Response from chat server : %v", response)

	// Send message to the client
	fmt.Println("Sending message to client")
	//TODO: Define a commone json template for all pub sub messages for chat
	userId := fmt.Sprintf("%v", sendMessageRequest.Message.UserId_2)
	fmt.Printf("User ID : %v \n", userId)
	Client := (*h.config.ClientConns)[userId]
	if Client == nil {
		fmt.Println("Client not found")
		return pk, nil
	}

	// Client.WritePacket(packets.NewPackets(packets.FixedHeader{Type: packets.Publish, Qos: 0}, []byte(fmt.Sprintf("%v", response))))
	Client.WritePacket(packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Publish,
			Qos:  0,
		},
		Payload: []byte(fmt.Sprintf("%v", response)),
	})

	// _ = h.config.Server.Publish(utils.ClientConnSub+id, []byte(fmt.Sprintf("%v", response)), false, 0)

	pkx := pk
	// if string(pk.Payload) == "hello" {
	// 	pkx.Payload = []byte("hello world")
	// 	// h.Log.Info("received modified packet from client", "client", cl.ID, "payload", string(pkx.Payload))
	// }

	return pkx, nil
}

func (h *ChatMQTTHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	// h.Log.Info("published to client", "client", cl.ID, "payload", string(pk.Payload))
}
