package router

import (
	"fatawa-app-gcp/backend/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func getDocumentByID(c *gin.Context) {
	docID := c.Param("id")

	doc, err := db.GetDocumentByID(c, docID)
	if err != nil {
		if strings.Contains(err.Error(), "document with ID") && strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func createDocument(c *gin.Context) {
	var doc db.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Extract the document ID from the URL parameter
	docID := c.Param("id")
	if docID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	// Check if the document already exists
	if _, err := db.GetDocumentByID(c, docID); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Document already exists"})
		return
	}

	// Pass the document ID to the CreateDocument function
	if err := db.CreateDocument(c, docID, doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document created successfully", "document_id": docID})
}

func updateDocument(c *gin.Context) {
	var updatedDoc db.Document
	if err := c.ShouldBindJSON(&updatedDoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	docID := c.Param("id") // Assuming the document ID is passed as a URL parameter
	if err := db.UpdateDocument(c, docID, updatedDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func deleteDocument(c *gin.Context) {
	docID := c.Param("id") // Assuming the document ID is passed as a URL parameter

	if err := db.DeleteDocument(c, docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

func searchDocuments(c *gin.Context) {
	searchQuery := c.Query("query")

	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Search query is required"})
		return
	}

	docs, err := db.SearchDocuments(c, searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to search documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

func getAllDocuments(c *gin.Context) {
	docs, err := db.GetAllDocuments(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}
