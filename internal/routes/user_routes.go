package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {

	//Set function handlers for user routes
	rg.GET("/users/:id", getUsersByIdHandler)
	rg.POST("/users", createUserHandler)
}

func getUsersByIdHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"userId": "Hello World",
	})
}

func createUserHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"status" : "User Created Successfully !",
	})
}