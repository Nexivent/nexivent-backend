package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// GET /evento/:eventoId/tipos-ticket
func (a *Api) ListarTiposTicketPorEvento(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	resp, e := a.BllController.TipoTicket.ListarTiposTicketPorEvento(eventoID)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}

// POST /evento/:eventoId/tipos-ticket
func (a *Api) CrearTipoTicket(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.TipoTicketTicketRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	req.EventoID = eventoID

	usuarioCreacion := int64(1)

	resp, e := a.BllController.TipoTicket.CrearTipoTicket(req, usuarioCreacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusCreated, resp)
}

// PUT /tipos-ticket/:tipoTicketId
func (a *Api) ActualizarTipoTicket(c echo.Context) error {
	idStr := c.Param("tipoTicketId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.TipoTicketUpdateRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuarioModificacion := int64(1)

	resp, e := a.BllController.TipoTicket.ActualizarTipoTicket(id, req, usuarioModificacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}
