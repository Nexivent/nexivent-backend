package controller

import (
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
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

func (ec *EventoController) FetchEventosFeed(usuarioId *int64) (*schemas.EventosPaginados, *errors.Error) {
	return ec.EventoAdapter.FetchPostgresqlEventosFeed(usuarioId)
}
func (ec *EventoController) FetchEventosConInteraccionesFeed(usuarioId *int64) (*schemas.EventosPaginados, *errors.Error) {
	return ec.EventoAdapter.FetchPostgresqlEventosConInteraccionesFeed(usuarioId)
}

func (ec *EventoController) FetchEventosWithFilters(
	categoriaID *int64,
	organizadorID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time,
	estado *int16,
	soloFuturos bool) (*schemas.EventosPaginados, *errors.Error) {
	return ec.EventoAdapter.FetchPostgresqlEventosWithFilters(
		categoriaID,
		organizadorID,
		titulo,
		descripcion,
		lugar,
		fecha,
		horaInicio,
		estado,
		soloFuturos)
}

// GetEventoById retrieves an event by its ID with all related entities
func (ec *EventoController) GetEventoById(eventoID int64) (*schemas.EventoResponse, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlEventoById(eventoID)
}

// GetReporteEvento genera el reporte de un evento por ID + filtros opcionales.
func (ec *EventoController) GetReporteEvento(organizadorId int64,
	eventoID *int64,
	fechaDesde *time.Time,
	fechaHasta *time.Time,
) ([]*schemas.EventoReporte, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlReporteEvento(organizadorId, eventoID, fechaDesde, fechaHasta)
}

// GetReporteEventosOrganizador genera el reporte resumido de todos los eventos de un organizador.
func (ec *EventoController) GetReporteEventosOrganizador(
	organizadorID int64,
	fechaDesde *time.Time,
	fechaHasta *time.Time,
) ([]schemas.EventoOrganizadorReporte, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlReporteEventosOrganizador(organizadorID, fechaDesde, fechaHasta)
}

// GenerarReporteAdministrativo genera el reporte global BI para administradores
func (ec *EventoController) GenerarReporteAdministrativo(req schemas.AdminReportRequest) (*schemas.AdminReportResponse, *errors.Error) {
	return ec.EventoAdapter.GenerarReporteAdministrativo(req)
}

func (ec *EventoController) GetEventoDetalle(eventoId int64) (*schemas.EventoDetalleDTO, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlEventoDetalle(eventoId)
}

func (ec *EventoController) EditarEventoFull(eventoID int64, req schemas.EditarEventoFullRequest) (*schemas.EventoResponse, *errors.Error) {
	return ec.EventoAdapter.EditarEventoFull(eventoID, &req)
}

func (c *EventoController) EditarEvento(req *schemas.EditarEventoRequest) (*schemas.EventoDetalleDTO, *errors.Error) {
	return c.EventoAdapter.EditarEvento(req)
}

func (ec *EventoController) ObtenerTransaccionesPorEvento(eventoId int64) ([]model.OrdenDeCompra, *errors.Error) {
	return ec.EventoAdapter.ObtenerTransaccionesPorEvento(eventoId)
}

func (ec *EventoController) PostInteraccionUsuarioEvento(req schemas.InteraccionConEventoRequest) (*schemas.InteraccionConEventoResponse, *errors.Error) {
	return ec.EventoAdapter.PostPostgresqlInteraccionUsuarioEvento(req)
}

func (ec *EventoController) PutInteraccionUsuarioEvento(req schemas.InteraccionConEventoRequest) (*schemas.InteraccionConEventoResponse, *errors.Error) {
	return ec.EventoAdapter.PutPostgresqlInteraccionUsuarioEvento(req)
}

func (ec *EventoController) GetAsistentesPorEvento(eventoID int64) ([]map[string]interface{}, *errors.Error) {
	// Validar que el evento existe
	_, err := ec.EventoAdapter.GetPostgresqlEventoById(eventoID)
	if err != nil {
		ec.Logger.Errorf("❌ [CONTROLLER] Evento no encontrado: %d", eventoID)
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	asistentes, err := ec.EventoAdapter.GetAsistentesPorEvento(eventoID)
	if err != nil {
		return nil, err
	}

	return asistentes, nil
}
