package controller

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type UsuarioController struct {
	Logger         logging.Logger
	DB             *repository.NexiventPsqlEntidades
	GoogleClientID string
}

type GoogleUser struct {
	Email 	      string
	Name 	      string
	Picture       string
	Sub           string
	VerifiedEmail bool
}

func (uc *UsuarioController) RegisterUsuario(usuario *model.Usuario) (model.Usuario, *errors.Error) {
	err := uc.DB.Usuario.PostgresqlDB.Transaction(func(tx *gorm.DB) error {
		txRepo := uc.DB

		// Verificar si ya existe un usuario con el mismo correo
		existingUser, err := txRepo.Usuario.ObtenerUsuarioPorCorreo(usuario.Correo)
		switch {
		case err == gorm.ErrRecordNotFound:
			// No existe un usuario con este correo, continuar
		case err != nil:
			uc.Logger.Error(fmt.Sprintf("Error al verificar correo existente: %v", err))
			return err
		case existingUser != nil:
			uc.Logger.Warn(fmt.Sprintf("Intento de registro con correo existente: %s", usuario.Correo))
			return fmt.Errorf("el correo existe")
		}

		// Verificar si ya existe un usuario con el mismo número de documento
        existingDoc, err := txRepo.Usuario.ObtenerUsuarioPorNumDocumento(usuario.NumDocumento)
        switch {
        case err == gorm.ErrRecordNotFound:
            // No existe un usuario con este documento, continuar
        case err != nil:
            uc.Logger.Error(fmt.Sprintf("Error al verificar documento existente: %v", err))
            return err
        case existingDoc != nil:
            uc.Logger.Warn(fmt.Sprintf("Intento de registro con documento existente: %s", usuario.NumDocumento))
            return fmt.Errorf("el número de documento ya está registrado")
        }

		// Crear el usuario
		uc.Logger.Info(fmt.Sprintf("Creando usuario: %s", usuario.Correo))
		err = txRepo.Usuario.CrearUsuario(usuario)
		if err != nil {
			uc.Logger.Error(fmt.Sprintf("Error al crear usuario: %v", err))
			return err
		}

		// Asignar el rol por defecto "usuario"
		defaultRole, err := txRepo.Roles.ObtenerRolPorNombre("ASISTENTE")
		if err != nil {
			return err
		}
		if defaultRole == nil {
			return fmt.Errorf("rol por defecto 'asistente' no encontrado")
		}
		if usuario.TipoDocumento == "RUC_PERSONA" || usuario.TipoDocumento == "RUC_EMPRESA" {
			defaultRole, err = txRepo.Roles.ObtenerRolPorNombre("ORGANIZADOR")
			if err != nil {
				return err
			}
			usuario.EstadoDeCuenta = 0	
		}
		rolAsignado, err := txRepo.RolesUsuario.AsignarRolAUsuario(usuario.ID, defaultRole.ID, usuario.ID)
		if err != nil {
			return err
		}
		usuario.RolesAsignados = []model.RolUsuario{*rolAsignado}

		return nil
	})
	
	if err != nil {
		uc.Logger.Error(fmt.Sprintf("Error en transacción de registro: %v", err))
        // Retornar errores más específicos
        if err.Error() == "el correo ya está registrado" {
            return model.Usuario{}, &errors.Error{
                Code:    "DUPLICATE_EMAIL",
                Message: "El correo electrónico ya está registrado",
            }
        }
        if err.Error() == "el número de documento ya está registrado" {
            return model.Usuario{}, &errors.Error{
                Code:    "DUPLICATE_DOCUMENT",
                Message: "El número de documento ya está registrado",
            }
        }
        if err.Error() == "rol por defecto 'ASISTENTE' no encontrado" {
            return model.Usuario{}, &errors.InternalServerError.Default
        }
        
        return model.Usuario{}, &errors.InternalServerError.Default
	}
	return *usuario, nil
}

func (uc *UsuarioController) GetUsuario(id int64) (*model.Usuario, *errors.Error) {
	var usuario *model.Usuario
	var transactionErr error

	err := uc.DB.Usuario.PostgresqlDB.Transaction(func(tx *gorm.DB) error {
		txRepo := uc.DB
		var err error
		usuario, err = txRepo.Usuario.ObtenerUsuarioBasicoPorID(id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				transactionErr = fmt.Errorf("user not found")
				return transactionErr
			}
			transactionErr = err
			return err
		}

		// Validar que usuario no sea nil
        if usuario == nil {
            uc.Logger.Error(fmt.Sprintf("ObtenerUsuarioBasicoPorID retornó nil para ID %d", id))
            transactionErr = fmt.Errorf("user is nil")
            return transactionErr
        }

		comentarios := txRepo.Comentario.PostgresqlDB.Model(&model.Comentario{}).Where("usuario_id = ?", id).Find(&usuario.Comentarios)
		if comentarios.Error != nil {
			return comentarios.Error
		}

		ordenes := txRepo.OrdenDeCompra.PostgresqlDB.Model(&model.OrdenDeCompra{}).Where("usuario_id = ?", id).Find(&usuario.Ordenes)
		if ordenes.Error != nil {
			return ordenes.Error
		}

		// Falta cupones
		roles, err := txRepo.RolesUsuario.ListarRolesDeUsuario(usuario.ID)
		if err != nil {
			return err
		}
		usuario.RolesAsignados = roles

		return nil
	})
	if err != nil {
		uc.Logger.Error(fmt.Sprintf("Error en GetUsuario para ID %d: %v", id, err))
        
        if transactionErr != nil && transactionErr.Error() == "user not found" {
            return nil, &errors.ObjectNotFoundError.UserNotFound
        }
        if transactionErr != nil && transactionErr.Error() == "user is nil" {
            return nil, &errors.InternalServerError.Default
        }
        return nil, &errors.InternalServerError.Default
	}

	// Validación final por seguridad
    if usuario == nil {
        uc.Logger.Error(fmt.Sprintf("Usuario sigue siendo nil después de la transacción para ID %d", id))
        return nil, &errors.InternalServerError.Default
    }

    uc.Logger.Infof("Usuario obtenido exitosamente: ID=%d, Nombre=%s, Correo=%s", usuario.ID, usuario.Nombre, usuario.Correo)
	
	return usuario, nil
}

func (uc *UsuarioController) GetUsuarioConRoles(id int64) ([]*model.Usuario, *errors.Error) {

	usuarios, err := uc.DB.Usuario.ObtenerUsuariosPorRolID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.UserNotFound
		}
		return nil, &errors.InternalServerError.Default
	}
	return usuarios, nil
}

func (uc *UsuarioController) AuthenticateUsuario(correo, contrasenha string) (*model.Usuario, *errors.Error) {
	usuario, err := uc.DB.Usuario.ObtenerUsuarioPorCorreo(correo)
	if usuario.TipoDocumento == "RUC_PERSONA" || usuario.TipoDocumento == "RUC_EMPRESA" {
		return nil, &errors.AuthenticationError.InvalidCredentials
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.UserNotFound
		}
		return nil, &errors.InternalServerError.Default
	}

	ok, err := model.VerifyPassword(contrasenha, usuario.Contrasenha)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}

	if !ok {
		return nil, &errors.AuthenticationError.InvalidCredentials
	}

	return usuario, nil
}	

func (uc *UsuarioController) AuthenticateOrganizador(ruc, contrasenha string) (*model.Usuario, *errors.Error) {
	usuario, err := uc.DB.Usuario.ObtenerUsuarioPorNumDocumento(ruc)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.UserNotFound
		}
		return nil, &errors.InternalServerError.Default
	}
	if usuario.TipoDocumento != "RUC_PERSONA" && usuario.TipoDocumento != "RUC_EMPRESA" {
        uc.Logger.Warnf("Intento de login organizador con documento incorrecto: %s", usuario.TipoDocumento)
        return nil, &errors.AuthenticationError.InvalidCredentials
    }
	ok, err := model.VerifyPassword(contrasenha, usuario.Contrasenha)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}

	if !ok {
		return nil, &errors.AuthenticationError.InvalidCredentials
	}

	return usuario, nil
}

func (u *UsuarioController) VerifyGoogleToken(idToken string) (*GoogleUser, error) {
    ctx := context.Background()
	// Validar el token usando el Client ID de Google
	payload, err := idtoken.Validate(ctx, idToken, u.GoogleClientID)
    if err != nil {
        u.Logger.Errorf("Error validating Google token: %v", err)
        return nil, err
    }

    // Extraer información del payload
    email, _ := payload.Claims["email"].(string)
    name, _ := payload.Claims["name"].(string)
    picture, _ := payload.Claims["picture"].(string)
    sub, _ := payload.Claims["sub"].(string)
    emailVerified, _ := payload.Claims["email_verified"].(bool)

    if !emailVerified {
        return nil, fmt.Errorf("email not verified")
    }

    return &GoogleUser{
        Email:   email,
        Name:    name,
        Picture: picture,
        Sub:     sub,
        VerifiedEmail: emailVerified,
    }, nil
}

func (u *UsuarioController) GetOrCreateGoogleUser(googleUser *GoogleUser) (*model.Usuario, *errors.Error) {
    // Buscar usuario por correo
    usuario, err := u.DB.Usuario.ObtenerUsuarioPorCorreo(googleUser.Email)
    
    if err == nil && usuario != nil {
        // Usuario existe, actualizar último acceso
        u.Logger.Infof("Usuario existente de Google: %s", googleUser.Email)
        return usuario, nil
    }

    // Crear nuevo usuario con Google
    u.Logger.Infof("Creando nuevo usuario de Google: %s", googleUser.Email)
    
    nuevoUsuario := &model.Usuario{
        Nombre:          googleUser.Name,
        Correo:          googleUser.Email,
        TipoDocumento:   "GOOGLE",
        NumDocumento:    googleUser.Sub, // Usar Sub de Google como documento único
        EstadoDeCuenta:  1, // Verificado automáticamente
        Estado:          1, // Activo
        Contrasenha:     "", // Sin contraseña para usuarios de Google
    }

	// Registrar usuario
	usuarioCreado, newErr := u.RegisterUsuario(nuevoUsuario)
	if newErr != nil {
		return nil, newErr
	}

	return &usuarioCreado, nil
}

// ActivarUsuario activa un usuario (estado = 1)
func (uc *UsuarioController) ActivarUsuario(usuarioID int64, updatedBy int64) *errors.Error {
	// Verificar que el usuario que realiza la modificación existe
	_, err := uc.DB.Usuario.ObtenerUsuarioBasicoPorID(updatedBy)
	if err != nil {
		uc.Logger.Errorf("Usuario modificador no encontrado: %v", err)
		return &errors.Error{
			Code:    "INVALID_MODIFIER",
			Message: "Usuario modificador no encontrado",
		}
	}

	// Verificar que el usuario a modificar existe
	_, err = uc.DB.Usuario.ObtenerUsuarioBasicoPorID(usuarioID)
	if err != nil {
		uc.Logger.Errorf("Usuario a modificar no encontrado: %v", err)
		return &errors.ObjectNotFoundError.UserNotFound
	}

	// Actualizar el estado a 1 (activo)
	estado := int16(1)
	_, err = uc.DB.Usuario.ActualizarUsuario(
		usuarioID,
		nil,     // nombre
		nil,     // tipoDocumento
		nil,     // numDocumento
		nil,     // correo
		nil,     // contrasenha
		nil,     // telefono
		nil,     // estadoDeCuenta
		nil,     // codigoVerificacion
		nil,     // fechaExpiracionCodigo
		&estado, // estado = 1 (activo)
		updatedBy,
	)

	if err != nil {
		uc.Logger.Errorf("Error activando usuario %d: %v", usuarioID, err)
		return &errors.InternalServerError.Default
	}

	uc.Logger.Infof("Usuario %d activado exitosamente por usuario %d", usuarioID, updatedBy)
	return nil
}

// DesactivarUsuario desactiva un usuario (estado = 0)
func (uc *UsuarioController) DesactivarUsuario(usuarioID int64, updatedBy int64) *errors.Error {
	// Verificar que el usuario que realiza la modificación existe
	_, err := uc.DB.Usuario.ObtenerUsuarioBasicoPorID(updatedBy)
	if err != nil {
		uc.Logger.Errorf("Usuario modificador no encontrado: %v", err)
		return &errors.Error{
			Code:    "INVALID_MODIFIER",
			Message: "Usuario modificador no encontrado",
		}
	}

	// Usar la función existente DesactivarUsuario
	err = uc.DB.Usuario.DesactivarUsuario(usuarioID, updatedBy)
	if err != nil {
		uc.Logger.Errorf("Error desactivando usuario %d: %v", usuarioID, err)
		return &errors.ObjectNotFoundError.UserNotFound
	}

	uc.Logger.Infof("Usuario %d desactivado exitosamente por usuario %d", usuarioID, updatedBy)
	return nil
}