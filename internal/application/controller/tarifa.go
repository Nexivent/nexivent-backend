package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type TarifaController struct {
	Logger  logging.Logger
	Adapter *adapter.TarifaAdapter
}

func NewTarifaController(
	logger logging.Logger,
	a *adapter.TarifaAdapter,
) *TarifaController {
	return &TarifaController{
		Logger:  logger,
		Adapter: a,
	}
}

func (c *TarifaController) CrearTarifa(req schemas.TarifaRequest, usuarioCreacion int64) (*schemas.TarifaResponse, *errors.Error) {
	return c.Adapter.CrearTarifa(&req, usuarioCreacion)
}

func (c *TarifaController) ActualizarTarifa(id int64, req schemas.TarifaUpdateRequest, usuarioModificacion int64) (*schemas.TarifaResponse, *errors.Error) {
	return c.Adapter.ActualizarTarifa(id, &req, usuarioModificacion)
}

func (c *TarifaController) ListarTarifasPorIDs(ids []int64) ([]schemas.TarifaResponse, *errors.Error) {
	return c.Adapter.ListarTarifasPorIDs(ids)
}
