package router

import (
	"context"
	"fatawa-app-gcp/auth"
	"fatawa-app-gcp/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		v1.GET("/documents/:id", auth.AuthenticateUser(), getDocumentByID)
		v1.POST("/documents", auth.AuthenticateUser(), createDocument)
		v1.GET("/documents/search", auth.AuthenticateUser(), searchDocuments)
		v1.POST("/verify", auth.AuthenticateUser())
	}

	return r
}
