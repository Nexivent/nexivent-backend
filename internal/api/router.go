package api

import (
	"fmt"
	"os"
	"strings"

	config "github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// HealthCheck godoc
// @Summary 			Health Check
// @Description 		Check the health status of the API
// @Tags 				Health
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} map[string]string "API is healthy"
// @Router 				/health-check/ [get]
func (a *Api) HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "OK",
	})
}

func (a *Api) RegisterRoutes(configEnv *config.ConfigEnv) {
	allowOrigins := []string{"http://localhost:3000", "http://localhost:3001", "https://accounts.google.com"}

	if extraOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); extraOrigins != "" {
		if extraOrigins == "*" {
			allowOrigins = []string{"*"}
		} else {
			for _, origin := range strings.Split(extraOrigins, ",") {
				trimmed := strings.TrimSpace(origin)
				if trimmed != "" {
					allowOrigins = append(allowOrigins, trimmed)
				}
			}
		}
	}

	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Requested-With", "X-CSRF-Token"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Authorization",
		},
		MaxAge: 86400,
	}
	a.Echo.Use(middleware.CORSWithConfig(corsConfig))

	// Enable Swagger if configured
	if configEnv.EnableSwagger {
		a.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("server")))
	}

	// ===== PUBLIC ENDPOINTS =====
	healthCheck := a.Echo.Group("/health-check")
	healthCheck.GET("/", a.HealthCheck)

	// Document validation endpoint
	a.Echo.POST("/validar-documento", a.ValidarDocumento)

	// Usuario endpoints
	a.Echo.POST("/register", a.RegisterUsuario)
	a.Echo.POST("/google-auth", a.GoogleAuth)
	a.Echo.POST("/verify-email", a.VerifyEmail)

	// Autenticación
	a.Echo.POST("/login", a.AuthenticateUsuario)
	a.Echo.POST("/loginorg", a.AuthenticateOrganizador)
	a.Echo.POST("/logout", a.Logout)

	a.Echo.GET("/usuario/:id", a.GetUsuario)

	// Eventos endpoints
	a.Echo.GET("/evento/", a.FetchEventos)
	a.Echo.GET("/evento/:eventoId/", a.GetEvento)
	a.Echo.POST("/evento/", a.CreateEvento) //falta usuario creacion
	a.Echo.GET("/evento/filter", a.FetchEventosWithFilters)
	a.Echo.GET("/evento/reporte", a.GetReporteEvento)
	a.Echo.GET("/organizador/:organizadorId/eventos/reporte", a.GetReporteEventosOrganizador)
	a.Echo.GET("/api/events/:id/summary", a.GetEventoSummary)

	a.Echo.GET("/categorias/", a.FetchCategorias)
	a.Echo.POST("/categoria/", a.CreateCategoria)
	a.Echo.GET("/categoria/:categoriaId/", a.GetCategoria)
	// 2. Reporte Administrativo Global (Dashboard BI)
	a.Echo.POST("/api/admin/reports", a.GetAdminReports)
	// Media uploads
	a.Echo.POST("/media/upload-url", a.GenerateUploadURL)
	//Cupones
	//Cupon
	a.Echo.POST("/cupon/:usuarioCreacion", a.CreateCupon)
	a.Echo.PUT("/cupon/:usuarioModificacion", a.UpdateCupon)
	a.Echo.GET("/cupon/organizador/:organizadorId", a.FetchCuponPorOrganizador)
	a.Echo.GET("/cupon/validar", a.ValidateCupon)
	a.Echo.POST("/cupon/usuario", a.CreateUsuarioCuponForOrdenCompra)

	//Orden de compra
	a.Echo.POST("/orden_de_compra/hold", a.CrearSesionOrdenTemporal)
	a.Echo.GET("/orden_de_compra/:orderId/hold", a.ObtenerEstadoHold)
	a.Echo.POST("/orden_de_compra/:orderId/confirm", a.ConfirmarOrden)

	// Perfiles de persona
	a.Echo.GET("/evento/:eventoId/perfiles", a.ListarPerfilesPorEvento)
	a.Echo.POST("/evento/:eventoId/perfiles", a.CrearPerfilPersona)
	a.Echo.PUT("/perfiles/:perfilId", a.ActualizarPerfilPersona)

	// Sectores
	a.Echo.GET("/evento/:eventoId/sectores", a.ListarSectoresPorEvento)
	a.Echo.POST("/evento/:eventoId/sectores", a.CrearSector)
	a.Echo.PUT("/sectores/:sectorId", a.ActualizarSector)

	// Tipos de ticket
	a.Echo.GET("/evento/:eventoId/tipos-ticket", a.ListarTiposTicketPorEvento)
	a.Echo.POST("/evento/:eventoId/tipos-ticket", a.CrearTipoTicket)
	a.Echo.PUT("/tipos-ticket/:tipoTicketId", a.ActualizarTipoTicket)

	// Tarifas
	a.Echo.POST("/tarifas", a.CrearTarifa)
	a.Echo.PUT("/tarifas/:tarifaId", a.ActualizarTarifa)

	// Tickets
	a.Echo.POST("/api/tickets/issue", a.EmitirTickets)
	a.Echo.POST("/api/tickets/cancel", a.CancelarTickets)

	//Roles
	a.Echo.GET("/rol/:nombre/name", a.GetRolPorNombre)
	a.Echo.GET("/rol/:usuarioId/user", a.GetRolPorUsuario)
	a.Echo.PUT("/rol/:userId/update/:rolId", a.UpdateRol)
	a.Echo.GET("/roles/", a.FetchRoles)

	//roles_usuario
	a.Echo.GET("/api/users/:id/roles", a.ListarRolesDeUsuario)
	a.Echo.POST("/api/roles/assign", a.CreateRolUser)
	a.Echo.DELETE("/api/roles/revoke", a.DeleteRolUser)
	a.Echo.GET("/api/users", a.ListarUsuariosPorRol)

	// Gestión de estado de usuarios
	a.Echo.POST("/api/users/:id/status", a.CambiarEstadoUsuario)
	a.Echo.POST("/api/users/:id/activate", a.ActivarUsuario)
	a.Echo.POST("/api/users/:id/deactivate", a.DesactivarUsuario)

}

func (a *Api) RunApi(configEnv *config.ConfigEnv) {
	a.RegisterRoutes(configEnv)

	// Start the server
	port := configEnv.MainPort
	if port == "" {
		port = "8080"
	}
	a.Logger.Infoln(fmt.Sprintf("Nexivent server running on port %s", port))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", port)))
}
