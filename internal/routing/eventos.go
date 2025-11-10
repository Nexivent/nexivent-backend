package routing

import (
	"fmt"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data"
	"github.com/Nexivent/nexivent-backend/internal/validator"
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

	app.Logger.Info("fetching event", "id", id)

	evento := data.Evento{
		ID: uint64(id),
	}

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
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", evento.ID))

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"evento": evento}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
