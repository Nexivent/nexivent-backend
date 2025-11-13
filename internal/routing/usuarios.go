package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
	"gorm.io/gorm"
)

func getUsuarios(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var input struct {
		Rol string
		util.Filters
	}

	input.Rol = internal.ReadString(r.URL.Query(), "rol", "")
	input.Page = uint64(internal.ReadInt(r.URL.Query(), "page", 1, nil))
	input.PageSize = uint64(internal.ReadInt(r.URL.Query(), "page_size", 20, nil))
	input.Sort = internal.ReadString(r.URL.Query(), "sort", "usuario_id")

	// Validate the query parameters.
	v := validator.New()
	if util.ValidateFilters(v, input.Filters); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	var usuarios []*model.Usuario
	var rol *model.Rol
	if input.Rol != "" {
		err := app.Repository.DB.Transaction(func(tx *gorm.DB) error {
			txRepo := app.Repository.WithTx(tx)
			var err error
			rol, err = txRepo.Roles.ObtenerRolPorNombre(input.Rol)
			if err != nil {
				return err
			}

			usuarios, err = app.Repository.Usuarios.ObtenerUsuariosPorRolID(rol.ID)
			if err != nil {
				return err
			}

			for i := range usuarios {
				usuarios[i].RolesAsignados, err = txRepo.RolesUsuarios.ListarRolesDeUsuario(usuarios[i].ID)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}
	} else {
		var err error
		usuarios, err = app.Repository.Usuarios.ObtenerUsuarios()
		if err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		for i := range usuarios {
			usuarios[i].RolesAsignados, err = app.Repository.RolesUsuarios.ListarRolesDeUsuario(usuarios[i].ID)
			if err != nil {
				app.ServerErrorResponse(w, r, err)
				return
			}
		}
	}

	err := internal.WriteJSON(w, http.StatusOK, internal.Envelope{"usuarios": usuarios}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func getUsuarioPorCorreo(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Leer el correo de los query parameters
	qs := r.URL.Query()
	correo := internal.ReadString(qs, "correo", "")

	if correo == "" {
		app.BadRequestResponse(w, r, nil)
		return
	}

	usuario, err := app.Repository.Usuarios.ObtenerUsuarioPorCorreo(correo)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"usuario": usuario}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func getUsuarioPorID(w http.ResponseWriter, r *http.Request) {
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

	usuario, err := app.Repository.Usuarios.ObtenerUsuarioBasicoPorID(id)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = internal.WriteJSON(w, http.StatusOK, internal.Envelope{"usuario": usuario}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
