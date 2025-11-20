package controller

import (
	"fmt"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type UsuarioController struct {
	Logger logging.Logger
	DB     *repository.NexiventPsqlEntidades
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
	err := uc.DB.Usuario.PostgresqlDB.Transaction(func(tx *gorm.DB) error {
		txRepo := uc.DB
		usuario, err := txRepo.Usuario.ObtenerUsuarioBasicoPorID(id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("user not found")
			}
			return err
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
		if err.Error() == "user not found" {
			return nil, &errors.ObjectNotFoundError.UserNotFound
		}
		return nil, &errors.InternalServerError.Default
	}
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