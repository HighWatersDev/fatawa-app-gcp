package main

import (
	"context"
	"fatawa-app-gcp/backend/auth"
	"fatawa-app-gcp/backend/db"
	"fatawa-app-gcp/backend/server"
	"log"
)

func main() {
	ctx := context.Background()

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

	// Start server
	log.Println("Starting server...")
	err = server.StartServer(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
