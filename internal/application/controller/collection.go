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
	Logger        logging.Logger
	Evento        *EventoController
	Usuario       *UsuarioController
	Categoria     *CategoriaController
	Media         *MediaController
	Cupon         *CuponController
	Comentario    *ComentarioController
	Orden         *OrdenDeCompraController
	PerfilPersona *PerfilPersonaController
	Sector        *SectorController
	TipoTicket    *TipoTicketController
	Tarifa        *TarifaController
	Ticket        *TicketController
	Token         *TokenController
	Rol           *RolController
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
	ordenAdapter := adapter.NewOrdenDeCompraAdapter(logger, daoPostgresql)
	perfilAdapter := adapter.NewPerfilPersonaAdapter(logger, daoPostgresql)
	sectorAdapter := adapter.NewSectorAdapter(logger, daoPostgresql)
	tipoTicketAdapter := adapter.NewTipoTicketAdapter(logger, daoPostgresql)
	tarifaAdapter := adapter.NewTarifaAdapter(logger, daoPostgresql)
	ticketAdapter := adapter.NewTicketAdapter(logger, daoPostgresql)
	rolAdapter := adapter.NewRolAdapter(logger, daoPostgresql)

	// Services
	s3Storage, storageErr := storage.NewS3Storage(logger, configEnv)
	if storageErr != nil {
		logger.Warnln("S3 storage not initialized:", storageErr)
	}

	// Create controllers
	eventoController := NewEventoController(logger, eventoAdapter)
	categoriaController := NewCategoriaController(logger, categoriaAdapter)
	cuponController := NewCuponController(logger, cuponAdapter)
	ordenController := NewOrdenDeCompraController(logger, ordenAdapter)
	perfilController := NewPerfilPersonaController(logger, perfilAdapter)
	sectorController := NewSectorController(logger, sectorAdapter)
	tipoTicketController := NewTipoTicketController(logger, tipoTicketAdapter)
	tarifaController := NewTarifaController(logger, tarifaAdapter)
	ticketController := NewTicketController(logger, ticketAdapter)
	rolController := NewRolController(logger, rolAdapter)

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
		Usuario: &UsuarioController{
			Logger: logger,
			DB:     daoPostgresql,
		},
		Comentario: &ComentarioController{
			Logger: logger,
			DB:     daoPostgresql,
		},
		Orden:         ordenController,
		PerfilPersona: perfilController,
		Sector:        sectorController,
		TipoTicket:    tipoTicketController,
		Tarifa:        tarifaController,
		Ticket:        ticketController,
		Token: &TokenController{
			Logger: logger,
			DB:     daoPostgresql,
		},
		Rol: rolController,
	}, nexiventPsqlDB
}
