package main

import (
	"github.com/Loui27/nexivent-backend/internal/api"
	"github.com/Loui27/nexivent-backend/internal/config"
	"github.com/Loui27/nexivent-backend/logging"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger("NexiventBackendServer", "Version 1.0", logging.FormatText, 4)

	// Load configuration from environment
	configEnv := config.NuevoConfigEnv(logger)

	// Run the API service
	api.RunService(configEnv, logger)
}
