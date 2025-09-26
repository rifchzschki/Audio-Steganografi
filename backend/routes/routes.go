package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/controllers"
)

// SetupRouter mendefinisikan semua route di aplikasi.
func SetupRouter() *gin.Engine {
	// 1. Inisialisasi router Gin
	router := gin.Default()

	// 2. Definisi Endpoint (Routes)
	
	// Route Grup API
	api := router.Group("/api")
	{
		// Endpoint GET /api/hello akan ditangani oleh controllers.HandleHello
		api.GET("/hello", controllers.HandleHello) 
		
		// Endpoint GET /api/welcome
		api.GET("/welcome", controllers.HandleWelcome)
	}

	return router
}