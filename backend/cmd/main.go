package main

import (
	"context"
	"fatawa-app-gcp/backend/router"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a base context
	ctx := context.Background()

	// Initialize router
	r := router.SetupRouter(ctx)
	if r == nil {
		log.Fatal("Failed to setup router")
	}

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Printf("Server error: %v", err)
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Println("Server exiting")
}