package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/controller"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

func (a *Api) RegisterUsuario(c echo.Context) error {

	var input struct {
		//para reenvio de codigo de verif
		Resend    bool  `json:"resend"`
		UsuarioID int64 `json:"usuario_id"`
		//campos de registro normal
		Nombre        string  `json:"nombre"`
		TipoDocumento string  `json:"tipo_documento"`
		NumDocumento  string  `json:"num_documento"`
		Correo        string  `json:"correo"`
		Email         string  `json:"email"`
		Contrasenha   string  `json:"contrasenha"`
		Contrasena    string  `json:"contrasena"`
		Telefono      *string `json:"telefono"`
		CuentaDeBanco *string `json:"cuenta_de_banco"`
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
	esOrganizador := input.TipoDocumento == "RUC_EMPRESA"

	var usuario model.Usuario = model.Usuario{
		Nombre:        input.Nombre,
		TipoDocumento: input.TipoDocumento,
		NumDocumento:  input.NumDocumento,
		Correo:        input.Correo,
		Contrasenha:   input.Contrasenha,
		Telefono:      input.Telefono,
		Estado:        1,
	}
	// Solo organizador puede tener cuenta de banco
	if esOrganizador && input.CuentaDeBanco != nil && *input.CuentaDeBanco != "" {
		usuario.CuentaDeBanco = input.CuentaDeBanco
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
		"codigo_verificacion":   codigo,
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
		Interaccion   []model.Interaccion   `json:"interacion"`
		Ordenes       []model.OrdenDeCompra `json:"ordenes"`
		Roles         []model.RolUsuario    `json:"roles"`
	}

	response.ID = usuario.ID
	response.Nombre = usuario.Nombre
	response.TipoDocumento = usuario.TipoDocumento
	response.NumDocumento = usuario.NumDocumento
	response.Correo = usuario.Correo
	response.Telefono = usuario.Telefono
	response.Interaccion = usuario.Interaccion
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Datos inválidos",
		})
	}

	usuario, newErr := a.BllController.Usuario.AuthenticateUsuario(input.Correo, input.Contrasenha)
	if newErr != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": newErr.Message,
		})
	}

	if usuario.Estado != 1 {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":   "ACCOUNT_DISABLED",
			"message": "Tu cuenta ha sido deshabilitada. Contacta al soporte.",
		})
	}

	// Obtener roles desde el controller (adapter -> repository)
	rolesResp, rolErr := a.BllController.RolUsuario.GetUserRoles(usuario.ID)
	if rolErr != nil {
		rolesResp = &schemas.RolUsuarioResponse{IDUsuario: usuario.ID, Roles: []schemas.RolResponse{}}
	}

	// extraer slice de roles para uso y respuesta
	var roles []schemas.RolResponse
	if rolesResp != nil && len(rolesResp.Roles) > 0 {
		roles = rolesResp.Roles
	} else {
		roles = []schemas.RolResponse{}
	}

	// Determinar rol principal
	var rolPrincipal string
	if len(roles) > 0 {
		for _, r := range roles {
			// comparar con el nombre que usas para admin
			if r.Nombre == "ADMINISTRADOR" {
				rolPrincipal = "ADMINISTRADOR"
				break
			}
		}
		if rolPrincipal == "" {
			rolPrincipal = roles[0].Nombre
		}
	} else {
		rolPrincipal = "ASISTENTE"
	}

	// Generar token
	token, err := a.BllController.Token.CreateToken(usuario.ID, 24*time.Hour, "authentication")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "TOKEN_GENERATION_ERROR",
			"message": "Error al generar el token de autenticación",
		})
	}

	response := map[string]interface{}{
		"message": "Autenticación exitosa",
		"token": map[string]interface{}{
			"token":  token.Plaintext,
			"expiry": time.Now().Add(24 * time.Hour).Unix(),
		},
		"usuario": map[string]interface{}{
			"id":             usuario.ID,
			"correo":         usuario.Correo,
			"nombre":         usuario.Nombre,
			"num_documento":  usuario.NumDocumento,
			"telefono":       usuario.Telefono,
			"tipo_documento": usuario.TipoDocumento,
			"roles":          roles,
			"rol_principal":  rolPrincipal,
		},
	}
	return c.JSON(http.StatusOK, response)
}

func (a *Api) AuthenticateOrganizador(c echo.Context) error {

	var input struct {
		Ruc         string `json:"ruc"`
		Contrasenha string `json:"contrasenha"`
	}

	if err := c.Bind(&input); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Validar que los campos no estén vacíos
	if input.Ruc == "" || input.Contrasenha == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "EMPTY_CREDENTIALS",
			"message": "RUC y contraseña son requeridos",
		})
	}

	// Validar formato de RUC (11 dígitos numéricos)
	if len(input.Ruc) != 11 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "INVALID_RUC_FORMAT",
			"message": "El RUC debe tener 11 dígitos",
		})
	}

	usuario, newErr := a.BllController.Usuario.AuthenticateOrganizador(input.Ruc, input.Contrasenha)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Validar que el usuario exista
	if usuario == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "INVALID_CREDENTIALS",
			"message": "Credenciales incorrectas",
		})
	}

	if usuario.Estado != 1 {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":   "ACCOUNT_DISABLED",
			"message": "Tu cuenta ha sido deshabilitada. Contacta al soporte.",
		})
	}

	// Validar que el usuario esté activo
	if usuario.Estado != 1 {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":   "ACCOUNT_DISABLED",
			"message": "Tu cuenta ha sido deshabilitada. Contacta al soporte.",
		})
	}

	// Generar token
	token, err := a.BllController.Token.CreateToken(usuario.ID, 24*time.Hour, "authentication")
	if err != nil {
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
			"estado_cuenta":  usuario.EstadoDeCuenta,
		},
	})
}

func (a *Api) GoogleAuth(c echo.Context) error {
	var input struct {
		AccessToken   string `json:"access_token"`
		IdToken       string `json:"id_token"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
		Sub           string `json:"sub"`
		TipoDocumento string `json:"tipo_documento"`
		NumDocumento  string `json:"num_documento"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "INVALID_REQUEST",
			"message": "Solicitud inválida",
		})
	}

	var googleUser *controller.GoogleUser

	if input.Email != "" && input.EmailVerified && input.Sub != "" {
		googleUser = &controller.GoogleUser{
			Email:         input.Email,
			Name:          input.Name,
			Picture:       input.Picture,
			Sub:           input.Sub,
			VerifiedEmail: input.EmailVerified,
		}
	} else if input.IdToken != "" {
		var err error
		googleUser, err = a.BllController.Usuario.VerifyGoogleToken(input.IdToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "INVALID_GOOGLE_TOKEN",
				"message": "Token de Google inválido o expirado",
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "MISSING_DATA",
			"message": "Se requiere id_token o datos validados de Google",
		})
	}

	// Verificar que el email esté verificado por Google
	if !googleUser.VerifiedEmail {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "EMAIL_NOT_VERIFIED",
			"message": "El correo de Google no está verificado",
		})
	}

	usuarioExistente, err := a.BllController.Usuario.DB.Usuario.ObtenerUsuarioPorCorreo(googleUser.Email)

	if usuarioExistente.Estado != 1 {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":   "ACCOUNT_DISABLED",
			"message": "Tu cuenta ha sido deshabilitada. Contacta al soporte.",
		})
	}

	if err == nil && usuarioExistente != nil {
		token, tokenErr := a.BllController.Token.CreateToken(usuarioExistente.ID, 24*time.Hour, "authentication")
		if tokenErr != nil {
			a.Logger.Errorf("Error al generar token: %v", tokenErr)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "TOKEN_GENERATION_ERROR",
				"message": "Error al generar token de autenticación",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Autenticación con Google exitosa",
			"token": map[string]interface{}{
				"token":  token.Plaintext,
				"expiry": token.Expiry.Unix(),
			},
			"usuario": map[string]interface{}{
				"id":               usuarioExistente.ID,
				"nombre":           usuarioExistente.Nombre,
				"correo":           usuarioExistente.Correo,
				"tipo_documento":   usuarioExistente.TipoDocumento,
				"num_documento":    usuarioExistente.NumDocumento,
				"telefono":         usuarioExistente.Telefono,
				"estado_de_cuenta": usuarioExistente.EstadoDeCuenta,
			},
			"is_new_user": false,
		})
	}

	var nuevoUsuario model.Usuario

	nuevoUsuario.Nombre = googleUser.Name
	nuevoUsuario.Correo = googleUser.Email
	nuevoUsuario.TipoDocumento = input.TipoDocumento
	nuevoUsuario.NumDocumento = input.NumDocumento
	nuevoUsuario.EstadoDeCuenta = 1
	nuevoUsuario.Estado = 1
	nuevoUsuario.Contrasenha = ""
	nuevoUsuario.Telefono = nil

	usuarioRegistrado, newErr := a.BllController.Usuario.RegisterUsuario(&nuevoUsuario)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	token, tokenErr := a.BllController.Token.CreateToken(usuarioRegistrado.ID, 24*time.Hour, "authentication")
	if tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "TOKEN_GENERATION_ERROR",
			"message": "Error al generar token de autenticación",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Registro con Google exitoso",
		"token": map[string]interface{}{
			"token":  token.Plaintext,
			"expiry": token.Expiry.Unix(),
		},
		"usuario": map[string]interface{}{
			"id":             usuarioRegistrado.ID,
			"nombre":         usuarioRegistrado.Nombre,
			"correo":         usuarioRegistrado.Correo,
			"tipo_documento": usuarioRegistrado.TipoDocumento,
			"num_documento":  usuarioRegistrado.NumDocumento,
			"telefono":       usuarioRegistrado.Telefono,
		},
		"is_new_user":       true,
		"skip_verification": true,
	})
}

func (a *Api) VerifyEmail(c echo.Context) error {
	var input struct {
		UsuarioID int64 `json:"usuario_id"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "INVALID_REQUEST",
			"message": "Cuerpo de solicitud inválido",
		})
	}

	if input.UsuarioID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "INVALID_USER_ID",
			"message": "ID de usuario inválido",
		})
	}

	usuarioPrev, usuarioErr := a.BllController.Usuario.GetUsuario(input.UsuarioID)
	if usuarioErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "USER_NOT_FOUND",
			"message": "Usuario no encontrado",
		})
	}
	var estadoCuenta int
	if usuarioPrev.TipoDocumento == "RUC_PERSONA" || usuarioPrev.TipoDocumento == "RUC_EMPRESA" {
		estadoCuenta = 0
	} else {
		estadoCuenta = 1
	}
	// Marcar como verificado
	err := a.BllController.Usuario.DB.Usuario.PostgresqlDB.Model(&model.Usuario{}).
		Where("usuario_id = ?", input.UsuarioID).
		Updates(map[string]interface{}{
			"estado_de_cuenta":        estadoCuenta,
			"codigo_verificacion":     nil,
			"fecha_expiracion_codigo": nil,
		}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "VERIFICATION_ERROR",
			"message": "Error al verificar el usuario",
		})
	}

	// Obtener usuario actualizado
	usuario, usuarioErr := a.BllController.Usuario.GetUsuario(input.UsuarioID)
	if usuarioErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "USER_NOT_FOUND",
			"message": "Usuario no encontrado",
		})
	}

	if usuario == nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "USER_NULL",
			"message": "Error: usuario es nil",
		})
	}

	// Validar campos críticos
	if usuario.Nombre == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "INVALID_USER_DATA",
			"message": "Usuario con datos incompletos",
		})
	}

	if usuario.Correo == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "INVALID_USER_DATA",
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

	// Preparar respuesta con validaciones
	usuarioResponse := map[string]interface{}{
		"id":     usuario.ID,
		"nombre": usuario.Nombre,
		"correo": usuario.Correo,
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

// ===================================================
// GESTIÓN DE ESTADO DE USUARIOS
// ===================================================

// ActivarUsuario activa un usuario (estado = 1)
// POST /api/users/:id/activate
func (a *Api) ActivarUsuario(c echo.Context) error {
	idParam := c.Param("id")
	usuarioID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		a.Logger.Warnf("ID de usuario inválido: %s", idParam)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "ID de usuario inválido",
		})
	}

	// TODO: Obtener del contexto/sesión
	updatedBy := int64(1)

	apiErr := a.BllController.Usuario.ActivarUsuario(usuarioID, updatedBy)
	if apiErr != nil {
		a.Logger.Errorf("Error activando usuario %d: %v", usuarioID, apiErr)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": apiErr.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Usuario activado correctamente",
	})
}

// DesactivarUsuario desactiva un usuario (estado = 0)
// POST /api/users/:id/deactivate
func (a *Api) DesactivarUsuario(c echo.Context) error {
	idParam := c.Param("id")
	usuarioID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		a.Logger.Warnf("ID de usuario inválido: %s", idParam)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "ID de usuario inválido",
		})
	}

	updatedBy := int64(1)

	apiErr := a.BllController.Usuario.DesactivarUsuario(usuarioID, updatedBy)
	if apiErr != nil {
		a.Logger.Errorf("Error desactivando usuario %d: %v", usuarioID, apiErr)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": apiErr.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Usuario desactivado correctamente",
	})
}

// CambiarEstadoUsuario cambia el estado de un usuario
// POST /api/users/:id/status
// Body: { "estado": 1 } o { "estado": 0 }
func (a *Api) CambiarEstadoUsuario(c echo.Context) error {
	idParam := c.Param("id")
	usuarioID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		a.Logger.Warnf("ID de usuario inválido: %s", idParam)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "ID de usuario inválido",
		})
	}

	var request struct {
		Estado int16 `json:"estado"`
	}
	if err := c.Bind(&request); err != nil {
		a.Logger.Warnf("Error parseando request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "El campo 'estado' es requerido (0 o 1)",
		})
	}

	if request.Estado != 0 && request.Estado != 1 {
		a.Logger.Warnf("Estado inválido: %d", request.Estado)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "El estado debe ser 0 (inactivo) o 1 (activo)",
		})
	}

	updatedBy := int64(1)

	// Llamar al controller directamente
	var apiErr *errors.Error
	if request.Estado == 1 {
		apiErr = a.BllController.Usuario.ActivarUsuario(usuarioID, updatedBy)
	} else {
		apiErr = a.BllController.Usuario.DesactivarUsuario(usuarioID, updatedBy)
	}

	if apiErr != nil {
		a.Logger.Errorf("Error cambiando estado del usuario %d: %v", usuarioID, apiErr)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": apiErr.Message,
		})
	}

	estadoTexto := "activado"
	if request.Estado == 0 {
		estadoTexto = "desactivado"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Usuario %s correctamente", estadoTexto),
	})
}

// Dentro de tu archivo api_usuario.go (o donde tengas los handlers)

type actualizarPasswordRequest struct {
	NuevaContrasenha string `json:"nuevaContrasenha"`
}

func (a *Api) ActualizarContrasenha(c echo.Context) error {
	// 1) ID del path
	idParam := c.Param("id")
	usuarioID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		a.Logger.Warnf("ID de usuario inválido: %s", idParam)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "ID de usuario inválido",
		})
	}

	// 2) Body con nueva contraseña
	var req actualizarPasswordRequest
	if err := c.Bind(&req); err != nil {
		a.Logger.Warnf("Body inválido al actualizar contraseña de usuario %d: %v", usuarioID, err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Body inválido",
		})
	}

	if req.NuevaContrasenha == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "La nueva contraseña es obligatoria",
		})
	}

	// 3) updatedBy (luego lo sacas del token/JWT)
	updatedBy := usuarioID

	// 4) Llamar a la capa de negocio
	apiErr := a.BllController.Usuario.ActualizarContrasenha(usuarioID, req.NuevaContrasenha, updatedBy)
	if apiErr != nil {
		a.Logger.Errorf("Error actualizando contraseña de usuario %d: %v", usuarioID, apiErr)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": apiErr.Message,
		})
	}

	// 5) Respuesta OK
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Contraseña actualizada correctamente",
	})
}
