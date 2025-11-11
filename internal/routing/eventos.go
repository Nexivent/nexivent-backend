package routing

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data"
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
	"github.com/google/uuid"
)

func getEvento(w http.ResponseWriter, r *http.Request) {
	// Obtener la aplicación del contexto
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := internal.ReadIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	evento, err := app.Models.Eventos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.ServerErrorResponse(w, r, err)
		}
		return
	}
	// Evento encontrado, devolver como respuesta JSON

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"evento": evento}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func postEvento(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var evento data.Evento
	err := internal.ReadJSON(w, r, &evento)
	if err != nil {
		app.BadRequestResponse(w, r, err)
	}

	v := validator.New()
	if data.ValidateEvento(v, &evento); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Eventos.Insert(&evento)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/eventos/%s", evento.ID))

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"evento": evento}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}

func updateEvento(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := internal.ReadIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	evento, err := app.Models.Eventos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.ServerErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		CategoriaID      *uuid.UUID         `json:"categoria_id"`
		Titulo           *string            `json:"titulo"`
		Descripcion      *string            `json:"descripcion"`
		Lugar            *string            `json:"lugar"`
		EventoEstado     *util.EstadoEvento `json:"evento_estado"`
		CantMeGusta      *int               `json:"cant_me_gusta"`
		CantNoInteresa   *int               `json:"cant_no_interesa"`
		CantVendidoTotal *int               `json:"cant_vendido_total"`
		TotalRecaudado   *float64           `json:"total_recaudado"`
		Estado           *util.Estado       `json:"estado"`
	}

	if input.CategoriaID != nil {
		evento.CategoriaID = *input.CategoriaID
	}

	if input.Titulo != nil {
		evento.Titulo = *input.Titulo
	}

	if input.Descripcion != nil {
		evento.Descripcion = *input.Descripcion
	}

	if input.Lugar != nil {
		evento.Lugar = *input.Lugar
	}

	if input.EventoEstado != nil {
		evento.EventoEstado = *input.EventoEstado
	}
	if input.CantMeGusta != nil {
		evento.CantMeGusta = *input.CantMeGusta
	}

	if input.CantNoInteresa != nil {
		evento.CantNoInteresa = *input.CantNoInteresa
	}
	
	if input.CantVendidoTotal != nil {
		evento.CantVendidoTotal = *input.CantVendidoTotal
	}
	
	if input.TotalRecaudado != nil {
		evento.TotalRecaudado = *input.TotalRecaudado
	}
	
	if input.Estado != nil {
		evento.Estado = *input.Estado
	}	

	err = internal.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateEvento(v, evento); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Eventos.Update(evento)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"evento": evento}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
