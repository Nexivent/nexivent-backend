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
	router.HandlerFunc(http.MethodGet, "/v1/evento/:id", getEvento)
	router.HandlerFunc(http.MethodGet, "/v1/eventos", getEventos)
	router.HandlerFunc(http.MethodPost, "/v1/evento", postEvento)
	router.HandlerFunc(http.MethodPatch, "/v1/evento/:id", patchEvento)

	// Aplicar el middleware para inyectar la aplicación en el contexto
	return middleware.InjectApplication(app)(router)
}
