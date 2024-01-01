package auth

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"net/http"
	"strings"
)

var authClient *auth.Client

// LoginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response from the Login endpoint
type LoginResponse struct {
	Token string `json:"token"`
}

type LoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// InitializeAuthClient initializes the Firebase Auth client
func InitializeAuthClient(ctx context.Context) error {
	opt := option.WithCredentialsFile("backend/salafifatawa-firestore.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	authClient, err = app.Auth(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AuthenticateUser is a middleware function that verifies the user's ID token
func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID token from Authorization header
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		fmt.Println("idToken: ", idToken)

		// Verify ID token
		token, err := authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		// Check if token is nil before accessing its properties
		if token == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			ctx.Abort()
			return
		}
		fmt.Println("token: ", token.UID)

		// Set user ID as context value
		ctx.Set("userID", token.UID)
		ctx.Next()
	}
}
