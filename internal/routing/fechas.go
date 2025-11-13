package routing

import (
	"net/http"
	"time"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

func getFechaPorFecha(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	v := validator.New()
	qs := r.URL.Query()

	// Leer el parámetro de fecha de la query string
	fecha := internal.ReadTime(qs, time.Time{}, v)

	// Validar que se proporcionó una fecha
	v.Check(!fecha.IsZero(), "fecha", "la fecha es obligatoria")

	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// Buscar la fecha en el repositorio
	fechaObtenida, err := app.Repository.Fechas.ObtenerPorFecha(fecha)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	// Verificar si se encontró la fecha
	if fechaObtenida == nil || fechaObtenida.ID == 0 {
		app.NotFoundResponse(w, r)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"fecha": fechaObtenida}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
