package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/labstack/echo/v4"
)

func (a *Api) RegisterUsuario(c echo.Context) error {

	var input struct {
		//para reenvio de codigo de verif
		Resend 	  bool    `json:"resend"`
		UsuarioID int64   `json:"usuario_id"`
		//campos de registro normal
		Nombre        string  `json:"nombre"`
		TipoDocumento string  `json:"tipo_documento"`
		NumDocumento  string  `json:"num_documento"`
		Correo        string  `json:"correo"`
		Email         string  `json:"email"` // Alias para correo
		Contrasenha   string  `json:"contrasenha"`
		Contrasena    string  `json:"contrasena"` // Alias para contrasenha
		Telefono      *string `json:"telefono"`
	}

	if err := c.Bind(&input); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

    if input.Resend && input.UsuarioID > 0 {
        // Obtener usuario existente
        usuario, err := a.BllController.Usuario.GetUsuario(input.UsuarioID)
        if err != nil {
            return errors.HandleError(*err, c)
        }

        // Generar nuevo código
        codigo, codeErr := repository.GenerarCodigoVerificacion()
        if codeErr != nil {
            a.Logger.Errorf("Error al generar código: %v", codeErr)
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                "error":   "CODE_GENERATION_ERROR",
                "message": "Error al generar código de verificación",
            })
        }

        // Actualizar código en BD
        expira := time.Now().Add(15 * time.Minute)
        updateErr := a.BllController.Usuario.DB.Usuario.ActualizarCodigoVerificacion(usuario.ID, codigo, expira)
        if updateErr != nil {
            a.Logger.Errorf("Error al actualizar código: %v", updateErr)
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                "error":   "CODE_UPDATE_ERROR",
                "message": "Error al actualizar código de verificación",
            })
        }

        a.Logger.Infof("Código reenviado para usuario %d: %s", usuario.ID, codigo)

        // Retornar nuevo código
        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "Código reenviado exitosamente",
            "usuario": map[string]interface{}{
                "id":     usuario.ID,
                "nombre": usuario.Nombre,
                "correo": usuario.Correo,
            },
            "codigo_verificacion": codigo,
        })
    }

	// Usar email si correo está vacío
	if input.Correo == "" && input.Email != "" {
		input.Correo = input.Email
	}

	// Usar contrasena si contrasenha está vacío
	if input.Contrasenha == "" && input.Contrasena != "" {
		input.Contrasenha = input.Contrasena
	}

	var usuario model.Usuario = model.Usuario{
		Nombre:        input.Nombre,
		TipoDocumento: input.TipoDocumento,
		NumDocumento:  input.NumDocumento,
		Correo:        input.Correo,
		Contrasenha:   input.Contrasenha,
		Telefono:      input.Telefono,
		Estado:        1,
	}

	password, err := model.HashPassword(input.Contrasenha)
	if err != nil {
		return errors.HandleError(errors.InternalServerError.PasswordHashingFailed, c)
	}
	usuario.Contrasenha = password

	usuarioRegistrado, newErr := a.BllController.Usuario.RegisterUsuario(&usuario)
    if newErr != nil {
        return errors.HandleError(*newErr, c)
    }

	// Generar código de verificación
    codigo, codeErr := repository.GenerarCodigoVerificacion()
    if codeErr != nil {
        a.Logger.Errorf("Error al generar código: %v", codeErr)
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "CODE_GENERATION_ERROR",
            "message": "Error al generar código de verificación",
        })
    }

    expira := time.Now().Add(15 * time.Minute)
    updateErr := a.BllController.Usuario.DB.Usuario.ActualizarCodigoVerificacion(usuarioRegistrado.ID, codigo, expira)
    if updateErr != nil {
        a.Logger.Errorf("Error al guardar código: %v", updateErr)
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "CODE_SAVE_ERROR",
            "message": "Error al guardar código de verificación",
        })
    }

	a.Logger.Infof("Usuario registrado: %d - Código generado: %s", usuarioRegistrado.ID, codigo)

    // Retornar código al frontend para que lo envíe por email
    return c.JSON(http.StatusCreated, map[string]interface{}{
        "message": "Usuario registrado exitosamente",
        "usuario": map[string]interface{}{
            "id":             usuarioRegistrado.ID,
            "nombre":         usuarioRegistrado.Nombre,
            "correo":         usuarioRegistrado.Correo,
            "tipo_documento": usuarioRegistrado.TipoDocumento,
            "num_documento":  usuarioRegistrado.NumDocumento,
        },
        "codigo_verificacion": codigo,
        "requiere_verificacion": true,
    })
}

func (a *Api) GetUsuario(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.HandleError(errors.BadRequestError.InvalidIDParam, c)
	}

	usuario, newErr := a.BllController.Usuario.GetUsuario(int64(id))
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	var response struct {
		ID            int64                 `json:"id"`
		Nombre        string                `json:"nombre"`
		TipoDocumento string                `json:"tipo_documento"`
		NumDocumento  string                `json:"num_documento"`
		Correo        string                `json:"correo"`
		Telefono      *string               `json:"telefono"`
		Comentario    []model.Comentario    `json:"comentarios"`
		Ordenes       []model.OrdenDeCompra `json:"ordenes"`
		Roles         []model.RolUsuario    `json:"roles"`
	}

	response.ID = usuario.ID
	response.Nombre = usuario.Nombre
	response.TipoDocumento = usuario.TipoDocumento
	response.NumDocumento = usuario.NumDocumento
	response.Correo = usuario.Correo
	response.Telefono = usuario.Telefono
	response.Comentario = usuario.Comentarios
	response.Ordenes = usuario.Ordenes
	response.Roles = usuario.RolesAsignados

	return c.JSON(http.StatusOK, response)
}

func (a *Api) AuthenticateUsuario(c echo.Context) error {

	var input struct {
		Correo      string `json:"correo"`
		Contrasenha string `json:"contrasenha"`
	}

	if err := c.Bind(&input); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	usuario, newErr := a.BllController.Usuario.AuthenticateUsuario(input.Correo, input.Contrasenha)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Generar token
    token, err := a.BllController.Token.CreateToken(usuario.ID, 24*time.Hour, "authentication")
    if err != nil {
        a.Logger.Errorf("Error al generar token para usuario %d: %v", usuario.ID, err)
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "TOKEN_GENERATION_ERROR",
            "message": "Error al generar el token de autenticación",
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Autenticación exitosa",
        "token": map[string]interface{}{
            "token":  token.Plaintext,
            "expiry": token.Expiry.Unix(),
        },
        "usuario": map[string]interface{}{
            "id":             usuario.ID,
            "nombre":         usuario.Nombre,
            "correo":         usuario.Correo,
            "tipo_documento": usuario.TipoDocumento,
            "num_documento":  usuario.NumDocumento,
            "telefono":       usuario.Telefono,
        },
    })
}

func (a *Api) VerifyEmail(c echo.Context) error {
    var input struct {
        UsuarioID int64 `json:"usuario_id"`
    }

    if err := c.Bind(&input); err != nil {
        a.Logger.Error(fmt.Sprintf("Error al parsear request: %v", err))
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "error": "INVALID_REQUEST",
            "message": "Cuerpo de solicitud inválido",
        })
    }

    if input.UsuarioID == 0 {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "error":   "INVALID_USER_ID",
            "message": "ID de usuario inválido",
        })
    }

    a.Logger.Infof("Verificando email para usuario ID: %d", input.UsuarioID)

    // Marcar como verificado
    err := a.BllController.Usuario.DB.Usuario.PostgresqlDB.Model(&model.Usuario{}).
        Where("usuario_id = ?", input.UsuarioID).
        Updates(map[string]interface{}{
            "estado_de_cuenta":        1,
            "codigo_verificacion":     nil,
            "fecha_expiracion_codigo": nil,
        }).Error

    if err != nil {
        a.Logger.Error(fmt.Sprintf("Error al verificar usuario: %v", err))
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "VERIFICATION_ERROR",
            "message": "Error al verificar el usuario",
        })
    }

    // Obtener usuario actualizado
    usuario, usuarioErr := a.BllController.Usuario.GetUsuario(input.UsuarioID)
    if usuarioErr != nil {
        a.Logger.Error(fmt.Sprintf("Error al obtener usuario: %v", usuarioErr))
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "USER_NOT_FOUND",
            "message": "Usuario no encontrado",
        })
    }

    if usuario == nil {
        a.Logger.Error("GetUsuario retornó nil")
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error":   "USER_NULL",
            "message": "Error: usuario es nil",
        })
    }

    // Log para debugging
    a.Logger.Infof("Usuario obtenido: ID=%d, Nombre=%s, Correo=%s", usuario.ID, usuario.Nombre, usuario.Correo)

    // Validar campos críticos
    if usuario.Nombre == "" {
        a.Logger.Error("Usuario sin nombre")
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error": "INVALID_USER_DATA",
            "message": "Usuario con datos incompletos",
        })
    }

    if usuario.Correo == "" {
        a.Logger.Error("Usuario sin correo")
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "error": "INVALID_USER_DATA",
            "message": "Usuario sin correo",
        })
    }
    
    // Manejar telefono (puede ser nil)
    telefono := ""
    if usuario.Telefono != nil {
        telefono = *usuario.Telefono
    }

    // Generar token simple (sin dependencias)
    tokenString := fmt.Sprintf("nexivent_auth_%d_%d", input.UsuarioID, time.Now().Unix())
    tokenExpiry := time.Now().Add(24 * time.Hour)

    a.Logger.Infof("Usuario %d verificado exitosamente", input.UsuarioID)

    // Preparar respuesta con validaciones
    usuarioResponse := map[string]interface{}{
        "id":      usuario.ID,
        "nombre":  usuario.Nombre,
        "correo":  usuario.Correo,
    }

    // Agregar campos opcionales solo si existen
    if usuario.TipoDocumento != "" {
        usuarioResponse["tipo_documento"] = usuario.TipoDocumento
    }

    if usuario.NumDocumento != "" {
        usuarioResponse["num_documento"] = usuario.NumDocumento
    }

    if telefono != "" {
        usuarioResponse["telefono"] = telefono
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Email verificado exitosamente",
        "token": map[string]interface{}{
            "token":  tokenString,
            "expiry": tokenExpiry.Unix(),
        },
        "usuario": usuarioResponse,
    })
}

func (a *Api) Logout(c echo.Context) error {
    // Obtener el token del header Authorization
    authHeader := c.Request().Header.Get("Authorization")
    if authHeader == "" {
        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "Sesión cerrada exitosamente",
        })
    }

    // Extraer el token (formato: "Bearer <token>")
    var token string
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        token = authHeader[7:]
    }

    a.Logger.Infof("Usuario cerró sesión con token: %s", token)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Sesión cerrada exitosamente",
    })
}