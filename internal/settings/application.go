package settings

import (
	"log/slog"
)

type Application struct {
	Config Config
	Logger *slog.Logger
}
