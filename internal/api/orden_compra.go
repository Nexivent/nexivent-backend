package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// -----------------------------------------------------------------------------
// POST /api/orders/hold
// -----------------------------------------------------------------------------

// @Summary      Crear sesión de compra temporal
// @Description  Crea una orden en estado TEMPORAL con expiración (hold).
// @Tags         Orden
// @Accept       json
// @Produce      json
// @Param        request body schemas.CrearOrdenTemporalRequest true "Datos de la reserva"
// @Success      201 {object} schemas.CrearOrdenTemporalResponse "Created"
// @Failure      400 {object} errors.Error "Bad Request"
// @Failure      409 {object} errors.Error "Conflict"
// @Failure      422 {object} errors.Error "Unprocessable Entity"
// @Failure      500 {object} errors.Error "Internal Server Error"
// @Router       /api/orders/hold [post]
func (a *Api) CrearSesionOrdenTemporal(c echo.Context) error {
	var req schemas.CrearOrdenTemporalRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// En producción, podrías sobreescribir req.IdUsuario con el user del token:
	// userID := c.Get("userID").(int64)
	// req.IdUsuario = userID

	resp, errBll := a.BllController.Orden.CrearSesionOrdenTemporal(req)
	if errBll != nil {
		return errors.HandleError(*errBll, c)
	}
	return c.JSON(http.StatusCreated, resp)
}

// -----------------------------------------------------------------------------
// GET /api/orders/{orderId}/hold
// -----------------------------------------------------------------------------

// @Summary      Obtener estado del hold de una orden
// @Description  Verifica si la reserva sigue activa y cuánto tiempo resta.
// @Tags         Orden
// @Accept       json
// @Produce      json
// @Param        orderId path int true "ID de la orden"
// @Success      200 {object} schemas.ObtenerHoldResponse "OK"
// @Failure      404 {object} errors.Error "Not Found"
// @Failure      410 {object} errors.Error "Gone"
// @Failure      422 {object} errors.Error "Unprocessable Entity"
// @Failure      500 {object} errors.Error "Internal Server Error"
// @Router       /api/orders/{orderId}/hold [get]
func (a *Api) ObtenerEstadoHold(c echo.Context) error {
	orderIdStr := c.Param("orderId")
	orderID, parseErr := strconv.ParseInt(orderIdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	resp, errBll := a.BllController.Orden.ObtenerEstadoHold(orderID)
	if errBll != nil {
		return errors.HandleError(*errBll, c)
	}
	return c.JSON(http.StatusOK, resp)
}

// -----------------------------------------------------------------------------
// POST /api/orders/{orderId}/confirm
// -----------------------------------------------------------------------------

// @Summary      Confirmar orden de compra
// @Description  Verifica el pago y actualiza la orden a estado CONFIRMADA.
// @Tags         Orden
// @Accept       json
// @Produce      json
// @Param        orderId path int true "ID de la orden"
// @Param        request body schemas.ConfirmarOrdenRequest true "Datos de pago"
// @Success      200 {object} schemas.ConfirmarOrdenResponse "OK"
// @Failure      400 {object} errors.Error "Bad Request"
// @Failure      402 {object} errors.Error "Payment Required"
// @Failure      409 {object} errors.Error "Conflict"
// @Failure      410 {object} errors.Error "Gone"
// @Failure      422 {object} errors.Error "Unprocessable Entity"
// @Failure      500 {object} errors.Error "Internal Server Error"
// @Router       /api/orders/{orderId}/confirm [post]
func (a *Api) ConfirmarOrden(c echo.Context) error {
	orderIdStr := c.Param("orderId")
	orderID, parseErr := strconv.ParseInt(orderIdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.ConfirmarOrdenRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	resp, errBll := a.BllController.Orden.ConfirmarOrden(orderID, req)
	if errBll != nil {
		return errors.HandleError(*errBll, c)
	}
	return c.JSON(http.StatusOK, resp)
}
