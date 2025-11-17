package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	adapter "github.com/Nexivent/nexivent-backend/internal/application/adapter"
	schemas "github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type OrdenDeCompraController struct {
	Logger       logging.Logger
	OrdenAdapter *adapter.OrdenDeCompra
}

func NewOrdenDeCompraController(
	logger logging.Logger,
	ordenAdapter *adapter.OrdenDeCompra,
) *OrdenDeCompraController {
	return &OrdenDeCompraController{
		Logger:       logger,
		OrdenAdapter: ordenAdapter,
	}
}

// POST /api/orders/hold
func (oc *OrdenDeCompraController) CrearSesionOrdenTemporal(
	req schemas.CrearOrdenTemporalRequest,
) (*schemas.CrearOrdenTemporalResponse, *errors.Error) {
	return oc.OrdenAdapter.CrearSesionOrdenTemporal(&req)
}

// GET /api/orders/{orderId}/hold
func (oc *OrdenDeCompraController) ObtenerEstadoHold(
	orderID int64,
) (*schemas.ObtenerHoldResponse, *errors.Error) {
	return oc.OrdenAdapter.ObtenerEstadoHold(orderID)
}

// POST /api/orders/{orderId}/confirm
func (oc *OrdenDeCompraController) ConfirmarOrden(
	orderID int64,
	req schemas.ConfirmarOrdenRequest,
) (*schemas.ConfirmarOrdenResponse, *errors.Error) {
	return oc.OrdenAdapter.ConfirmarOrden(orderID, &req)
}
