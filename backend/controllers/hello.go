package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
)

// HandleHello adalah handler untuk endpoint /hello.
func HandleHello(c *gin.Context) {
	// Membuat objek response menggunakan model yang telah didefinisikan
	response := models.MessageResponse{
		Status:  http.StatusOK,
		Message: "Hello Proper World from Go Gin Controller!",
	}

	// Mengembalikan response JSON
	c.JSON(http.StatusOK, response)
}

// HandleWelcome adalah handler untuk endpoint /welcome, menunjukkan endpoint lain.
func HandleWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"greeting": "Welcome to the structured Go Gin API!",
	})
}