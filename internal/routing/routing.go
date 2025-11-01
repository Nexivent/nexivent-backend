package routing

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal/middleware"
	"github.com/Nexivent/nexivent-backend/internal/settings"
	"github.com/julienschmidt/httprouter"
)

func Routes(app *settings.Application) http.Handler {
	router := httprouter.New()
	
	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.NotFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)
	
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/eventos/:id", getEvento)
	router.HandlerFunc(http.MethodPut, "/v1/eventos/", postEvento)

	// Aplicar el middleware para inyectar la aplicación en el contexto
	return middleware.InjectApplication(app)(router)
}
