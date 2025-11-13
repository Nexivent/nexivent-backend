package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
)

func getTicketsPorOrden(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del path parameter
	ordenID, err := internal.ReadIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	tickets, err := app.Repository.Tickets.ObtenerTicketsPorOrden(int64(ordenID))
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"tickets": tickets}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
