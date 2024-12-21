package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestPicsumImageSource_GetURL(t *testing.T) {
	source := &PicsumImageSource{}

	tests := []struct {
		imageType ImageType
		width     string
		height    string
		expected  string
	}{
		{Nature, "300", "400", "https://picsum.photos/300/400?nature"},
		{Random, "300", "400", "https://picsum.photos/300/400"},
		{Nature, "", "", "https://picsum.photos/500/500?nature"},
		{Random, "", "", "https://picsum.photos/500/500"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			url := source.GetURL(tt.imageType, tt.width, tt.height)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func TestImageFetcher_GetImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	imageSource := &PicsumImageSource{}
	imageFetcher := NewImageFetcher(imageSource)

	router.GET("/image", func(c *gin.Context) {
		imageType := ImageType(c.DefaultQuery("type", string(Random)))
		imageFetcher.GetImage(c, imageType)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/image?type=nature", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
}

func TestImageFetcher_GetImageWithResolution(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	imageSource := &PicsumImageSource{}
	imageFetcher := NewImageFetcher(imageSource)

	router.GET("/image/:width/:height", func(c *gin.Context) {
		width := c.Param("width")
		height := c.Param("height")
		imageType := ImageType(c.DefaultQuery("type", string(Random)))
		imageFetcher.GetImageWithResolution(c, imageType, width, height)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/image/300/400?type=nature", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
}

func TestImageFetcher_Cache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	imageSource := &PicsumImageSource{}
	imageFetcher := NewImageFetcher(imageSource)

	// Mock cache
	imageFetcher.cache = cache.New(5*time.Minute, 10*time.Minute)
	imageFetcher.cache.Set("https://picsum.photos/300/400?nature", []byte("cached image data"), cache.DefaultExpiration)

	router.GET("/image/:width/:height", func(c *gin.Context) {
		width := c.Param("width")
		height := c.Param("height")
		imageType := ImageType(c.DefaultQuery("type", string(Random)))
		imageFetcher.GetImageWithResolution(c, imageType, width, height)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/image/300/400?type=nature", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
	assert.Equal(t, "cached image data", w.Body.String())
}
