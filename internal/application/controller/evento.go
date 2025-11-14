package controller

import (
	"github.com/Loui27/nexivent-backend/errors"
	"github.com/Loui27/nexivent-backend/internal/application/adapter"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
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

// GetEventoById retrieves an event by its ID with all related entities
func (ec *EventoController) GetEventoById(eventoID int64) (*schemas.EventoResponse, *errors.Error) {
	return ec.EventoAdapter.GetPostgresqlEventoById(eventoID)
}
