package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
)

func HandleHello(c *gin.Context) {
	response := models.NewHelloResponse(true, "Hello, World!")
	c.JSON(http.StatusOK, response)
}
