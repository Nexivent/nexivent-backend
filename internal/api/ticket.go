package api

import (
	"net/http"
	"strconv"
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// POST /api/tickets/issue

// @Summary      Emitir tickets para una orden confirmada
// @Description  Genera tickets individuales con código QR único
// @Tags         Ticket
// @Accept       json
// @Produce      json
// @Param        request body schemas.EmitirTicketsRequest true "Datos para emitir tickets"
// @Success      201 {object} schemas.EmitirTicketsResponse "Tickets generados"
// @Failure      404 {object} map[string]string "Orden no encontrada"
// @Failure      422 {object} errors.Error "Datos inválidos"
// @Failure      500 {object} errors.Error "Error interno"
// @Router       /api/tickets/issue [post]
func (a *Api) EmitirTickets(c echo.Context) error {
	var req schemas.EmitirTicketsRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	resp, ferr := a.BllController.Ticket.EmitirTicketsConInfo(req)
	if ferr != nil {
		if *ferr == errors.ObjectNotFoundError.EventoNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Orden no encontrada o no confirmada",
			})
		}
		return errors.HandleError(*ferr, c)
	}

	return c.JSON(http.StatusCreated, resp)
}

// POST /api/tickets/cancel

// @Summary      Cancelar uno o varios tickets.
// @Description  Cancela tickets (no USADOS ni ya CANCELADOS), actualiza stock y devuelve resumen.
// @Tags         Ticket
// @Accept       json
// @Produce      json
// @Param        request body schemas.TicketCancelRequest true "Cancelar Tickets Request"
// @Success      200 {object} schemas.TicketCancelResponse "OK"
// @Failure      404 {object} map[string]map[string]string "Error al cancelar los tickets"
// @Failure      422 {object} errors.Error "Unprocessable Entity"
// @Failure      500 {object} errors.Error "Internal Server Error"
// @Router       /api/tickets/cancel [post]
func (a *Api) CancelarTickets(c echo.Context) error {
	var req schemas.TicketCancelRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	resp, ferr := a.BllController.Ticket.CancelarTickets(req)
	if ferr != nil {
		if *ferr == errors.ObjectNotFoundError.EventoNotFound {
			return c.JSON(http.StatusNotFound, map[string]map[string]string{
				"error": {
					"message": "Ocurrió un error al cancelar los tickets.",
				},
			})
		}
		return errors.HandleError(*ferr, c)
	}

	return c.JSON(http.StatusOK, resp)
}


func (a *Api) GetTicketsByUser(c echo.Context) error {
	IdStr := c.Param("id")
	idUser, parseErr := strconv.ParseInt(IdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.Ticket.ObtenerTicketsPorUsuario(idUser)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}