package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type PerfilPersonaController struct {
	Logger  logging.Logger
	Adapter *adapter.PerfilPersonaAdapter
}

func NewPerfilPersonaController(
	logger logging.Logger,
	a *adapter.PerfilPersonaAdapter,
) *PerfilPersonaController {
	return &PerfilPersonaController{
		Logger:  logger,
		Adapter: a,
	}
}

func (c *PerfilPersonaController) CrearPerfilPersona(req schemas.PerfilPersonaRequest, usuarioCreacion int64) (*schemas.PerfilPersonaResponse, *errors.Error) {
	return c.Adapter.CrearPerfilPersona(&req, usuarioCreacion)
}

func (c *PerfilPersonaController) ActualizarPerfilPersona(id int64, req schemas.PerfilPersonaUpdateRequest, usuarioModificacion int64) (*schemas.PerfilPersonaResponse, *errors.Error) {
	return c.Adapter.ActualizarPerfilPersona(id, &req, usuarioModificacion)
}

func (c *PerfilPersonaController) ListarPerfilesPorEvento(eventoID int64) ([]schemas.PerfilPersonaResponse, *errors.Error) {
	return c.Adapter.ListarPerfilesPorEvento(eventoID)
}
