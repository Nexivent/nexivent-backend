package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/labstack/echo/v4"
)

func (a *Api) RegisterUsuario(c echo.Context) error {

	var input struct {
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

	// Generar token JWT después del registro exitoso
    token, tokenErr := a.BllController.Token.CreateToken(usuarioRegistrado.ID, 24*time.Hour, "authentication")
    if tokenErr != nil {
        a.Logger.Errorf("Error al generar token para usuario %d: %v", usuarioRegistrado.ID, tokenErr)
        return errors.HandleError(*tokenErr, c)
    }

	// Estructurar la respuesta con el token y los datos del usuario
    var response struct {
        ID            int64       `json:"id"`
        Nombre        string      `json:"nombre"`
        TipoDocumento string      `json:"tipo_documento"`
        NumDocumento  string      `json:"num_documento"`
        Correo        string      `json:"correo"`
        Telefono      *string     `json:"telefono"`
        Token         model.Token `json:"token"`
    }

	response.ID = usuarioRegistrado.ID
	response.Nombre = usuarioRegistrado.Nombre
	response.TipoDocumento = usuarioRegistrado.TipoDocumento
	response.NumDocumento = usuarioRegistrado.NumDocumento
	response.Correo = usuarioRegistrado.Correo
	response.Telefono = usuarioRegistrado.Telefono
	response.Token = *token

	return c.JSON(http.StatusCreated, response)
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

	var response struct {
		ID            int64       `json:"id"`
		Nombre        string      `json:"nombre"`
		TipoDocumento string      `json:"tipo_documento"`
		NumDocumento  string      `json:"num_documento"`
		Correo        string      `json:"correo"`
		Telefono      *string     `json:"telefono"`
		Token         model.Token `json:"token"`
	}

	response.ID = usuario.ID
	response.Nombre = usuario.Nombre
	response.TipoDocumento = usuario.TipoDocumento
	response.NumDocumento = usuario.NumDocumento
	response.Correo = usuario.Correo
	response.Telefono = usuario.Telefono
	token, err := a.BllController.Token.CreateToken(usuario.ID, 24*time.Hour, "authentication")
	if err != nil {
		return errors.HandleError(*err, c)
	}
	response.Token = *token

	return c.JSON(http.StatusOK, response)
}
