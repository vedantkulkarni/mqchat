package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	middleware "github.com/vedantkulkarni/mqchat/api/middlewares"
	"github.com/vedantkulkarni/mqchat/gen/proto"
	"github.com/vedantkulkarni/mqchat/pkg/utils"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	userService *proto.UserGRPCServiceClient
}

func NewAuthHandler(userService *proto.UserGRPCServiceClient) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}

}

func (a *AuthHandler) RegisterAuthRoutes(auth fiber.Router) error {

	auth.Post("/login/", a.login)

	return nil
}

func (a *AuthHandler) login(c fiber.Ctx) error {
	request := &proto.User{}
	err := c.Bind().Body(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !validatePassword(request.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters",
		})
	}

	//Check if user exists
	user, err := (*a.userService).GetUser(c.Context(), &proto.GetUserRequest{
		By:    "email",
		Email: request.Email,
	})

	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		fiberStatus := utils.CheckGRPCError(*errStatus)

		return utils.WriteJson(fiberStatus.Code, nil, &utils.ApiError{
			Code:    fiberStatus.Code,
			Message: fiberStatus.Message,
			Details: fiberStatus.Message,
		}, c)

	}

	//print both passwords
	fmt.Printf("User Password: %v\nRequest Password: %v\n", user.User.Password, request.Password)

	if user.User.Password != request.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Password",
		})
	}

	token, err := middleware.GenerateToken(strconv.Itoa(int(user.User.Id)))

	if err != nil {
		return utils.WriteJson(fiber.StatusInternalServerError, nil, utils.InternalServerApiError, c)

	}

	utils.WriteJson(
		fiber.StatusOK,
		fiber.Map{
			"token": token,
		},	
		nil,
		c,
	)

	return nil
}

//Helpers

// Validate password
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}
