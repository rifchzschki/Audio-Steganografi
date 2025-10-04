package controllers

import (
	"fmt"
    "net/http"
    "os"
    "path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
	"github.com/rifchzschki/Audio-Steganografi/backend/service/decoder"
)

func HandleDecode(c *gin.Context){
	debug := false

	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
    if err != nil {
        resp := models.NewExtractResponse(false, "Failed to parse form", "", "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }
	stegoFile, stegoHeader, err := c.Request.FormFile("stegoFile")
    if err != nil {
        resp := models.NewExtractResponse(false, "No stego file uploaded", "", "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }
	defer stegoFile.Close()

	key := c.PostForm("key")
    useRandomStart := c.PostForm("useRandomStart") == "true"
    outputFileName := c.PostForm("outputFileName")

	if key == "" {
        resp := models.NewExtractResponse(false, "Key is required", "", "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

	tempDir := "./tmp"
    os.MkdirAll(tempDir, 0755)

	stegoPath := filepath.Join(tempDir, stegoHeader.Filename)
    if err := c.SaveUploadedFile(stegoHeader, stegoPath); err != nil {
        resp := models.NewExtractResponse(false, "Failed to save uploaded file", "", "")
        c.JSON(http.StatusInternalServerError, resp)
        return
    }
    defer os.Remove(stegoPath)

	outputDir := "./output"
    os.MkdirAll(outputDir, 0755)
    
	extractedFile, err := decoder.DecodeFile(stegoPath, key, outputFileName, useRandomStart, debug)
    if err != nil {
        resp := models.NewExtractResponse(false, err.Error(), "", "")
        c.JSON(http.StatusBadRequest, resp)
        return
    }

	finalPath := filepath.Join(outputDir, filepath.Base(extractedFile))
    if extractedFile != finalPath {
        if err := os.Rename(extractedFile, finalPath); err != nil {
            data, readErr := os.ReadFile(extractedFile)
            if readErr == nil {
                os.WriteFile(finalPath, data, 0644)
                os.Remove(extractedFile)
            }
        }
    }

    resp := models.NewExtractResponse(true, "Decode Success", extractedFile, filepath.Base(extractedFile))
    c.JSON(http.StatusOK, resp)

}

func HandleDownloadExtracted(c *gin.Context) {
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
    c.File(filePath)
}