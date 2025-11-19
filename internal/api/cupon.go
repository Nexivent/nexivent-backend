package api

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// @Summary 			Create Cupon.
// @Description 		Create a new cupon with validation and ownership.
// @Tags 				Cupon
// @Accept 				json
// @Produce 			json
// @Param               usuarioCreacion path int true "ID del usuario que crea el cupón"
// @Param               request body schemas.CuponResquest true "Create Cupon Request"
// @Success 			201 {object} schemas.CuponResponse "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/cupon/{usuarioCreacion} [post]
func (a *Api) CreateCupon(c echo.Context) error {
	usuarioCreacionParam := c.Param("usuarioCreacion")
	usuarioCreacionId, err := strconv.ParseInt(usuarioCreacionParam, 10, 64)

	if err != nil || usuarioCreacionId <= 0 {
		return errors.HandleError(errors.BadRequestError.InvalidUpdatedByValue, c)
	}

	var request schemas.CuponResquest
	result := c.Bind(&request)

	if result != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Cupon.CreateCupon(request, usuarioCreacionId)

	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Cupon.
// @Description 		Update an existing cupon with validation and ownership.
// @Tags 				Cupon
// @Accept 				json
// @Produce 			json
// @Param               usuarioModificacion path int true "ID del usuario que realiza la modificación"
// @Param               request body schemas.CuponResquest true "Update Cupon Request"
// @Success 			200 {object} schemas.CuponResponse "Updated"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/cupon/{usuarioModificacion} [put]
func (a *Api) UpdateCupon(c echo.Context) error {
	usuarioModParam := c.Param("usuarioModificacion")
	usuarioModId, err := strconv.ParseInt(usuarioModParam, 10, 64)

	if err != nil || usuarioModId <= 0 {
		return errors.HandleError(errors.BadRequestError.InvalidUpdatedByValue, c)
	}

	var request schemas.CuponResquest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Cupon.UpdateCupon(request, usuarioModId)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Cupones by Organizador.
// @Description 		Fetch all cupones belonging to a specific organizer.
// @Tags 				Cupon
// @Accept 				json
// @Produce 			json
// @Param               organizadorId path int true "ID del organizador"
// @Success 			200 {object} schemas.CuponesOrganizator "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/cupon/organizador/{organizadorId} [get]
func (a *Api) FetchCuponPorOrganizador(c echo.Context) error {
	organizadorParam := c.Param("organizadorId")
	organizadorId, err := strconv.ParseInt(organizadorParam, 10, 64)

	if err != nil || organizadorId <= 0 {
		return errors.HandleError(errors.BadRequestError.InvalidUpdatedByValue, c)
	}

	response, newErr := a.BllController.Cupon.FetchCuponPorOrganizador(organizadorId)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}
