package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// GET /evento/:eventoId/perfiles
func (a *Api) ListarPerfilesPorEvento(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	resp, e := a.BllController.PerfilPersona.ListarPerfilesPorEvento(eventoID)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}

// POST /evento/:eventoId/perfiles
func (a *Api) CrearPerfilPersona(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoID, err := strconv.ParseInt(eventoIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.PerfilPersonaRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	req.EventoID = eventoID

	usuarioCreacion := int64(1) // luego lo cambias por el ID del token

	resp, e := a.BllController.PerfilPersona.CrearPerfilPersona(req, usuarioCreacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusCreated, resp)
}

// PUT /perfiles/:perfilId
func (a *Api) ActualizarPerfilPersona(c echo.Context) error {
	perfilIdStr := c.Param("perfilId")
	perfilID, err := strconv.ParseInt(perfilIdStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.PerfilPersonaUpdateRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuarioModificacion := int64(1)

	resp, e := a.BllController.PerfilPersona.ActualizarPerfilPersona(perfilID, req, usuarioModificacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}
