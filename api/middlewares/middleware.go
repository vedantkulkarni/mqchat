package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)


func AuthMiddleware(c fiber.Ctx) error {

	fmt.Printf("AuthMiddleware\n")
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	t, err:= ValidateToken(token)
	if err != nil || !t {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()

}

