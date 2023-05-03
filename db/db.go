package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"strings"

	"google.golang.org/api/option"
)

type Document struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

var client *firestore.Client

// InitializeFirestoreClient initializes the Firestore client
func InitializeFirestoreClient(ctx context.Context, projectID string) error {
	creds := option.WithCredentialsFile("server/salafifatawa-firestore.json")
	firestoreClient, err := firestore.NewClient(ctx, projectID, creds)
	if err != nil {
		return err
	}

	client = firestoreClient
	return nil
}

// CreateDocument creates a new document in Firestore
func CreateDocument(ctx context.Context, doc Document) (string, error) {
	docRef, _, err := client.Collection("salafifatawa").Add(ctx, map[string]interface{}{
		"title":    doc.Title,
		"author":   doc.Author,
		"question": doc.Question,
		"answer":   doc.Answer,
	})

	if err != nil {
		return "", err
	}

	return docRef.ID, nil
}

// GetDocumentByID retrieves a single document by ID
func GetDocumentByID(ctx context.Context, docID string) (Document, error) {
	// Get document reference
	docRef := client.Collection("salafifatawa").Doc(docID)

	// Get document data
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return Document{}, err
	}

	// Parse document data into Document struct
	docData := docSnapshot.Data()
	doc := Document{
		ID:       docSnapshot.Ref.ID,
		Title:    docData["title"].(string),
		Author:   docData["author"].(string),
		Question: docData["question"].(string),
		Answer:   docData["answer"].(string),
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
		if err == iterator.Done {
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
