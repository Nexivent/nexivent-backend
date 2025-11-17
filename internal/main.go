package main

import (
	"github.com/Nexivent/nexivent-backend/internal/api"
	_ "github.com/Nexivent/nexivent-backend/internal/api/docs"
	"github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/logging"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger("NexiventBackendServer", "Version 1.0", logging.FormatText, 4)

	// Load configuration from environment
	configEnv := config.NuevoConfigEnv(logger)

	// Run the API service
	api.RunService(configEnv, logger)
}
