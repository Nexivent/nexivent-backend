package adapter

import (
	//"fmt"
	//"time"

	"github.com/Loui27/nexivent-backend/errors"
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
	//"github.com/Loui27/nexivent-backend/utils/convert"
	//"gorm.io/gorm"
)

type Categoria struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewCategoriaAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Categoria {
	return &Categoria{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// CreatePostgresqlEvento creates a new event with all related entities
func (e *Categoria) CreatePostgresqlCategoria(categoriaReq *schemas.CategoriaRequest) (*schemas.CategoriaResponse, *errors.Error) {
	// Start a transaction
	tx := e.DaoPostgresql.Categoria.PostgresqlDB.Begin()
	if tx.Error != nil {
		e.logger.Errorf("Failed to begin transaction: %v", tx.Error)
		return nil, &errors.BadRequestError.CategoriaNotCreated
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the main event model
	categoriaModel := &model.Categoria{
		Descripcion:     categoriaReq.Descripcion,
		Estado:          categoriaReq.Estado,
		Nombre:          categoriaReq.Nombre,
	}

	// Create the event
	if err := tx.Create(categoriaModel).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("Failed to create event: %v", err)
		return nil, &errors.BadRequestError.CategoriaNotCreated
	}

	
	// Build the response
	response := &schemas.CategoriaResponse{
		ID: categoriaModel.ID,
		Nombre: categoriaModel.Nombre,
		//Descripcion: categoriaModel.Descripcion,

		
	}

	return response, nil
}

// FetchPostgresqlEventos retrieves the upcoming events without filters
func (e *Categoria) FetchPostgresqlCategorias() ([] schemas.CategoriaResponse, *errors.Error) {
	categorias, err := e.DaoPostgresql.Categoria.ObtenerCategorias()
	if err != nil {
		e.logger.Errorf("Failed to fetch categorias: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	categoriasResponse := make([]schemas.CategoriaResponse, len(categorias))

	for i, c := range categorias{
		categoriasResponse[i] = schemas.CategoriaResponse{
			ID: c.ID,
			Nombre: c.Nombre,
		}
	}

	return categoriasResponse, nil
}

// GetPostgresqlEventoById gets an event by ID
/*
func (e *Evento) GetPostgresqlCategoriaById(categoriaID int64) (*schemas.CategoriaResponse, *errors.Error) {
	var eventoModel model.Evento

	// Use preload to fetch all related entities
	result := e.DaoPostgresql.Categoria.PostgresqlDB.
		Preload("Perfiles").
		Preload("Sectores").
		Preload("TiposTicket").
		Preload("Fechas.Fecha").
		First(&eventoModel, eventoID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.CategoriaNotFound
		}
		e.logger.Errorf("Failed to get evento: %v", result.Error)
		return nil, &errors.BadRequestError.CategoriaNotFound
	}

	// Build event dates response
	eventDates := []schemas.EventDateResponse{}
	for _, ef := range eventoModel.Fechas {
		if ef.Fecha != nil {
			eventDates = append(eventDates, schemas.EventDateResponse{
				IdFechaEvento: ef.ID,
				IdFecha:       ef.FechaID,
				Fecha:         ef.Fecha.FechaEvento.Format("2006-01-02"),
				HoraInicio:    ef.HoraInicio.Format("15:04"),
				HoraFin:       "", // You may need to add HoraFin to the model
			})
		}
	}

	// Build perfiles response
	perfiles := []schemas.PerfilResponse{}
	for _, p := range eventoModel.Perfiles {
		perfiles = append(perfiles, schemas.PerfilResponse{
			ID:    fmt.Sprintf("%d", p.ID),
			Label: p.Nombre,
		})
	}

	// Build sectores response
	sectores := []schemas.SectorResponse{}
	for _, s := range eventoModel.Sectores {
		sectores = append(sectores, schemas.SectorResponse{
			ID:        fmt.Sprintf("%d", s.ID),
			Nombre:    s.SectorTipo,
			Capacidad: s.TotalEntradas,
		})
	}

	// Build tipos ticket response
	tiposTicket := []schemas.TipoTicketResponse{}
	for _, t := range eventoModel.TiposTicket {
		tiposTicket = append(tiposTicket, schemas.TipoTicketResponse{
			ID:    fmt.Sprintf("%d", t.ID),
			Label: t.Nombre,
		})
	}

	// Fetch and build precios (this requires fetching tarifas)
	// For now, returning empty map - you may want to implement this based on your needs
	precios := make(schemas.PreciosSector)

	response := &schemas.EventoResponse{
		IdEvento:          eventoModel.ID,
		IdOrganizador:     eventoModel.OrganizadorID,
		IdCategoria:       eventoModel.CategoriaID,
		Titulo:            eventoModel.Titulo,
		Descripcion:       eventoModel.Descripcion,
		Lugar:             eventoModel.Lugar,
		Estado:            convert.MapEstadoToString(eventoModel.EventoEstado),
		Likes:             eventoModel.CantMeGusta,
		NoInteres:         eventoModel.CantNoInteresa,
		CantVendidasTotal: eventoModel.CantVendidoTotal,
		TotalRecaudado:    eventoModel.TotalRecaudado,
		ImagenPortada:     eventoModel.ImagenPortada,
		ImagenLugar:       eventoModel.ImagenEscenario,
		VideoUrl:          eventoModel.VideoPresentacion,
		EventDates:        eventDates,
		Perfiles:          perfiles,
		Sectores:          sectores,
		TiposTicket:       tiposTicket,
		Precios:           precios,
		Metadata: schemas.MetadataResponse{
			CreadoPor:           fmt.Sprintf("%d", *eventoModel.UsuarioCreacion),
			FechaCreacion:       eventoModel.FechaCreacion.Format(time.RFC3339),
			UltimaActualizacion: eventoModel.FechaCreacion.Format(time.RFC3339),
			Version:             1,
		},
	}

	return response, nil
}
*/