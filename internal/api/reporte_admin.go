package api

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// GetAdminReports GET /api/admin/reports (Usando Body para filtros complejos)
func (a *Api) GetAdminReports(c echo.Context) error {
	var req schemas.AdminReportRequest

	// Bind Body
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Call Adapter
	resp, err := a.BllController.Evento.GenerarReporteAdministrativo(req)

	// Manejo de errores
	if err != nil {
		// Si es el caso especial de "NO_DATA" (204)
		if err.Code == "NO_DATA" {
			// Opción A: Retornar 204 (Sin body según estándar)
			// return c.NoContent(http.StatusNoContent)

			// Opción B: Retornar 200/404 con JSON de error (según tu ejemplo)
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    err.Code,
					"message": err.Message,
				},
			})
		}
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, resp)
}
