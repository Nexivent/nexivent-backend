package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	appcontext "github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data"
)

func getEvent(w http.ResponseWriter, r *http.Request) {
	// Obtener la aplicación del contexto
	app := appcontext.GetApplication(r.Context())
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
