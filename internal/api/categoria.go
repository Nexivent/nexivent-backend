package api

import (
	"net/http"
	"strconv"

	"github.com/Loui27/nexivent-backend/errors"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// @Summary 			Get Evento.
// @Description 		Gets an event given its id.
// @Tags 				Evento
// @Accept 				json
// @Produce 			json
// @Param               eventoId    path   int  true  "Evento ID"
// @Success 			200 {object} schemas.EventoResponse "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/evento/{eventoId}/ [get]

func (a *Api) GetCategoria(c echo.Context) error {
	categoriaIdStr := c.Param("categoriaId")
	categoriaId, parseErr := strconv.ParseInt(categoriaIdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.Categoria.GetCategoriaById(categoriaId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}


// @Summary 			Fetch Eventos.
// @Description 		Fetches all available events.
// @Tags 				Evento
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} schemas.EventosPaginados "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/evento/ [get]
func (a *Api) FetchCategorias(c echo.Context) error {
	response, err := a.BllController.Categoria.FetchCategorias()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Evento.
// @Description 		Create a new event with all related entities.
// @Tags 				Evento
// @Accept 				json
// @Produce 			json
// @Param               request body schemas.EventoRequest true "Create Evento Request"
// @Success 			201 {object} schemas.EventoResponse "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/evento/ [post]
func (a *Api) CreateCategoria(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `usuarioCreacion` param)
	//usuarioCreacion := int64(1) // Hardcoded for now

	var request schemas.CategoriaRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Categoria.CreateCategoria(request)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}
