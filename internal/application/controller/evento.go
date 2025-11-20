package controller

import (
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type EventoController struct {
	Logger        logging.Logger
	EventoAdapter *adapter.Evento
}

// NewEventoController creates a new controller for event operations
func NewEventoController(
	logger logging.Logger,
	eventoAdapter *adapter.Evento,
) *EventoController {
	return &EventoController{
		Logger:        logger,
		EventoAdapter: eventoAdapter,
	}
}

// CreateEvento creates a new event with all related entities
func (ec *EventoController) CreateEvento(
	eventoReq schemas.EventoRequest,
	usuarioCreacion int64,
) (*schemas.EventoResponse, *errors.Error) {
	return ec.EventoAdapter.CreatePostgresqlEvento(&eventoReq, usuarioCreacion)
}

// FetchEventos retrieves the list of available events
func (ec *EventoController) FetchEventos() (*schemas.EventosPaginados, *errors.Error) {
	return ec.EventoAdapter.FetchPostgresqlEventos()
}

func (ec *EventoController) FetchEventosWithFilters(
	categoriaID *int64,
	organizadorID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time) (*schemas.EventosPaginados, *errors.Error) {
	return ec.EventoAdapter.FetchPostgresqlEventosWithFilters(
		categoriaID,
		organizadorID,
		titulo,
		descripcion,
		lugar,
		fecha,
		horaInicio)
}

// GetEventoById retrieves an event by its ID with all related entities
func (ec *EventoController) GetEventoById(eventoID int64) (*schemas.EventoResponse, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlEventoById(eventoID)
}

// GetReporteEvento genera el reporte de un evento por ID + filtros opcionales.
func (ec *EventoController) GetReporteEvento(
	eventoID *int64,
	fechaDesde *time.Time,
	fechaHasta *time.Time,
) ([]*schemas.EventoReporte, *errors.Error) {

	return ec.EventoAdapter.GetPostgresqlReporteEvento(eventoID, fechaDesde, fechaHasta)
}

// GenerarReporteAdministrativo genera el reporte global BI para administradores
func (ec *EventoController) GenerarReporteAdministrativo(req schemas.AdminReportRequest) (*schemas.AdminReportResponse, *errors.Error) {
	return ec.EventoAdapter.GenerarReporteAdministrativo(req)
}
