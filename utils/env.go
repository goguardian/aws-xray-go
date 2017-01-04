package utils

import "os"

// GetenvOrDefault returns an environment variable value if not empty, otherwise
// it returns the default value.
func GetenvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return def
}
