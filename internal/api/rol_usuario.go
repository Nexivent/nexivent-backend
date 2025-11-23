package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// @Summary 			Get Rol.
// @Description 		Gets an event given its id.
// @Tags 				Rol
// @Accept 				json
// @Produce 			json
// @Param               nombre    path   string  true  "nombre"
// @Success 			200 {object} schemas.EventoResponse "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/evento/{eventoId}/ [get]

func (a *Api) ListarRolesDeUsuario(c echo.Context) error {
	userStr := c.Param("id")

	userId, parseErr := strconv.ParseInt(userStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.RolUsuario.GetUserRoles(userId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

func (a *Api) CreateRolUser(c echo.Context) error {

	var request schemas.RolUsuarioRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.RolUsuario.AsignarRolUser(request)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

func (a *Api) DeleteRolUser(c echo.Context) error {
	var request schemas.RolUsuarioRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.RolUsuario.RevokeRolUser(request)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

func (a *Api) ListarUsuariosPorRol(c echo.Context) error {
	var rol *int64
	
	// Si se proporciona el parámetro "rol", parsearlo
	if rolStr := c.QueryParam("rol"); rolStr != "" {
		parsed, err := strconv.ParseInt(rolStr, 10, 64)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		rol = &parsed
	}
	// Si rol es nil, significa que queremos TODOS los usuarios

	// Pasar el puntero directamente (puede ser nil)
	response, err := a.BllController.RolUsuario.GetUsersByRol(rol)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}