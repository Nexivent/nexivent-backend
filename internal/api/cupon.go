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
