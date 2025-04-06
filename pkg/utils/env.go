package utils

import "os"

// GetEnvOrDefault returns the value of the environment variable or the default value if the variable is not set
func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
