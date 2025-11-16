package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
)

func getCategorias(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categorias, err := app.Repository.Categorias.ObtenerCategorias()
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"categorias": categorias}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func postCategoria(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var input struct {
		Nombre        string `json:"nombre"`
		Descripcion   string `json:"descripcion"`
	}

	err := internal.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	categoria := &model.Categoria{
		Nombre:        input.Nombre,
		Descripcion:   input.Descripcion,
	}

	err = app.Repository.Categorias.CrearCategoria(categoria)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"categoria": categoria}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}