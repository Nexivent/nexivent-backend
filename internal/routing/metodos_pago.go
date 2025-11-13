package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
)

func getMetodosPago(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	metodosPago, err := app.Repository.MetodosPago.ListarMetodosActivos()
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"metodosPago": metodosPago}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func getMetodoPagoPorID(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del path parameter
	id, err := internal.ReadIDParam(r)
	if err != nil {
		app.NotFoundResponse(w, r)
		return
	}

	metodoPago, err := app.Repository.MetodosPago.ObtenerMetodoDePagoBasicoPorID(int64(id))
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"metodoPago": metodoPago}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
