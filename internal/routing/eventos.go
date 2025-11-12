package routing

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
	"gorm.io/gorm"
)

func getEvento(w http.ResponseWriter, r *http.Request) {
	// // Obtener la aplicación del contexto
	// app := context.GetApplication(r.Context())
	// if app == nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// id, err := internal.ReadIDParam(r)
	// if err != nil {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// evento, err := app.Repository.Eventos.(id)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, model.ErrRecordNotFound):
	// 		app.NotFoundResponse(w, r)
	// 	default:
	// 		app.ServerErrorResponse(w, r, err)
	// 	}
	// 	return
	// }
	// // Evento encontrado, devolver como respuesta JSON

	// err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"evento": evento}, nil)
	// if err != nil {
	// 	app.ServerErrorResponse(w, r, err)
	// }
}

func postEvento(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var evento model.Evento
	err := internal.ReadJSON(w, r, &evento)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if model.ValidateEvento(v, &evento); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// Usar una transacción para todas las operaciones
	err = app.Repository.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := app.Repository.WithTx(tx)

		for i := range evento.Fechas {
			fecha := model.Fecha{
				FechaEvento: evento.Fechas[i].HoraInicio.Truncate(24 * time.Hour),
			}
			if err := txRepo.Fechas.CrearFecha(&fecha); err != nil {
				return err
			}
			evento.Fechas[i].FechaID = fecha.ID
		}

		if err := txRepo.Eventos.CrearEvento(&evento); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/eventos/%d", evento.ID))

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"evento": evento}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func patchEvento(w http.ResponseWriter, r *http.Request) {
	// app := context.GetApplication(r.Context())
	// if app == nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// id, err := internal.ReadIDParam(r)
	// if err != nil {
	// 	app.NotFoundResponse(w, r)
	// 	return
	// }

	// evento, err := app.Models.Eventos.Get(id)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, model.ErrRecordNotFound):
	// 		app.NotFoundResponse(w, r)
	// 	default:
	// 		app.ServerErrorResponse(w, r, err)
	// 	}
	// 	return
	// }

	// var input struct {
	// 	CategoriaID      *uuid.UUID         `json:"categoria_id"`
	// 	Titulo           *string            `json:"titulo"`
	// 	Descripcion      *string            `json:"descripcion"`
	// 	Lugar            *string            `json:"lugar"`
	// 	EventoEstado     *util.EstadoEvento `json:"evento_estado"`
	// 	CantMeGusta      *int               `json:"cant_me_gusta"`
	// 	CantNoInteresa   *int               `json:"cant_no_interesa"`
	// 	CantVendidoTotal *int               `json:"cant_vendido_total"`
	// 	TotalRecaudado   *float64           `json:"total_recaudado"`
	// 	Estado           *util.Estado       `json:"estado"`
	// }

	// if input.CategoriaID != nil {
	// 	evento.CategoriaID = *input.CategoriaID
	// }

	// if input.Titulo != nil {
	// 	evento.Titulo = *input.Titulo
	// }

	// if input.Descripcion != nil {
	// 	evento.Descripcion = *input.Descripcion
	// }

	// if input.Lugar != nil {
	// 	evento.Lugar = *input.Lugar
	// }

	// if input.EventoEstado != nil {
	// 	evento.EventoEstado = *input.EventoEstado
	// }
	// if input.CantMeGusta != nil {
	// 	evento.CantMeGusta = *input.CantMeGusta
	// }

	// if input.CantNoInteresa != nil {
	// 	evento.CantNoInteresa = *input.CantNoInteresa
	// }

	// if input.CantVendidoTotal != nil {
	// 	evento.CantVendidoTotal = *input.CantVendidoTotal
	// }

	// if input.TotalRecaudado != nil {
	// 	evento.TotalRecaudado = *input.TotalRecaudado
	// }

	// if input.Estado != nil {
	// 	evento.Estado = *input.Estado
	// }

	// err = internal.ReadJSON(w, r, &input)
	// if err != nil {
	// 	app.BadRequestResponse(w, r, err)
	// 	return
	// }

	// v := validator.New()
	// if model.ValidateEvento(v, evento); !v.Valid() {
	// 	app.FailedValidationResponse(w, r, v.Errors)
	// 	return
	// }

	// err = app.Models.Eventos.Patch(evento)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, model.ErrEditConflict):
	// 		app.EditConflictResponse(w, r)
	// 	default:
	// 		app.ServerErrorResponse(w, r, err)
	// 	}
	// 	return
	// }

	// err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"evento": evento}, nil)
	// if err != nil {
	// 	app.ServerErrorResponse(w, r, err)
	// }
}

func getEventos(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var input struct {
		Titulo string
		util.Filters
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Titulo = internal.ReadString(qs, "titulo", "")

	input.Filters.CategoriaID = uint64(internal.ReadInt(qs, "categoria_id", 0, v))
	input.Filters.Lugar = internal.ReadString(qs, "lugar", "")
	input.Filters.Fecha = internal.ReadTime(qs, time.Time{}, v)
	input.Filters.Descripcion = internal.ReadString(qs, "descripcion", "")

	input.Filters.Page = uint64(internal.ReadInt(qs, "page", 1, v))
	input.Filters.PageSize = uint64(internal.ReadInt(qs, "page_size", 20, v))
	input.Filters.Sort = internal.ReadString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "time", "-id", "-title", "-time"}

	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if util.ValidateFilters(v, input.Filters); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	eventos, err := app.Repository.Eventos.ObtenerEventosDisponiblesConFiltros(
		&input.Filters.CategoriaID,
		&input.Titulo,
		&input.Descripcion,
		&input.Lugar,
		nil,
		nil,
	)

	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"eventos": eventos}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
