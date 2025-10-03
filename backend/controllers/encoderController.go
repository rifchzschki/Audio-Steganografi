package controllers

import (
	"fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
	"github.com/rifchzschki/Audio-Steganografi/backend/service/encoder"
)

func HandleEncode(c *gin.Context){
	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
    if err != nil {
        resp := models.NewStegoResponse(false, "Failed to parse form", 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

	audioFile, audioHeader, err := c.Request.FormFile("audioFile")
    if err != nil {
        resp := models.NewStegoResponse(false, "No audio file uploaded", 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }
    defer audioFile.Close()

	secretFile, secretHeader, err := c.Request.FormFile("secretFile")
    if err != nil {
        resp := models.NewStegoResponse(false, "No secret file uploaded", 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }
    defer secretFile.Close()

	key := c.PostForm("key")
    lsbBitsStr := c.PostForm("lsbBits")
    useEncryption := c.PostForm("useEncryption") == "true"
    useRandomStart := c.PostForm("useRandomStart") == "true"

	if key == "" {
        resp := models.NewStegoResponse(false, "Key is required", 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

    lsbBits, err := strconv.Atoi(lsbBitsStr)
    if err != nil || (lsbBits != 1 && lsbBits != 2 && lsbBits != 3 && lsbBits != 4) {
        resp := models.NewStegoResponse(false, "Invalid LSB bits (must be 1, 2, 3, or 4)", 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

	tempDir := "./tmp"
    os.MkdirAll(tempDir, 0755)

	audioPath := filepath.Join(tempDir, audioHeader.Filename)
    secretPath := filepath.Join(tempDir, secretHeader.Filename)

	if err := c.SaveUploadedFile(audioHeader, audioPath); err != nil {
        resp := models.NewStegoResponse(false, "Failed to save audio file", 0.0, "")
        c.JSON(http.StatusInternalServerError, resp)
        return
    }
    defer os.Remove(audioPath)

	if err := c.SaveUploadedFile(secretHeader, secretPath); err != nil {
        resp := models.NewStegoResponse(false, "Failed to save secret file", 0.0, "")
        c.JSON(http.StatusInternalServerError, resp)
        return
    }
    defer os.Remove(secretPath)

	outputDir := "./output"
    os.MkdirAll(outputDir, 0755)

	outputMP3 := filepath.Join(outputDir, "stego_"+audioHeader.Filename)

	outputName, psnrVal, _, err := encoder.EncodeFile(audioPath, secretPath, outputMP3, key, lsbBits, useEncryption, useRandomStart)
    if err != nil {
        resp := models.NewStegoResponse(false, err.Error(), 0.0, "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

	resp := models.NewStegoResponse(true, "Encode Success", psnrVal, filepath.Base(outputName))
    c.JSON(http.StatusOK, resp)
}

func HandleDownloadStego(c *gin.Context) {
    filename := c.Param("filename")
    if filename == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
        return
    }

    filePath := filepath.Join("./output", filename)
    
    // Check if file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    // Set headers for download
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    c.Header("Content-Type", "audio/mpeg")
    c.File(filePath)
}

func HandlePlayStego(c *gin.Context) {
    filename := c.Param("filename")
    if filename == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
        return
    }

    filePath := filepath.Join("./output", filename)
    
    // Check if file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    // Set headers for audio streaming
    c.Header("Content-Type", "audio/mpeg")
    c.Header("Accept-Ranges", "bytes")
    c.File(filePath)
}