package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
)

func getRoles(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	roles, err := app.Repository.Roles.ObtenerRoles()
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"roles": roles}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func getRolPorNombre(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Leer el nombre de los query parameters
	qs := r.URL.Query()
	nombre := internal.ReadString(qs, "nombre", "")

	if nombre == "" {
		app.BadRequestResponse(w, r, nil)
		return
	}

	rol, err := app.Repository.Roles.ObtenerRolPorNombre(nombre)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"rol": rol}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
