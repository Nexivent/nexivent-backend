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
				// Use the builtin recover function to check if there has been a panic or
				// not.
				if err := recover(); err != nil {
					// If there was a panic, set a "Connection: close" header on the
					// response. This acts as a trigger to make Go's HTTP server
					// automatically close the current connection after a response has been
					// sent.
					w.Header().Set("Connection", "close")
					// The value returned by recover() has the type any, so we use
					// fmt.Errorf() to normalize it into an error and call our
					// serverErrorResponse() helper. In turn, this will log the error using
					// our custom Logger type at the ERROR level and send the client a 500
					// Internal Server Error response.
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
