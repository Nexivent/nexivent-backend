package controller

import (
	"gorm.io/gorm"

	config "github.com/Loui27/nexivent-backend/internal/config"
	"github.com/Loui27/nexivent-backend/internal/application/adapter"
	"github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/logging"
)

type ControllerCollection struct {
	Logger  logging.Logger
	Evento  *EventoController
}

// Creates BLL controller collection
func NewControllerCollection(
	logger logging.Logger,
	configEnv *config.ConfigEnv,
) (*ControllerCollection, *gorm.DB) {
	// Create DAO layer
	daoPostgresql, nexiventPsqlDB := repository.NewNexiventPsqlEntidades(
		logger,
		configEnv,
	)

	// Create adapters
	eventoAdapter := adapter.NewEventoAdapter(logger, daoPostgresql)

	// Create controllers
	eventoController := NewEventoController(logger, eventoAdapter)

	return &ControllerCollection{
		Logger:  logger,
		Evento:  eventoController,
	}, nexiventPsqlDB
}
