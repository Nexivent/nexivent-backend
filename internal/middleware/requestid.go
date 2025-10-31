package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	appcontext "github.com/Nexivent/nexivent-backend/internal/context"
)

// generateRequestID genera un ID único aleatorio para el request
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// RequestID es un middleware que genera y añade un ID único a cada request
// Útil para logging, debugging y tracing de requests a través del sistema
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Primero, intentar obtener el request ID del header (si el cliente lo envió)
		requestID := r.Header.Get("X-Request-ID")

		// Si no existe, generar uno nuevo
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Añadir el request ID al contexto
		ctx := appcontext.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		// Añadir el request ID como header de respuesta
		// para que el cliente pueda usarlo en reportes de errores
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}
