package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type TipoTicketController struct {
	Logger  logging.Logger
	Adapter *adapter.TipoTicketAdapter
}

func NewTipoTicketController(
	logger logging.Logger,
	a *adapter.TipoTicketAdapter,
) *TipoTicketController {
	return &TipoTicketController{
		Logger:  logger,
		Adapter: a,
	}
}

func (c *TipoTicketController) CrearTipoTicket(req schemas.TipoTicketTicketRequest, usuarioCreacion int64) (*schemas.TipoTicketTicketResponse, *errors.Error) {
	return c.Adapter.CrearTipoTicket(&req, usuarioCreacion)
}

func (c *TipoTicketController) ActualizarTipoTicket(id int64, req schemas.TipoTicketUpdateRequest, usuarioModificacion int64) (*schemas.TipoTicketTicketResponse, *errors.Error) {
	return c.Adapter.ActualizarTipoTicket(id, &req, usuarioModificacion)
}

func (c *TipoTicketController) ListarTiposTicketPorEvento(eventoID int64) ([]schemas.TipoTicketTicketResponse, *errors.Error) {
	return c.Adapter.ListarTiposTicketPorEvento(eventoID)
}
