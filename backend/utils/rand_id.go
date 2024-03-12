package utils

import (
	"github.com/google/uuid"
	"strings"
)

func generateID() string {
	// Generate a UUID
	uniqueID := uuid.New()
	// Convert UUID to string and remove hyphens
	safeID := strings.ReplaceAll(uniqueID.String(), "-", "")
	return safeID
}
