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

func (a *Api) GetRolPorNombre(c echo.Context) error {
	nombreStr := c.Param("nombre")

	response, err := a.BllController.Rol.GetRolPorNombre(nombreStr)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}


func (a *Api) GetRolPorUsuario(c echo.Context) error {
	usuarioIdStr := c.Param("usuarioId")
	usuarioId, parseErr := strconv.ParseInt(usuarioIdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.Rol.GetRolPorUsuario(usuarioId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Roles.
// @Description 		Fetches all available categorias.
// @Tags 				Categoria
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} schemas.CategoriaResponse "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/categorias/ [get]
func (a *Api) FetchRoles(c echo.Context) error {
	response, err := a.BllController.Rol.FetchRoles()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}


// @Summary 			Update Community.
// @Description 		Update the community information.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Param               request body schemas.UpdateCommunityRequest true "Update Community Request"
// @Success 			200 {object} schemas.Community "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/{communityId}/ [patch]
func (a *Api) UpdateRol(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)

	var rolID *int64
	if catStr := c.QueryParam("rolId"); catStr != "" {
		parsed, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		rolID = &parsed
	}

	var updatedBy *int64
	if catStr := c.QueryParam("userId"); catStr != "" {
		parsed, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		updatedBy = &parsed
	}

	var request *schemas.RolRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	request.UsuarioModificacion=updatedBy
	response, newErr := a.BllController.Rol.ActualizarRol(*request, *rolID)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}
