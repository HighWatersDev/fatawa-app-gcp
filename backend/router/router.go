package router

import (
	"context"
	"fatawa-app-gcp/backend/auth"
	"fatawa-app-gcp/backend/db"
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
		v1.GET("/documents/:id", auth.AuthenticateUser(), GetDocumentByID)
		v1.POST("/documents", auth.AuthenticateUser(), CreateDocument)
		v1.PUT("/documents", auth.AuthenticateUser(), UpdateDocument)
		v1.GET("/documents/search", auth.AuthenticateUser(), SearchDocuments)
		v1.GET("/documents/all", auth.AuthenticateUser(), GetAllDocuments)
		v1.POST("/verify", auth.AuthenticateUser())
		v1.DELETE("/documents/:id", auth.AuthenticateUser(), DeleteDocument)
	}

	return r
}
