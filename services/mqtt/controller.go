package mqtt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/logger"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	clientChatTopic string = "mqchat/client/chat/"
)

type ChatHookOptions struct {
	Server                 *mqtt.Server
	ChatGRPCClient         *proto.ChatServiceClient
	ChatGRPCMessagesClient *proto.ChatService_GetMessagesClient
}

type ChatMQTTHook struct {
	mqtt.HookBase
	config *ChatHookOptions
}

func (h *ChatMQTTHook) ID() string {
	return "chat-hook"
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
	//TODO: Basic checks if the the userID is authentic, if it exisits or if its already connected

	// Get messages
	getMessages(cl, pk, *h.config.ChatGRPCClient)

	return nil
}

func (h *ChatMQTTHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {

	//Clean up code
	h.config.Server.Unsubscribe(clientChatTopic+cl.ID, 0)
	h.config.Server.Clients.Delete(cl.ID)
}

func (h *ChatMQTTHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
}

func (h *ChatMQTTHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
}

func (h *ChatMQTTHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {

	return sendMessage(cl, pk, *h.config.ChatGRPCClient, h.config.Server)

	return pk, nil

}

func (h *ChatMQTTHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
}

// Helpers
func sendMessage(cl *mqtt.Client, pk packets.Packet, grpcClient proto.ChatServiceClient, server *mqtt.Server) (packets.Packet, error) {
	message := pk.Payload
	chatMessage := &proto.Message{}
	protojson.Unmarshal(message, chatMessage)

	fmt.Printf("Received message from client: %v\n", chatMessage)

	// Store chat message to a database
	sendMessageRequest := &proto.SendMessageRequest{
		Message: chatMessage,
	}
	response, err := grpcClient.SendMessage(context.Background(), sendMessageRequest)
	if err != nil || response == nil {
		utils.PublishError(cl, errors.New(fmt.Sprintf("an error occurred while sending the message: %v", err)))
		return pk, nil
	}

	fmt.Printf("Message sent successfully: %v\n", response)

	//TODO: Define a common json template for all pub sub messages for chat
	userId := fmt.Sprintf("%v", sendMessageRequest.Message.UserId_2)
	client, check := server.Clients.Get(userId)
	if !check || client == nil {

		//TODO : Maintain a queue service like RabbitMQ to store the queued messages instead of not allowing.
		utils.PublishError(cl, errors.New("oops! User is not connected at the moment"))
		return pk, nil
	}

	utils.PublishMessage(client, chatMessage)
	return pk, nil
}

func getMessages(cl *mqtt.Client, pk packets.Packet, grpcClient proto.ChatServiceClient) (packets.Packet, error) {
	l := logger.Get()

	message := pk.Payload
	getMessageRequest := &proto.GetMessagesRequest{}
	protojson.Unmarshal(message, getMessageRequest)

	stream, err := grpcClient.GetMessages(context.Background(), getMessageRequest)
	if err != nil {
		utils.PublishError(cl, errors.New("an error occurred while fetching the messages"))
	}

	messageCh := make(chan *proto.Message)
	errorCh := make(chan error)

	for {
		message, err := stream.Recv()
		if err != nil {
			l.Error().Err(err).Msg("error while receiving message")
			errorCh <- err
			break
		}

		messageCh <- message
	}

	var wg sync.WaitGroup

	select {
	case message := <-messageCh:
		wg.Add(1)
		go func() {
			utils.PublishMessage(cl, message)
			wg.Done()
		}()

	case err := <-errorCh:
		if err == io.EOF {
			l.Info().Msg("stream closed")
			stream.CloseSend()
		} else {
			l.Error().Err(err).Msg("error while receiving message")
		}
	}

	wg.Wait()

	return pk, nil
}
