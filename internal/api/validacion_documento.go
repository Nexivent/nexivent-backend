package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// ValidarDocumento godoc
// @Summary 			Validar documento (DNI o RUC)
// @Description 		Valida un documento de identidad consultando la API de Factiliza
// @Tags 				Validacion
// @Accept 				json
// @Produce 			json
// @Param 				request body schemas.ValidarDocumentoRequest true "Datos del documento a validar"
// @Success 			200 {object} schemas.ValidarDocumentoResponse
// @Failure 			400 {object} map[string]interface{}
// @Failure 			500 {object} map[string]interface{}
// @Router 				/validar-documento [post]
func (a *Api) ValidarDocumento(c echo.Context) error {
    a.Logger.Info("=== ValidarDocumento endpoint called ===")

    // Leer el body completo para debugging
    bodyBytes, _ := io.ReadAll(c.Request().Body)
    a.Logger.Info("Raw body:", string(bodyBytes))
    // Recrear el body para que Bind pueda leerlo
    c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

    var req schemas.ValidarDocumentoRequest

    if err := c.Bind(&req); err != nil {
        a.Logger.Error("Error binding request:", err)
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "Error al procesar los datos",
            "error":   err.Error(),
        })
    }

    a.Logger.Info("Parsed request - TipoDocumento:", req.TipoDocumento, "NumeroDocumento:", req.NumeroDocumento)

    // Validaciones
    if req.TipoDocumento == "" {
        a.Logger.Error("tipo_documento is empty")
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "El campo 'tipo_documento' es requerido",
        })
    }

    if req.NumeroDocumento == "" {
        a.Logger.Error("numero_documento is empty")
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "El campo 'numero_documento' es requerido",
        })
    }

    // Validar tipos permitidos
    validTypes := []string{"DNI", "CE", "RUC_PERSONA", "RUC_EMPRESA"}
    isValid := false
    for _, vt := range validTypes {
        if req.TipoDocumento == vt {
            isValid = true
            break
        }
    }

    if !isValid {
        a.Logger.Error("Invalid tipo_documento:", req.TipoDocumento)
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "Tipo de documento inválido. Use: DNI, CE, RUC_PERSONA o RUC_EMPRESA",
            "received": req.TipoDocumento,
        })
    }

    a.Logger.Info("Calling controller...")

    // Verificar que el controlador exista
    if a.BllController.ValidacionDocumento == nil {
        a.Logger.Error("ValidacionDocumentoController is nil!")
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "success": false,
            "message": "Controller no inicializado",
        })
    }

    response, err := a.BllController.ValidacionDocumento.ValidarDocumento(&req)
    if err != nil {
        a.Logger.Error("Error from controller:", err)
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "success": false,
            "message": "Error interno del servidor",
            "error":   err.Error(),
        })
    }

    a.Logger.Info("Controller response success:", response.Success)

    statusCode := http.StatusOK
    if !response.Success {
        statusCode = http.StatusBadRequest
    }

    return c.JSON(statusCode, response)
}