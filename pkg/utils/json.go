package utils

import (
	"github.com/gofiber/fiber/v3"
)

type ApiResponse struct {
	Data   interface{} `json:"data"`
	Error  *ApiError `json:"error"`
	Status string      `json:"status"`
}

type ApiError struct {
	Code   int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

var BadRequestApiError *ApiError = &ApiError{
	Code: fiber.ErrBadRequest.Code,
	Message: "Bad Request",
	Details: "Bad Request",
}

var InternalServerApiError *ApiError = &ApiError{
	Code: fiber.ErrInternalServerError.Code,
	Message: "Internal Server Error",
	Details: "Internal Server Error",
}


func WriteJson(code int, data interface{}, err *ApiError, c fiber.Ctx) error {
	apiResponse := &ApiResponse{}
	c.Status(code)
	var status string
	if err != nil {
		status = "error"
	} else {
		status = "success"
	}

	apiResponse.Status = status
	apiResponse.Data = data
	apiResponse.Error = err

	return c.JSON(
		apiResponse,
	)

}


