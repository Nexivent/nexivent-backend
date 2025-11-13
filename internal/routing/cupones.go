package routing

import (
	"errors"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/data/repository"
	datautil "github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

func getCupones(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cupones, err := app.Repository.Cupones.ObtenerCupones()
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"cupones": cupones}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func postCupon(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var cupon model.Cupon
	err := internal.ReadJSON(w, r, &cupon)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	err = app.Repository.Cupones.CrearCupon(&cupon)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"cupon": cupon}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func putCupon(w http.ResponseWriter, r *http.Request) {
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

	result := app.Repository.DB.Where("cupon_id = ?", id).First(&model.Cupon{})
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, repository.ErrRecordNotFound):
			app.NotFoundResponse(w, r)
		default:
			app.ServerErrorResponse(w, r, result.Error)
		}
		return
	}

	var input struct {
		ID          uint64          `json:"id"`
		Codigo      string          `json:"codigo"`
		Descripcion string          `json:"descripcion"`
		Tipo        string          `json:"tipo"`
		Valor       float64         `json:"valor"`
		EstadoCupon datautil.Estado `json:"estadoCupon"`
		EventoID    uint64          `json:"eventoId"`

		Usuarios []model.UsuarioCupon `json:"usuarios"`
	}

	err = internal.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	cupon := &model.Cupon{
		Descripcion: input.Descripcion,
		Tipo:        input.Tipo,
		Valor:       input.Valor,
		EstadoCupon: input.EstadoCupon,
	}

	v := validator.New()
	if model.ValidateCupon(v, cupon); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}
}
