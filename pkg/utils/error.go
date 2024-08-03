package utils

import (
	"log"


	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func CheckGRPCError(err status.Status) *fiber.Error {


	log.Printf("gRPC error: %s", err.Message())

	if err.Code() == codes.InvalidArgument {
		return fiber.NewError(fiber.StatusBadRequest, err.Message())
	} else if err.Code() == codes.NotFound {
		return fiber.NewError(fiber.StatusNotFound, err.Message())
	} else if err.Code() == codes.PermissionDenied {
		return fiber.NewError(fiber.StatusUnauthorized, err.Message())
	} else if err.Code() == codes.Unauthenticated {
		return fiber.NewError(fiber.StatusUnauthorized, err.Message())
	} else if err.Code() == codes.AlreadyExists {
		return fiber.NewError(fiber.StatusConflict, err.Message())
	} else if err.Code() == codes.Unimplemented {
		return fiber.NewError(fiber.StatusNotImplemented, err.Message())
	} else {
		return fiber.NewError(fiber.StatusInternalServerError, err.Message())
	}
	
	
}