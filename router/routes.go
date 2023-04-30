package router

import (
	"fatawa-app-gcp/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDocumentByID(c *gin.Context) {
	docID := c.Param("id")

	doc, err := db.GetDocumentByID(c, docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func createDocument(c *gin.Context) {
	var doc db.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Bad Request"})
		return
	}

	// Check if document already exists
	if _, err := db.GetDocumentByID(c, doc.ID); err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": "Document already exists"})
		return
	}

	docID, err := db.CreateDocument(c, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"document_id": docID})
}

func searchDocuments(c *gin.Context) {
	searchQuery := c.Query("search")

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
