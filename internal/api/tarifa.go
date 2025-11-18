package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// POST /tarifas
func (a *Api) CrearTarifa(c echo.Context) error {
	var req schemas.TarifaRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuarioCreacion := int64(1)

	resp, e := a.BllController.Tarifa.CrearTarifa(req, usuarioCreacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusCreated, resp)
}

// PUT /tarifas/:tarifaId
func (a *Api) ActualizarTarifa(c echo.Context) error {
	idStr := c.Param("tarifaId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var req schemas.TarifaUpdateRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuarioModificacion := int64(1)

	resp, e := a.BllController.Tarifa.ActualizarTarifa(id, req, usuarioModificacion)
	if e != nil {
		return errors.HandleError(*e, c)
	}
	return c.JSON(http.StatusOK, resp)
}
