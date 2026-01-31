package utility

import (
	"strings"

	"cloud.google.com/go/auth/credentials/idtoken"
)

func SplitNameStrict(fullName string) (string, string) {
	parts := strings.Fields(fullName)

	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// Helper function to safely get strings from claims
func GetClaim(key string, payload *idtoken.Payload) string {
	if val, ok := payload.Claims[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
