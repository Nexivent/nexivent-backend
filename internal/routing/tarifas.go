package routing

import (
	"net/http"
	"strconv"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

func getTarifasPorIDs(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	v := validator.New()
	qs := r.URL.Query()

	// Leer los IDs de las tarifas de la query string (formato: ?ids=1,2,3)
	idsStr := internal.ReadCSV(qs, "ids", []string{})

	if len(idsStr) == 0 {
		v.AddError("ids", "debe proporcionar al menos un ID de tarifa")
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// Convertir los strings a int64
	ids := make([]int64, 0, len(idsStr))
	for _, idStr := range idsStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			v.AddError("ids", "todos los IDs deben ser números enteros positivos")
			continue
		}
		ids = append(ids, id)
	}

	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	tarifas, err := app.Repository.Tarifas.ObtenerTarifasPorIDs(ids)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"tarifas": tarifas}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
