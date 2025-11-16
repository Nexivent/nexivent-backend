package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
	"github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	datautil "github.com/Nexivent/nexivent-backend/internal/data/model/util"
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
	input.SortSafeList = []string{"usuario_id", "-usuario_id", "nombre", "-nombre", "correo", "-correo", "fecha_creacion", "-fecha_creacion"}

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
			// rol, err = txRepo.Roles.ObtenerRolPorNombre(input.Rol)
			// if err != nil {
			// 	return err
			// }

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

		// for i := range usuarios {
		// 	usuarios[i].RolesAsignados, err = app.Repository.RolesUsuarios.ListarRolesDeUsuario(usuarios[i].ID)
		// 	if err != nil {
		// 		app.ServerErrorResponse(w, r, err)
		// 		return
		// 	}
		// }
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

func registerUsuario(w http.ResponseWriter, r *http.Request) {
	app := context.GetApplication(r.Context())
	if app == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var input struct {
		Nombre        string                 `json:"nombre"`
		TipoDocumento datautil.TipoDocumento `json:"tipoDocumento"`
		NumDocumento  string                 `json:"numDocumento"`
		Correo        string                 `json:"correo"`
		Password      string                 `json:"password"`
		Telefono      *string                `json:"telefono,omitempty"`
	}

	err := internal.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}
	usuario := &model.Usuario{
		Nombre:         input.Nombre,
		TipoDocumento:  input.TipoDocumento,
		NumDocumento:   input.NumDocumento,
		Correo:         input.Correo,
		Telefono:       input.Telefono,
		EstadoDeCuenta: datautil.Inactivo,
	}

	v := validator.New()
	model.ValidateUsuario(v, usuario, &input.Password)
	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	password, err := model.HashPassword(input.Password)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
	usuario.Password = []byte(password)

	// Crear el usuario en la base de datos
	err = app.Repository.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := app.Repository.WithTx(tx)

		// Verificar si ya existe un usuario con el mismo correo
		existingUser, err := txRepo.Usuarios.ObtenerUsuarioPorCorreo(input.Correo)
		switch {
			case err == gorm.ErrRecordNotFound:
				// No existe un usuario con este correo, continuar
			case err != nil:
				return err
			case existingUser != nil:
				v.AddError("correo", "ya existe un usuario con este correo")
				return nil
		}

		// Crear el usuario
		err = txRepo.Usuarios.CrearUsuario(usuario)
		if err != nil {
			return err
		}

		// // Asignar el rol por defecto "usuario"
		// defaultRole, err := txRepo.Roles.ObtenerRolPorNombre("usuario")
		// if err != nil {
		// 	return err
		// }
		// if defaultRole == nil {
		// 	return fmt.Errorf("rol por defecto 'usuario' no encontrado")
		// }

		// rolAsignado, err := txRepo.RolesUsuarios.AsignarRolAUsuario(usuario.ID, defaultRole.ID, usuario.ID)
		// if err != nil {
		// 	return err
		// }
		// usuario.RolesAsignados = []model.RolUsuario{*rolAsignado}

		return nil
	})
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = internal.WriteJSON(w, http.StatusCreated, internal.Envelope{"usuario": usuario}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
