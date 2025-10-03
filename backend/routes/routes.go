package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://fe-audio-steg.vercell.app"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour, 
	}

	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		api.GET("/hello", controllers.HandleHello) 

		api.POST("/encode", controllers.HandleEncode)
        api.GET("/download/stego/:filename", controllers.HandleDownloadStego)
        api.GET("/play/stego/:filename", controllers.HandlePlayStego)
		
		api.POST("/decode", controllers.HandleDecode)
        api.GET("/download/extracted/:filename", controllers.HandleDownloadExtracted)
	}

	return router
}