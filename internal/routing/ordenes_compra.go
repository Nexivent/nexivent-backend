package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
)

func getOrdenDeCompraPorID(w http.ResponseWriter, r *http.Request) {
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

	orden, err := app.Repository.OrdenesCompra.ObtenerOrdenBasica(int64(id))
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"orden": orden}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
