package utility

import (
	"strings"
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
func GetClaim(key string, claims map[string]interface{}) string {

	if val, ok := claims[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
