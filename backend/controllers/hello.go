package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
)

func HandleHello(c *gin.Context) {
	response := models.MessageResponse{
		Status:  http.StatusOK,
		Message: "Hello Cuy from Go Gin Controller!",
	}

	c.JSON(http.StatusOK, response)
}
