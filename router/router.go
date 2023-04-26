package router

import (
	"context"
	"fatawa-app-gcp/auth"
	"fatawa-app-gcp/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getDocumentByID(c *gin.Context) {
	docID := c.Param("id")

	doc, err := db.GetDocumentByID(c, docID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, doc)
}

func createDocument(c *gin.Context) {
	var doc db.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	docID, err := db.CreateDocument(c, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"document_id": docID})
}

func SetupRouter(ctx context.Context) *gin.Engine {
	r := gin.Default()

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
		v1.Use(auth.AuthenticateUser())

		v1.GET("/documents/:id", getDocumentByID)
		v1.POST("/documents", createDocument)
		v1.POST("login", auth.Login)
	}

	return r
}
