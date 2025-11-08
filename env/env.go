package env

import "os"

// GetVar gives the value of an environment variable or fallbacks to a default value.
func GetVar(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
