package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(user fiber.Router) error {
	user.Get("/", getUsers)
	user.Post("/", createUser)
	user.Get("/:uid", getUser)
	user.Put("/:uid", updeUser)
	user.Delete("/:uid", deleteUser)

	return nil
}

func getUsers(c fiber.Ctx) error {
	return c.SendString("Get Users")

}

func createUser(c fiber.Ctx) error {
	return c.SendString("Create User")

}

func getUser(c fiber.Ctx) error {
	return c.SendString("Get User")

}

func updeUser(c fiber.Ctx) error {
	return c.SendString("Update User")

}

func deleteUser(c fiber.Ctx) error {
	return c.SendString("Delete User")

}
