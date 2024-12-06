package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getImage(c *gin.Context) {
	url := "https://picsum.photos/500/500"
	// Download the image
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer response.Body.Close()

	// Set the content type to image/jpeg
	c.Header("Content-Type", "image/jpeg")

	// Copy the image directly to the response writer
	_, err = io.Copy(c.Writer, response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
		return
	}
}

func getImageResolution(c *gin.Context) {
	height := c.Param("height")
	width := c.Param("width")

	url := fmt.Sprintf("https://picsum.photos/%s/%s", width, height)
	// Download the image
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer response.Body.Close()

	// Set the content type to image/jpeg
	c.Header("Content-Type", "image/png")

	// Copy the image directly to the response writer
	_, err = io.Copy(c.Writer, response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
		return
	}
}

func main() {
	router := gin.Default()

	// Add the GET endpoint for the image
	router.GET("/getImage", getImage)
	router.GET("/getImage/:width/:height", getImageResolution)

	// Run the server
	router.Run("localhost:8080")
}
