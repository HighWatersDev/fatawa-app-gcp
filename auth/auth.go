package auth

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
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
	opt := option.WithCredentialsFile("server/salafifatawa-firestore.json")
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

func Login(c *gin.Context) {
	var json LoginData
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authClient.GetUserByEmail(context.Background(), json.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	//// Verify the password
	//signInReq := &firebaseauth.SignInWithEmailAndPasswordRequest{
	//	Email:    json.Email,
	//	Password: json.Password,
	//}
	//_, err = authClient.SignInWithEmailAndPassword(context.Background(), signInReq)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
	//	return
	//}

	// Create a custom token
	customToken, err := authClient.CustomToken(context.Background(), user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": customToken})
}

// AuthenticateUser is a middleware function that verifies the user's ID token
func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID token from Authorization header
		authHeader := ctx.Request.Header.Get("Authorization")
		idToken := strings.Replace(authHeader, "Bearer ", "", 1)

		// Verify ID token
		token, err := authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		// Set user ID as context value
		ctx.Set("userID", token.UID)
		ctx.Next()
	}
}
