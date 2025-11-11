package routing

import (
	"fmt"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
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

	err = app.Repository.Eventos.CrearEvento(&evento)
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
	input.Titulo = internal.ReadString(qs, "title", "")

	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = uint64(internal.ReadInt(qs, "page", 1, v))
	input.Filters.PageSize = uint64(internal.ReadInt(qs, "page_size", 20, v))

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = internal.ReadString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "year", "-id", "-title", "-year"}

	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}
