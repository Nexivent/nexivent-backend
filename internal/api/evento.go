package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
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
func (a *Api) GetEvento(c echo.Context) error {
	eventoIdStr := c.Param("eventoId")
	eventoId, parseErr := strconv.ParseInt(eventoIdStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.Evento.GetEventoById(eventoId)
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
func (a *Api) FetchEventos(c echo.Context) error {
	response, err := a.BllController.Evento.FetchEventos()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary      Fetch Eventos filtrados.
// @Description  Obtiene la lista de eventos disponibles aplicando filtros opcionales.
// @Tags         Evento
// @Accept       json
// @Produce      json
// @Param        categoriaId   query   int     false  "ID de categoría"
// @Param        organizadorId query   int     false  "ID de organizador"
// @Param        titulo        query   string  false  "Título del evento (coincidencia parcial)"
// @Param        descripcion   query   string  false  "Descripción (coincidencia parcial)"
// @Param        lugar         query   string  false  "Lugar del evento (coincidencia parcial)"
// @Param        fecha         query   string  false  "Fecha del evento (YYYY-MM-DD)"
// @Param        horaInicio    query   string  false  "Hora de inicio (HH:MM)"
// @Success      200  {object}  schemas.EventosPaginados  "OK"
// @Failure      400  {object}  errors.Error              "Bad Request"
// @Failure      404  {object}  errors.Error              "Not Found"
// @Failure      422  {object}  errors.Error              "Unprocessable Entity"
// @Failure      500  {object}  errors.Error              "Internal Server Error"
// @Router       /evento/filter [get]
func (a *Api) FetchEventosWithFilters(c echo.Context) error {
	fmt.Println("QUERY DEBUG:", c.QueryParams())

	var categoriaID *int64
	if catStr := c.QueryParam("categoriaId"); catStr != "" {
		parsed, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		categoriaID = &parsed
	}

	var organizadorID *int64
	if orgStr := c.QueryParam("organizadorId"); orgStr != "" {
		parsed, err := strconv.ParseInt(orgStr, 10, 64)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		organizadorID = &parsed
	}

	var titulo *string
	if t := c.QueryParam("titulo"); t != "" {
		titulo = &t
	}

	var descripcion *string
	if d := c.QueryParam("descripcion"); d != "" {
		descripcion = &d
	}

	var lugar *string
	if l := c.QueryParam("lugar"); l != "" {
		lugar = &l
	}

	var fecha *time.Time
	if fStr := c.QueryParam("fecha"); fStr != "" {
		parsed, err := time.Parse("2006-01-02", fStr)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidDateFormat, c)
		}
		fecha = &parsed
	}

	var horaInicio *time.Time
	if hStr := c.QueryParam("horaInicio"); hStr != "" {
		parsed, err := time.Parse("15:04", hStr)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidDateFormat, c)
		}
		horaInicio = &parsed
	}

	response, err := a.BllController.Evento.FetchEventosWithFilters(
		categoriaID,
		organizadorID,
		titulo,
		descripcion,
		lugar,
		fecha,
		horaInicio,
	)
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
func (a *Api) CreateEvento(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `usuarioCreacion` param)
	usuarioCreacion := int64(1) // Hardcoded for now///////////////////////

	// hashToken := c.Request().Header.Get("Authorization")
	// token, err := a.BllController.Token.ValidateToken(hashToken)
	// if err != nil {
	// 	return errors.HandleError(*err, c)
	// }

	// usuarioCreacion = token.UsuarioID

	var request schemas.EventoRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// if request.IdOrganizador != usuarioCreacion {
	// 	return errors.HandleError(errors.AuthenticationError.UnauthorizedUser, c)
	// }

	response, newErr := a.BllController.Evento.CreateEvento(request, usuarioCreacion)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}
