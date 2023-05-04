package server

import (
	"context"
	"fatawa-app-gcp/backend/router"
	"log"
)

func StartServer(ctx context.Context) error {
	// Initialize Gin router
	r := router.SetupRouter(ctx)

	// Start HTTP server
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
