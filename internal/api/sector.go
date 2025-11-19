package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// GET /evento/:eventoId/sectores
func (a *Api) ListarSectoresPorEvento(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	resp, e := a.BllController.Sector.ListarSectoresPorEvento(eventoID)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}

// POST /evento/:eventoId/sectores
func (a *Api) CrearSector(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.SectorTicketRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	req.EventoID = eventoID

	usuarioCreacion := int64(1)

	resp, e := a.BllController.Sector.CrearSector(req, usuarioCreacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusCreated, resp)
}

// PUT /sectores/:sectorId
func (a *Api) ActualizarSector(c echo.Context) error {
	sectorIdStr := c.Param("sectorId")
	sectorID, err := strconv.ParseInt(sectorIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.SectorUpdateRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuarioModificacion := int64(1)

	resp, e := a.BllController.Sector.ActualizarSector(sectorID, req, usuarioModificacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}
