package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Create image source and fetcher
	imageSource := &PicsumImageSource{}
	imageFetcher := NewImageFetcher(imageSource)

	// Add endpoints with type parameter
	router.GET("/image", func(c *gin.Context) {
		imageType := ImageType(c.DefaultQuery("type", string(Random)))
		imageFetcher.GetImage(c, imageType)
	})

	router.GET("/image/:width/:height", func(c *gin.Context) {
		width := c.Param("width")
		height := c.Param("height")
		imageType := ImageType(c.DefaultQuery("type", string(Random)))
		imageFetcher.GetImageWithResolution(c, imageType, width, height)
	})

	// Run the server
	router.Run(":8080")
}
