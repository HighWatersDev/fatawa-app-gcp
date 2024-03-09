package router

import (
	"crypto/rand"
	"encoding/base64"
	"fatawa-app-gcp/backend/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type DocumentRequest struct {
	IsSplitAudio   bool              `json:"isSplitAudio"`
	ParentDocID    string            `json:"parentDocID,omitempty"` // omitempty allows this field to be omitted for non-split audio
	DocID          string            `json:"docID,omitempty"`
	Document       db.Document       `json:"document,omitempty"`
	ParentDocument db.ParentDocument `json:"parentDocument,omitempty"`
}

type SearchRequest struct {
	Query        string `json:"query"`
	IsSplitAudio bool   `json:"isSplitAudio"`
	ParentDocID  string `json:"parentDocID,omitempty"`
}

// GenerateRandomSuffix generates a random string of a specified length.
func GenerateRandomSuffix(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func CreateDocument(c *gin.Context) {
	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var docID string
	var err error

	if req.IsSplitAudio {
		// Ensure parentDocID is provided for split audio documents
		if req.ParentDocID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent Document ID is required for split audio documents"})
			return
		}
		// Generate a random suffix for the child document ID
		suffix, err := GenerateRandomSuffix(3)
		if err != nil {
			return
		}
		childDocID := req.ParentDocID + "-" + suffix
		// Directly use Document from req for creating a child document
		docID, err = db.CreateDocument(c, req.IsSplitAudio, req.ParentDocID, childDocID, req.Document)
	} else {
		// Directly use ParentDocument from req for creating a parent document
		// ParentDocID is generated inside the CreateDocument function for parent documents
		docID, err = db.CreateDocument(c, req.IsSplitAudio, "", "", req.ParentDocument)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document created successfully", "document_id": docID})
}

func UpdateDocument(c *gin.Context) {
	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var err error

	if req.IsSplitAudio {
		if err := c.ShouldBindJSON(&req.Document); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		if err := db.UpdateDocument(c, req.IsSplitAudio, req.ParentDocID, req.DocID, req.Document); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
	} else {
		if err := c.ShouldBindJSON(&req.ParentDocument); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		if err := db.UpdateDocument(c, req.IsSplitAudio, req.ParentDocID, "", req.ParentDocument); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
	}
}

func GetDocumentByID(c *gin.Context) {
	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if req.IsSplitAudio {
		if err := c.ShouldBindJSON(&req.Document); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		doc, err := db.GetDocumentByID(c, req.IsSplitAudio, req.ParentDocID, req.DocID)
		if err != nil {
			if strings.Contains(err.Error(), "document with ID") && strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Server Error"})
			return
		}
		c.JSON(http.StatusOK, doc)
	} else {
		if err := c.ShouldBindJSON(&req.ParentDocument); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		parentDoc, err := db.GetDocumentByID(c, req.IsSplitAudio, req.ParentDocID, "")
		if err != nil {
			if strings.Contains(err.Error(), "document with ID") && strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Server Error"})
			return
		}
		c.JSON(http.StatusOK, parentDoc)
	}
}

func SearchDocuments(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "details": err.Error()})
		return
	}

	docs, err := db.SearchDocuments(c, req.IsSplitAudio, req.ParentDocID, req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to search documents", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

func GetAllDocuments(c *gin.Context) {
	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	docs, err := db.GetAllDocuments(c, req.ParentDocID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents", "details": err.Error()})
		return
	}

	fmt.Println(docs)
	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

func DeleteDocument(c *gin.Context) {
	docID := c.Param("id") // Assuming the document ID is passed as a URL parameter

	if docID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	// parentDocID is optional and used to indicate if the document is a child
	parentDocID := c.Query("parentDocID")

	if err := db.DeleteDocument(c, docID, parentDocID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
