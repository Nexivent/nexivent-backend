package env

import (
	"os"
	"path/filepath"
)

// Find .env path iterating through parent folders
func FindEnvPath() (string, error) {
	dir, _ := os.Getwd()

	for {
		// Check if the .env file exists in the current directory
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		// Move up to the parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", os.ErrNotExist
		}
		dir = parentDir
	}
}
