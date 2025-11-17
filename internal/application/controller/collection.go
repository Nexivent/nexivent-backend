package controller

import (
	"gorm.io/gorm"

	"github.com/Loui27/nexivent-backend/internal/application/adapter"
	"github.com/Loui27/nexivent-backend/internal/application/service/storage"
	config "github.com/Loui27/nexivent-backend/internal/config"
	"github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/logging"
)

type ControllerCollection struct {
	Logger    logging.Logger
	Evento    *EventoController
	Categoria *CategoriaController
	Media     *MediaController
	Cupon     *CuponController
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
	categoriaAdapter := adapter.NewCategoriaAdapter(logger, daoPostgresql)
	cuponAdapter := adapter.NewCuponAdapter(logger, daoPostgresql)

	// Services
	s3Storage, storageErr := storage.NewS3Storage(logger, configEnv)
	if storageErr != nil {
		logger.Warnln("S3 storage not initialized:", storageErr)
	}

	// Create controllers
	eventoController := NewEventoController(logger, eventoAdapter)
	categoriaController := NewCategoriaController(logger, categoriaAdapter)
	cuponController := NewCuponController(logger, cuponAdapter)

	var mediaController *MediaController
	if s3Storage != nil {
		mediaController = NewMediaController(logger, s3Storage)
	}

	return &ControllerCollection{
		Logger:    logger,
		Evento:    eventoController,
		Categoria: categoriaController,
		Media:     mediaController,
		Cupon:     cuponController,
	}, nexiventPsqlDB
}
