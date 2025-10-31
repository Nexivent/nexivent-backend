package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	appcontext "github.com/Nexivent/nexivent-backend/internal/context"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la aplicación del contexto
	app := appcontext.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	env := internal.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.Config.Env,
			"version":     "1.0",
		},
	}

	err := internal.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.Logger.Error("failed to write JSON response", "error", err)
		http.Error(w, "El servidor encontró un problema y no pudo procesar su solicitud", http.StatusInternalServerError)
	}
}
