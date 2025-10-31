package routing

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/Nexivent/nexivent-backend/internal/middleware"
	"github.com/Nexivent/nexivent-backend/internal/settings"
)

func Routes(app *settings.Application) http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/eventos/:id", getEvent)

	// Aplicar el middleware para inyectar la aplicación en el contexto
	return middleware.InjectApplication(app)(router)
}
