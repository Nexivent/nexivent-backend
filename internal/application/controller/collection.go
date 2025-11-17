package controller

import (
	"gorm.io/gorm"

	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/application/service/storage"
	config "github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/logging"
)

type ControllerCollection struct {
	Logger    logging.Logger
	Evento    *EventoController
	Usuario   *UsuarioController
	Categoria *CategoriaController
	Media     *MediaController
	Comentario *ComentarioController
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

	// Services
	s3Storage, storageErr := storage.NewS3Storage(logger, configEnv)
	if storageErr != nil {
		logger.Warnln("S3 storage not initialized:", storageErr)
	}

	// Create controllers
	eventoController := NewEventoController(logger, eventoAdapter)
	categoriaController := NewCategoriaController(logger, categoriaAdapter)
	var mediaController *MediaController
	if s3Storage != nil {
		mediaController = NewMediaController(logger, s3Storage)
	}

	return &ControllerCollection{
		Logger:    logger,
		Evento:    eventoController,
		Categoria: categoriaController,
		Media:     mediaController,
		Usuario: &UsuarioController{
			Logger: logger,
			DB:     daoPostgresql,
		},
		Comentario: &ComentarioController{
			Logger: logger,
			DB:     daoPostgresql,
		},
	}, nexiventPsqlDB
}
