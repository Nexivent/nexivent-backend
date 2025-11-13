package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal/middleware"
	"github.com/Nexivent/nexivent-backend/internal/settings"
	"github.com/julienschmidt/httprouter"
)

func Routes(app *settings.Application) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.NotFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthcheckHandler)

	// Rutas para EVENTOS
	// router.HandlerFunc(http.MethodGet, "/v1/evento/:id", getEvento)
	router.HandlerFunc(http.MethodGet, "/v1/eventos", getEventos)
	router.HandlerFunc(http.MethodPost, "/v1/evento", postEvento)
	// router.HandlerFunc(http.MethodPatch, "/v1/evento/:id", patchEvento)

	// Rutas para FECHAS
	router.HandlerFunc(http.MethodGet, "/v1/fecha", getFechaPorFecha)

	// Rutas para CATEGORIAS
	router.HandlerFunc(http.MethodGet, "/v1/categorias", getCategorias)

	// Rutas para USUARIOS
	router.HandlerFunc(http.MethodGet, "/v1/usuarios", getUsuarios)
	router.HandlerFunc(http.MethodGet, "/v1/usuario/:id", getUsuarioPorID)
	router.HandlerFunc(http.MethodGet, "/v1/usuario", getUsuarioPorCorreo)

	// Rutas para ROLES
	router.HandlerFunc(http.MethodGet, "/v1/roles", getRoles)
	router.HandlerFunc(http.MethodGet, "/v1/rol", getRolPorNombre)

	// Rutas para CUPONES
	router.HandlerFunc(http.MethodGet, "/v1/cupones", getCupones)
	router.HandlerFunc(http.MethodPost, "/v1/cupon/", postCupon)
	router.HandlerFunc(http.MethodPut, "/v1/cupon/:id", putCupon)

	// Rutas para METODOS DE PAGO
	router.HandlerFunc(http.MethodGet, "/v1/metodos-pago", getMetodosPago)
	router.HandlerFunc(http.MethodGet, "/v1/metodo-pago/:id", getMetodoPagoPorID)

	// Rutas para TARIFAS
	router.HandlerFunc(http.MethodGet, "/v1/tarifas", getTarifasPorIDs)

	// Rutas para ORDENES DE COMPRA
	router.HandlerFunc(http.MethodGet, "/v1/orden-compra/:id", getOrdenDeCompraPorID)

	// Rutas para TICKETS
	router.HandlerFunc(http.MethodGet, "/v1/tickets/orden/:id", getTicketsPorOrden)

	// Aplicar el middleware para inyectar la aplicación en el contexto
	return middleware.InjectApplication(app)(router)
}
