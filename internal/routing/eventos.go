package routing

import (
	"fmt"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data"
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

	// Ahora puedes usar app.Logger, app.Config, etc.
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

	fmt.Fprintf(w, "%+v\n", evento)
}
