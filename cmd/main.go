package main

import (
	"context"
	"fatawa-app-gcp/server"
	"log"

	"fatawa-app-gcp/auth"
	"fatawa-app-gcp/db"
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
