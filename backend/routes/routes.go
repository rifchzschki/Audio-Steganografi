package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/hello", controllers.HandleHello) 
	}

	return router
}