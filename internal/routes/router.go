package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	//Generate Router Group
	api := router.Group("/api")

	//Register Routes
	RegisterUserRoutes(api)
	RegisterMessageRoutes(api)
	RegisterSessionRoutes(api)

}
