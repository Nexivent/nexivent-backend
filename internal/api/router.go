package api

import (
	"fmt"

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
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // TODO: allow only authorized origins
		AllowCredentials: true,
	}
	a.Echo.Use(middleware.CORSWithConfig(corsConfig))

	// Enable Swagger if configured
	if configEnv.EnableSwagger {
		a.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("server")))
	}

	// ===== PUBLIC ENDPOINTS =====
	healthCheck := a.Echo.Group("/health-check")
	healthCheck.GET("/", a.HealthCheck)

	// Eventos endpoints
	a.Echo.GET("/evento/", a.FetchEventos)
	a.Echo.GET("/evento/:eventoId/", a.GetEvento)
	a.Echo.POST("/evento/", a.CreateEvento) //falta usuario creacion

	a.Echo.GET("/categorias/", a.FetchCategorias)
	a.Echo.POST("/categoria/", a.CreateCategoria)
	a.Echo.GET("/categoria/:categoriaId/", a.GetCategoria)

	// Media uploads
	a.Echo.POST("/media/upload-url", a.GenerateUploadURL)

	//Cupon
	a.Echo.POST("/cupon/:usuarioCreacion", a.CreateCupon)

	a.Echo.POST("/register", a.RegisterUsuario)
	a.Echo.GET("/usuario/:id", a.GetUsuario)

	a.Echo.POST("/orden_de_compra/hold", a.CrearSesionOrdenTemporal)
	a.Echo.GET("/orden_de_compra/:orderId/hold", a.ObtenerEstadoHold)
	a.Echo.POST("/orden_de_compra/:orderId/confirm", a.ConfirmarOrden)
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
