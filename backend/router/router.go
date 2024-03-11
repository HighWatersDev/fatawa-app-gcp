package router

import (
	"context"
	"fatawa-app-gcp/backend/auth"
	"fatawa-app-gcp/backend/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Initialize Firebase Auth client
	err := auth.InitializeAuthClient(ctx)
	if err != nil {
		panic(err)
	}

	// Initialize Firestore client
	err = db.InitializeFirestoreClient(ctx, "salafifatawa")
	if err != nil {
		panic(err)
	}

	// Define routes
	v1 := r.Group("/v1")
	{
		v1.GET("/documents/:id", auth.AuthenticateUser(), GetDocumentByID)
		v1.POST("/documents", auth.AuthenticateUser(), CreateDocument)
		v1.PUT("/documents", auth.AuthenticateUser(), UpdateDocument)
		v1.GET("/documents/search", auth.AuthenticateUser(), SearchDocuments)
		v1.GET("/documents/all", auth.AuthenticateUser(), GetAllDocuments)
		v1.POST("/verify", auth.AuthenticateUser())
		v1.DELETE("/documents/:id", auth.AuthenticateUser(), DeleteDocument)
	}

	// Define proxy routes to forward requests to the Python FastAPI server
	processor := r.Group("/v1/processor")
	{
		processorURL := "http://localhost:8000"                   // Adjust this to the URL of your Python server
		processor.Use(auth.AuthenticateUser())                    // Apply the same authentication middleware
		processor.Any("/*action", proxyToProcessor(processorURL)) // Proxy any request to /v1/processor/* to the Python server
	}

	return r
}
