package router

import (
	"context"
	"fatawa-app-gcp/backend/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"net/url"
)

// proxyToProcessor forwards requests to the Python server
func proxyToProcessor(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the target URL
		targetURL, _ := url.Parse(target)

		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Update the request URL
		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.URL.Host = targetURL.Host
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))
		c.Request.Host = targetURL.Host // This line is to ensure the target host is set correctly

		// ServeHTTP uses the Go net/http package to forward the request and write the response.
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupRouter(ctx context.Context) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Request", "Authorization", "Origin", "Accept", "X-Requested-With", "Content-Type"},
		AllowCredentials: true,
	}))

	// Initialize database connection using singleton
	_, err := db.GetDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}

	// Define routes
	v1 := r.Group("/v1")
	{
		// Audio routes
		v1.POST("/audio", CreateAudio)
		v1.GET("/audio/:id", GetAudio)
		v1.DELETE("/audio/:id", DeleteAudio)
		
		// Segment routes
		v1.POST("/segment", CreateSegment)
		v1.GET("/audio/:id/segments", GetAudioSegments)
		v1.PUT("/segment/:id/processed", UpdateSegmentProcessed)
		
		// QA routes
		v1.POST("/qa", CreateQA)
		v1.GET("/segment/:id/qa", GetSegmentQA)
	}

	// Define proxy routes to forward requests to the Python FastAPI server
	processor := r.Group("/v1/processor")
	{
		processorURL := "http://localhost:8000"              // Adjust this to the URL of your Python server
		processor.Any("/*action", proxyToProcessor(processorURL)) // Proxy any request to /v1/processor/* to the Python server
	}

	return r
}
