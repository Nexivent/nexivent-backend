package middleware

import (
	"fmt"
	"net/http"

	appcontext "github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/settings"
)

// InjectApplication es un middleware que inyecta la instancia de Application
// en el contexto de cada request HTTP
func InjectApplication(app *settings.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// Usa la función incorporada recover para verificar si ha habido un pánico
				// o no.
				if err := recover(); err != nil {
					// Si hubo un pánico, establece un header "Connection: close" en la
					// respuesta. Esto actúa como un disparador para hacer que el servidor HTTP
					// de Go cierre automáticamente la conexión actual después de que se haya
					// enviado una respuesta.
					w.Header().Set("Connection", "close")
					// El valor devuelto por recover() tiene el tipo any, por lo que usamos
					// fmt.Errorf() para normalizarlo en un error y llamar a nuestro
					// helper serverErrorResponse(). A su vez, esto registrará el error usando
					// nuestro tipo Logger personalizado en el nivel ERROR y enviará al cliente
					// una respuesta 500 Internal Server Error.
					app.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
				}
			}()

			// Añadir la aplicación al contexto del request
			ctx := appcontext.WithApplication(r.Context(), app)
			// Crear un nuevo request con el contexto actualizado
			r = r.WithContext(ctx)
			// Pasar al siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}
