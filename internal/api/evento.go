package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/utils/convert"
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

// @Summary      Feed de eventos del usuario
// @Description  Obtiene eventos activos para el feed, excluyendo los que ya tienen interacción del usuario
// @Tags         Evento
// @Accept       json
// @Produce      json
// @Param        usuarioId query int64 false "ID del usuario (opcional)"
// @Success      200 {object} schemas.EventosPaginados "OK"
// @Failure      400 {object} errors.Error "Bad Request"
// @Failure      422 {object} errors.Error "Unprocessable Entity"
// @Failure      500 {object} errors.Error "Internal Server Error"
// @Router       /feed/eventos [get]
func (a *Api) FetchEventosFeed(c echo.Context) error {
	// 1. Leer usuarioId del query param (ej: /feed/eventos?usuarioId=123)
	uidStr := c.QueryParam("usuarioId")

	var usuarioId *int64 = nil

	// 2. Si se envía, convertirlo a int64 y validar
	if uidStr != "" {
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil || uid <= 0 {
			// Si falla el parse o es <= 0 → error 422
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		usuarioId = &uid
	}

	// 3. Llamar a la lógica de negocio (tu función real)
	resp, newErr := a.BllController.Evento.FetchEventosFeed(usuarioId)
	if newErr != nil {
		// 4. Si la capa BLL devuelve error → responderlo
		return errors.HandleError(*newErr, c)
	}

	// 5. Todo OK → devolver JSON 200
	return c.JSON(http.StatusOK, resp)
}

func (a *Api) FetchEventosConInteraccionesFeed(c echo.Context) error {
	// 1. Leer usuarioId del query param (ej: /feed/eventos?usuarioId=123)
	uidStr := c.QueryParam("usuarioId")

	var usuarioId int64 = 0

	// 2. Si se envía, convertirlo a int64 y validar
	if uidStr != "" {
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil || uid <= 0 {
			// Si falla el parse o es <= 0 → error 422
			return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
		}
		usuarioId = uid
	}

	// 3. Llamar a la lógica de negocio (tu función real)
	resp, newErr := a.BllController.Evento.FetchEventosConInteraccionesFeed(&usuarioId)
	if newErr != nil {
		// 4. Si la capa BLL devuelve error → responderlo
		return errors.HandleError(*newErr, c)
	}

	// 5. Todo OK → devolver JSON 200
	return c.JSON(http.StatusOK, resp)
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
// @Param        estado        query   string  false  "Estado del evento (BORRADOR|PUBLICADO|CANCELADO)"
// @Param        soloFuturos   query   bool    false  "Si es true, solo eventos con fecha desde hoy"
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

	var estado *int16
	if eStr := c.QueryParam("estado"); eStr != "" {
		upper := strings.ToUpper(strings.TrimSpace(eStr))
		switch upper {
		case "BORRADOR", "PUBLICADO", "CANCELADO":
			val := convert.MapEstadoToInt16(upper)
			estado = &val
		default:
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
	}

	soloFuturos := false
	if sf := c.QueryParam("soloFuturos"); sf != "" {
		switch strings.ToLower(sf) {
		case "true", "1", "yes", "y":
			soloFuturos = true
		case "false", "0", "no", "n":
			soloFuturos = false
		default:
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
	}

	response, err := a.BllController.Evento.FetchEventosWithFilters(
		categoriaID,
		organizadorID,
		titulo,
		descripcion,
		lugar,
		fecha,
		horaInicio,
		estado,
		soloFuturos,
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

// @Summary             Obtener Reporte de Evento
// @Description         Genera el reporte detallado de un evento aplicando filtros opcionales.
// @Tags                Evento
// @Accept              json
// @Produce             json
// @Param               eventoID     query int64   false "ID del Evento"
// @Param               fechaDesde   query string  false "Fecha desde (YYYY-MM-DD)"
// @Param               fechaHasta   query string  false "Fecha hasta (YYYY-MM-DD)"
// @Success             200 {array}  schemas.EventoReporte "OK"
// @Failure             400 {object} errors.Error "Bad Request"
// @Failure             404 {object} errors.Error "Not Found"
// @Failure             422 {object} errors.Error "Unprocessable Entity"
// @Failure             500 {object} errors.Error "Internal Server Error"
// @Router              /evento/reporte [get]
func (a *Api) GetReporteEvento(c echo.Context) error {
	// TODO: Validar access token

	usuarioIDStr := c.Param("organizadorId")
	usuarioID, parseErr := strconv.ParseInt(usuarioIDStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	// --- Obtener query params ---
	var (
		eventoID   *int64
		fechaDesde *time.Time
		fechaHasta *time.Time
	)

	// eventoID
	if q := c.QueryParam("eventoId"); q != "" {
		val, err := strconv.ParseInt(q, 10, 64)
		if err != nil {
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
		eventoID = &val
	}

	// fechaDesde
	if q := c.QueryParam("fechaDesde"); q != "" {
		fd, err := time.Parse("2006-01-02", q)
		if err != nil {
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
		fechaDesde = &fd
	}

	// fechaHasta
	if q := c.QueryParam("fechaHasta"); q != "" {
		fh, err := time.Parse("2006-01-02", q)
		if err != nil {
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
		fechaHasta = &fh
	}

	// --- Llamar al controller ---
	response, newErr := a.BllController.Evento.GetReporteEvento(
		usuarioID,
		eventoID,
		fechaDesde,
		fechaHasta,
	)

	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary             Reporte por Organizador
// @Description         Devuelve el resumen de todos los eventos de un organizador.
// @Tags                Evento
// @Accept              json
// @Produce             json
// @Param               organizadorId path int   true  "ID del Organizador"
// @Param               fechaDesde    query string false "Fecha desde (YYYY-MM-DD)"
// @Param               fechaHasta    query string false "Fecha hasta (YYYY-MM-DD)"
// @Success             200 {array}  schemas.EventoOrganizadorReporte "OK"
// @Failure             400 {object} errors.Error "Bad Request"
// @Failure             404 {object} errors.Error "Not Found"
// @Failure             422 {object} errors.Error "Unprocessable Entity"
// @Failure             500 {object} errors.Error "Internal Server Error"
// @Router              /organizador/{organizadorId}/eventos/reporte [get]
func (a *Api) GetReporteEventosOrganizador(c echo.Context) error { ////////////////////////////////////////////////
	organizadorStr := c.Param("organizadorId")
	organizadorID, parseErr := strconv.ParseInt(organizadorStr, 10, 64)
	if parseErr != nil || organizadorID <= 0 {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	var fechaDesde *time.Time
	if q := c.QueryParam("fechaDesde"); q != "" {
		fd, err := time.Parse("2006-01-02", q)
		if err != nil {
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
		fechaDesde = &fd
	}

	var fechaHasta *time.Time
	if q := c.QueryParam("fechaHasta"); q != "" {
		fh, err := time.Parse("2006-01-02", q)
		if err != nil {
			return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
		}
		fechaHasta = &fh
	}

	resp, newErr := a.BllController.Evento.GetReporteEventosOrganizador(organizadorID, fechaDesde, fechaHasta)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Api) GetEventoSummary(c echo.Context) error {
	eventoIDStr := c.Param("id")
	eventoID, parseErr := strconv.ParseInt(eventoIDStr, 10, 64)
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	response, err := a.BllController.Evento.GetEventoDetalle(eventoID)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// EditarEventoFull reemplaza completamente un evento (solo borrador sin ventas).
func (a *Api) EditarEventoFull(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
	}

	var req schemas.EditarEventoFullRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	resp, errBll := a.BllController.Evento.EditarEventoFull(id, req)
	if errBll != nil {
		return errors.HandleError(*errBll, c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Api) EditarEvento(c echo.Context) error {
	// 1) Tomar el ID desde el path :id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
	}

	// 2) Parsear body
	var req schemas.EditarEventoRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// 3) Forzar que el ID venga del path (por seguridad)
	req.IdEvento = id

	// (Opcional) si luego lees el usuario del token:
	// if userID, ok := c.Get("userID").(int64); ok {
	//     req.UsuarioModificacion = &userID
	// }

	// 4) Llamar al BO / controller
	resp, errBll := a.BllController.Evento.EditarEvento(&req)
	if errBll != nil {
		return errors.HandleError(*errBll, c)
	}

	// 5) Devolver el detalle actualizado del evento
	return c.JSON(http.StatusOK, resp)
}

// @Summary             Crear interacción usuario–evento
// @Description         Registra una interacción de un usuario con un evento (like, vista)
// @Tags                Evento
// @Accept              json
// @Produce             json
// @Param               request body schemas.InteraccionConEventoRequest true "Datos de interacción con el evento"
// @Success             201 {object} schemas.InteraccionConEventoResponse "Interacción creada"
// @Failure             400 {object} errors.Error "Parámetros inválidos"
// @Failure             404 {object} errors.Error "Evento o Usuario no encontrado"
// @Failure             422 {object} errors.Error "Error en el request"
// @Failure             500 {object} errors.Error "Internal Server Error"
// @Router              /evento/interaccion [post]
func (a *Api) PostInteraccionUsuarioEvento(c echo.Context) error {
	var req schemas.InteraccionConEventoRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.BadRequestError.InvalidBodyFormat, c)
	}

	resp, newErr := a.BllController.Evento.PostInteraccionUsuarioEvento(req)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, resp)
}

// @Summary             Actualizar interacción usuario–evento
// @Description         Modifica una interacción existente de un usuario con un evento.
// @Tags                Evento
// @Accept              json
// @Produce             json
// @Param               request body schemas.InteraccionConEventoRequest true "Datos actualizados de la interacción"
// @Success             200 {object} schemas.InteraccionConEventoResponse "Interacción actualizada"
// @Failure             400 {object} errors.Error "Parámetros inválidos"
// @Failure             404 {object} errors.Error "Evento, Usuario o interacción no encontrado"
// @Failure             422 {object} errors.Error "Error en el request"
// @Failure             500 {object} errors.Error "Internal Server Error"
// @Router              /evento/interaccion [put]
func (a *Api) PutInteraccionUsuarioEvento(c echo.Context) error {
	var req schemas.InteraccionConEventoRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.BadRequestError.InvalidBodyFormat, c)
	}

	resp, newErr := a.BllController.Evento.PutInteraccionUsuarioEvento(req)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Api) GetAsistentesPorEvento(c echo.Context) error {
	eventoIDStr := c.Param("eventoId")
	eventoID, parseErr := strconv.ParseInt(eventoIDStr, 10, 64)
	if parseErr != nil {
		a.Logger.Errorf("❌ [API] Error parseando eventoId: %v", parseErr)
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	asistentes, err := a.BllController.Evento.GetAsistentesPorEvento(eventoID)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	response := map[string]interface{}{
		"success":    true,
		"asistentes": asistentes,
		"total":      len(asistentes),
	}

	return c.JSON(http.StatusOK, response)
}
