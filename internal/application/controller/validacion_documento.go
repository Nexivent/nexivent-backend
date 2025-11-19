package controller

import (
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type ValidacionDocumentoController struct {
	validacionAdapter *adapter.ValidacionDocumento
	logger            logging.Logger
}

func NewValidacionDocumentoController(
	validacionAdapter *adapter.ValidacionDocumento,
	logger logging.Logger,
) *ValidacionDocumentoController {
	return &ValidacionDocumentoController{
		validacionAdapter: validacionAdapter,
		logger:            logger,
	}
}

func (c *ValidacionDocumentoController) ValidarDocumento(req *schemas.ValidarDocumentoRequest) (*schemas.ValidarDocumentoResponse, error) {
	c.logger.Info("Controller: Validando documento", req.TipoDocumento, req.NumeroDocumento)
	return c.validacionAdapter.ValidarDocumento(req)
}