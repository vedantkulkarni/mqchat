package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	jsonUtils "github.com/vedantkulkarni/mqchat/pkg/utils"
)

type ChatHandler struct {
	//mqtt instance
	grpcChatClient proto.ChatServiceClient
}

func NewChatHandler(client proto.ChatServiceClient) *ChatHandler {
	return &ChatHandler{
		grpcChatClient: client,
	}
}

func (h *ChatHandler) RegisterChatRoutes(chat fiber.Router) error {
	chat.Get("/", h.GetMessages)
	chat.Post("/", h.SendMessage)

	return nil
}

func (h *ChatHandler) GetMessages(c fiber.Ctx) error {
	// Get the messages -> latest 20 msgs.
	initializeParam := c.Query("initialize")
	user_id_1 := c.Query("user_id_1")
	user_id_2 := c.Query("user_id_2")
	initialize := false

	if initializeParam == "1" {
		initialize = true
	} else {
		initialize = false
	}

	getMessages := new(proto.GetMessagesRequest)
	user_1, _ := strconv.Atoi(user_id_1)
	user_2, _ := strconv.Atoi(user_id_2)
	getMessages.Initialize = initialize
	getMessages.UserId_1 = int64(user_1)
	getMessages.UserId_2 = int64(user_2)

	response, err := h.grpcChatClient.GetMessages(c.Context(), getMessages)
	if err != nil {
		jsonUtils.WriteJson(
			fiber.ErrBadGateway.Code,
			nil,
			jsonUtils.BadRequestApiError,
			c,
		)
	}

	fmt.Printf("Sent message : %s", response)

	return jsonUtils.WriteJson(
		200,
		response,
		nil,
		c,
	)

}

func (h *ChatHandler) SendMessage(c fiber.Ctx) error {
	// Send message
	message := new(proto.Message)

	err := c.Bind().Body(&message)
	if err != nil {
		jsonUtils.WriteJson(
			fiber.ErrBadRequest.Code,
			nil,
			jsonUtils.BadRequestApiError,
			c,
		)
	}

	response, err := h.grpcChatClient.SendMessage(c.Context(), &proto.SendMessageRequest{
		Message: message,
	})
	if err != nil {
		jsonUtils.WriteJson(
			fiber.ErrBadGateway.Code,
			nil,
			jsonUtils.BadRequestApiError,
			c,
		)
	}

	fmt.Printf("Sent message : %s", response)

	return jsonUtils.WriteJson(
		200,
		response,
		nil,
		c,
	)

}
