package middleware

import (
	"net/http"

	appcontext "github.com/Nexivent/nexivent-backend/internal/context"
	"github.com/Nexivent/nexivent-backend/internal/settings"
)

// InjectApplication es un middleware que inyecta la instancia de Application
// en el contexto de cada request HTTP
func InjectApplication(app *settings.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Añadir la aplicación al contexto del request
			ctx := appcontext.WithApplication(r.Context(), app)
			// Crear un nuevo request con el contexto actualizado
			r = r.WithContext(ctx)
			// Pasar al siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}
