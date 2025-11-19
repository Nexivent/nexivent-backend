package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type SectorController struct {
	Logger  logging.Logger
	Adapter *adapter.SectorAdapter
}

func NewSectorController(
	logger logging.Logger,
	a *adapter.SectorAdapter,
) *SectorController {
	return &SectorController{
		Logger:  logger,
		Adapter: a,
	}
}

func (c *SectorController) CrearSector(req schemas.SectorTicketRequest, usuarioCreacion int64) (*schemas.SectorTicketResponse, *errors.Error) {
	return c.Adapter.CrearSector(&req, usuarioCreacion)
}

func (c *SectorController) ActualizarSector(id int64, req schemas.SectorUpdateRequest, usuarioModificacion int64) (*schemas.SectorTicketResponse, *errors.Error) {
	return c.Adapter.ActualizarSector(id, &req, usuarioModificacion)
}

func (c *SectorController) ListarSectoresPorEvento(eventoID int64) ([]schemas.SectorTicketResponse, *errors.Error) {
	return c.Adapter.ListarSectoresPorEvento(eventoID)
}
