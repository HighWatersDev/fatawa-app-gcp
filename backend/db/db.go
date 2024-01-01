package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"

	"google.golang.org/api/option"
)

type Document struct {
	ID       string `json:"id"`
	Audio    string `json:"audio"`
	Title    string `json:"title"`
	Topic    string `json:"topic"`
	Author   string `json:"author"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Complete bool   `json:"complete"`
}

var client *firestore.Client

// InitializeFirestoreClient initializes the Firestore client
func InitializeFirestoreClient(ctx context.Context, projectID string) error {
	creds := option.WithCredentialsFile("backend/salafifatawa-firestore.json")
	firestoreClient, err := firestore.NewClient(ctx, projectID, creds)
	if err != nil {
		return err
	}

	client = firestoreClient
	return nil
}

// CreateDocument creates a new document in Firestore
func CreateDocument(ctx context.Context, docID string, doc Document) error {
	docRef := client.Collection("salafifatawa").Doc(docID) // Create a document reference with a specified ID

	_, err := docRef.Set(ctx, map[string]interface{}{
		"audio":    doc.Audio,
		"title":    doc.Title,
		"topic":    doc.Topic,
		"author":   doc.Author,
		"question": doc.Question,
		"answer":   doc.Answer,
		"complete": doc.Complete,
	})

	if err != nil {
		return err
	}

	return nil
}

// UpdateDocument updates an existing document in Firestore
func UpdateDocument(ctx context.Context, docID string, updatedDoc Document) error {
	_, err := client.Collection("salafifatawa").Doc(docID).Set(ctx, map[string]interface{}{
		"audio":    updatedDoc.Audio,
		"title":    updatedDoc.Title,
		"topic":    updatedDoc.Topic,
		"author":   updatedDoc.Author,
		"question": updatedDoc.Question,
		"answer":   updatedDoc.Answer,
		"complete": updatedDoc.Complete,
	}, firestore.MergeAll)

	if err != nil {
		log.Printf("Failed to update document: %v", err)
		return err
	}

	log.Printf("Document with ID %s updated successfully", docID)
	return nil
}

// GetDocumentByID retrieves a single document by ID
func GetDocumentByID(ctx context.Context, docID string) (Document, error) {
	// Get document reference
	docRef := client.Collection("salafifatawa").Doc(docID)

	// Get document data
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		// Check if the error is because the document was not found
		if status.Code(err) == codes.NotFound {
			return Document{}, fmt.Errorf("document with ID %s not found: %w", docID, err)
		}
		return Document{}, err
	}

	// Ensure the document snapshot exists
	if !docSnapshot.Exists() {
		return Document{}, fmt.Errorf("document with ID %s not found", docID)
	}

	// Parse document data into Document struct
	docData := docSnapshot.Data()
	doc := Document{
		ID:       docSnapshot.Ref.ID,
		Audio:    docData["audio"].(string),
		Title:    docData["title"].(string),
		Topic:    docData["topic"].(string),
		Author:   docData["author"].(string),
		Question: docData["question"].(string),
		Answer:   docData["answer"].(string),
		Complete: docData["complete"].(bool),
	}

	return doc, nil
}

func SearchDocuments(c context.Context, searchQuery string) ([]Document, error) {
	var documents []Document

	q := client.Collection("salafifatawa").
		OrderBy("title", firestore.Asc).
		StartAt(strings.ToLower(searchQuery)).
		EndAt(strings.ToLower(searchQuery + "\uf8ff"))

	iter := q.Documents(c)

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return documents, err
		}

		var d Document
		err = doc.DataTo(&d)
		if err != nil {
			return documents, err
		}

		d.ID = doc.Ref.ID

		documents = append(documents, d)
	}

	return documents, nil
}

func GetAllDocuments(c context.Context) ([]Document, error) {
	var documents []Document

	iter := client.Collection("salafifatawa").Documents(c)

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return documents, err
		}

		var d Document
		err = doc.DataTo(&d)
		if err != nil {
			return documents, err
		}

		// Set the ID field to the document ID
		d.ID = doc.Ref.ID

		documents = append(documents, d)
	}

	return documents, nil
}

// DeleteDocument deletes a document in Firestore by ID
func DeleteDocument(ctx context.Context, docID string) error {
	_, err := client.Collection("salafifatawa").Doc(docID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete document with ID %s: %v", docID, err)
		return err
	}

	log.Printf("Document with ID %s deleted successfully", docID)
	return nil
}
