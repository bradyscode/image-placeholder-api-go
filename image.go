package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// ImageSource defines a strategy for generating image URLs
type ImageSource interface {
	GetURL(imageType ImageType, width, height string) string
}

// PicsumImageSource implements ImageSource
type PicsumImageSource struct{}

func (p *PicsumImageSource) GetURL(imageType ImageType, width, height string) string {
	switch imageType {
	case Nature:
		if width != "" && height != "" {
			return fmt.Sprintf("https://picsum.photos/%s/%s?nature", width, height)
		}
		return "https://picsum.photos/500/500?nature"
	default:
		if width != "" && height != "" {
			return fmt.Sprintf("https://picsum.photos/%s/%s", width, height)
		}
		return "https://picsum.photos/500/500"
	}
}

// ImageFetcher now uses the ImageSource interface
type ImageFetcher struct {
	source ImageSource
	cache  *cache.Cache
}

func NewImageFetcher(source ImageSource) *ImageFetcher {
	c := cache.New(5*time.Minute, 10*time.Minute) // 5 minutes expiry, 10 minutes cleanup interval
	return &ImageFetcher{source: source, cache: c}
}

func (f *ImageFetcher) GetImage(c *gin.Context, imageType ImageType) error {
	url := f.source.GetURL(imageType, "", "")
	return f.fetchAndSendImage(c, url)
}

func (f *ImageFetcher) GetImageWithResolution(c *gin.Context, imageType ImageType, width, height string) error {
	url := f.source.GetURL(imageType, width, height)

	// Check if the image is already cached
	cachedImage, found := f.cache.Get(url)
	if found {
		// Serve the cached image directly
		c.Writer.Header().Set("Content-Type", "image/jpeg")
		_, err := c.Writer.Write(cachedImage.([]byte))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send cached image"})
			return err
		}
		return nil
	}

	// If not found in cache, fetch the image
	imageData, err := f.fetchImage(url)
	if err != nil {
		return err
	}

	// Store the fetched image data in the cache (assuming it's []byte)
	f.cache.Set(url, imageData, cache.DefaultExpiration)

	// Send the image response
	c.Writer.Header().Set("Content-Type", "image/jpeg")
	_, err = c.Writer.Write(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
		return err
	}

	return nil
}

func (f *ImageFetcher) fetchAndSendImage(c *gin.Context, url string) error {
	// First, try to fetch the image from the cache
	cachedImage, found := f.cache.Get(url)
	if found {
		// Serve the cached image directly
		c.Writer.Header().Set("Content-Type", "image/jpeg")
		_, err := c.Writer.Write(cachedImage.([]byte))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send cached image"})
			return err
		}
		return nil
	}

	// If not found in cache, fetch the image
	imageData, err := f.fetchImage(url)
	if err != nil {
		return err
	}

	// Store the fetched image data in the cache (assuming it's []byte)
	f.cache.Set(url, imageData, cache.DefaultExpiration)

	// Send the image response
	c.Writer.Header().Set("Content-Type", "image/jpeg")
	_, err = c.Writer.Write(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
		return err
	}

	return nil
}

func (f *ImageFetcher) fetchImage(url string) ([]byte, error) {
	// Fetch image from the provided URL
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
