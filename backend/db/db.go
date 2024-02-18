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
}

type ParentDocument struct {
	ID       string `json:"id"`
	Audio    string `json:"audio"`
	Author   string `json:"author"`
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
func CreateDocument(ctx context.Context, isSplitAudio bool, parentDocID string, docID string, docData interface{}) (string, error) {
	var docRef *firestore.DocumentRef

	if isSplitAudio {
		// Convert docData to Document type
		doc, ok := docData.(Document)
		if !ok {
			return "", fmt.Errorf("invalid document data for split audio")
		}

		docRef = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").Doc(docID)
		_, err := docRef.Set(ctx, map[string]interface{}{
			"audio":    doc.Audio,
			"title":    doc.Title,
			"topic":    doc.Topic,
			"author":   doc.Author,
			"question": doc.Question,
			"answer":   doc.Answer,
		})
		if err != nil {
			return "", err
		}
	} else {
		// Convert docData to ParentDocument type
		parentDoc, ok := docData.(ParentDocument)
		if !ok {
			return "", fmt.Errorf("invalid document data for parent audio")
		}
		docRef = client.Collection("salafifatawa").NewDoc()
		_, err := docRef.Set(ctx, map[string]interface{}{
			"audio":    parentDoc.Audio,
			"author":   parentDoc.Author,
			"complete": parentDoc.Complete,
		})
		if err != nil {
			return "", err
		}
	}

	return docRef.ID, nil
}

// UpdateDocument updates an existing document in Firestore
func UpdateDocument(ctx context.Context, isSplitAudio bool, parentDocID string, docID string, docData interface{}) error {
	var docRef *firestore.DocumentRef

	if isSplitAudio {

		doc, ok := docData.(Document)
		if !ok {
			return fmt.Errorf("invalid document data for split audio")
		}
		docRef = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").Doc(docID)
		_, err := docRef.Set(ctx, map[string]interface{}{
			"audio":    doc.Audio,
			"title":    doc.Title,
			"topic":    doc.Topic,
			"author":   doc.Author,
			"question": doc.Question,
			"answer":   doc.Answer,
		}, firestore.MergeAll)
		if err != nil {
			log.Printf("Failed to update document: %v", err)
			return err
		}

		log.Printf("Document with ID %s updated successfully", docID)
		return nil
	} else {
		// Assume updateData is of type ParentDocument (for parent audio)
		parentDoc, ok := docData.(ParentDocument)
		if !ok {
			return fmt.Errorf("invalid document data for parent audio")
		}
		docRef = client.Collection("salafifatawa").Doc(docID)
		_, err := docRef.Set(ctx, map[string]interface{}{
			"audio":    parentDoc.Audio,
			"author":   parentDoc.Author,
			"complete": parentDoc.Complete,
		}, firestore.MergeAll)
		if err != nil {
			log.Printf("Failed to update document: %v", err)
			return err
		}

		log.Printf("Document with ID %s updated successfully", docID)
		return nil
	}
}

func GetDocumentByID(ctx context.Context, isSplitAudio bool, parentDocID string, docID string) (interface{}, error) {
	var docRef *firestore.DocumentRef

	if isSplitAudio {
		docRef = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").Doc(docID)
	} else {
		docRef = client.Collection("salafifatawa").Doc(docID)
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		// Check if the error is because the document was not found
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("document with ID %s not found: %w", docID, err)
		}
		return nil, err
	}

	if !docSnapshot.Exists() {
		return nil, fmt.Errorf("document with ID %s not found", docID)
	}

	// Dynamically choose the type to deserialize based on isSplitAudio
	if isSplitAudio {
		var doc Document
		err = docSnapshot.DataTo(&doc)
		if err != nil {
			return nil, err
		}
		doc.ID = docSnapshot.Ref.ID // Ensure the document's ID is included
		return doc, nil
	} else {
		var parentDoc ParentDocument
		err = docSnapshot.DataTo(&parentDoc)
		if err != nil {
			return nil, err
		}
		parentDoc.ID = docSnapshot.Ref.ID // Ensure the document's ID is included
		return parentDoc, nil
	}
}

func SearchDocuments(ctx context.Context, isSplitAudio bool, parentDocID string, searchQuery string) ([]interface{}, error) {
	var iter *firestore.DocumentIterator

	if isSplitAudio {
		// Search in the 'split-audios' subcollection
		iter = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").
			OrderBy("title", firestore.Asc).
			StartAt(strings.ToLower(searchQuery)).
			EndAt(strings.ToLower(searchQuery + "\uf8ff")).
			Documents(ctx)
	} else {
		// Search in the 'salafifatawa' collection
		iter = client.Collection("salafifatawa").
			OrderBy("title", firestore.Asc).
			StartAt(strings.ToLower(searchQuery)).
			EndAt(strings.ToLower(searchQuery + "\uf8ff")).
			Documents(ctx)
	}

	defer iter.Stop()

	var documents []interface{}
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var document interface{}
		if isSplitAudio {
			var splitAudio Document
			err := doc.DataTo(&splitAudio)
			if err != nil {
				return nil, err
			}
			splitAudio.ID = doc.Ref.ID
			document = splitAudio
		} else {
			var parentAudio ParentDocument
			err := doc.DataTo(&parentAudio)
			if err != nil {
				return nil, err
			}
			parentAudio.ID = doc.Ref.ID
			document = parentAudio
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func GetAllDocuments(ctx context.Context, parentDocID string) ([]interface{}, error) {
	var documents []interface{}
	var iter *firestore.DocumentIterator

	if parentDocID != "" {
		iter = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").Documents(ctx)
	} else {
		iter = client.Collection("salafifatawa").Documents(ctx)
	}

	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var document interface{}
		if parentDocID != "" {
			// Deserialize into the child document structure
			var childDoc Document // Assuming Document is the struct for child docs
			err := doc.DataTo(&childDoc)
			if err != nil {
				return nil, err
			}
			childDoc.ID = doc.Ref.ID
			document = childDoc
		} else {
			// Deserialize into the parent document structure
			var parentDoc ParentDocument // Assuming ParentDocument is the struct for parent docs
			err := doc.DataTo(&parentDoc)
			if err != nil {
				return nil, err
			}
			parentDoc.ID = doc.Ref.ID
			document = parentDoc
		}

		documents = append(documents, document)
	}

	return documents, nil
}

// DeleteDocument deletes a document in Firestore by ID
func DeleteDocument(ctx context.Context, docID string, parentDocID string) error {
	var docRef *firestore.DocumentRef

	if parentDocID != "" {
		docRef = client.Collection("salafifatawa").Doc(parentDocID).Collection("split-audios").Doc(docID)
	} else {
		docRef = client.Collection("salafifatawa").Doc(docID)
	}

	_, err := docRef.Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete document with ID %s: %v", docID, err)
		return err
	}

	log.Printf("Document with ID %s deleted successfully", docID)
	return nil
}
