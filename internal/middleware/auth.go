package middleware

import (
	"net/http"
	"strings"

	"github.com/Nexivent/nexivent-backend/internal/context"
)

// Authentication es un middleware que extrae y valida el token de autenticación
// y añade el userID al contexto si el token es válido
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token del header Authorization
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// No hay token, continuar sin autenticación
			next.ServeHTTP(w, r)
			return
		}

		// Esperar formato: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// TODO: Aquí validarías el token con tu sistema de autenticación
		// Por ahora, esto es solo un ejemplo
		userID, valid := validateToken(token)
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Añadir el token y userID al contexto
		ctx := context.WithToken(r.Context(), token)
		ctx = context.WithUserID(ctx, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RequireAuthentication es un middleware que REQUIERE autenticación
// Debe usarse después del middleware Authentication
func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar si hay un userID en el contexto
		_, ok := context.GetUserID(r.Context())
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateToken es una función placeholder para validar tokens
// TODO: Implementar con tu lógica real de autenticación (JWT, sesiones, etc.)
func validateToken(token string) (userID uint64, valid bool) {
	// Ejemplo simplificado - NO USAR EN PRODUCCIÓN
	// Aquí deberías:
	// 1. Validar el JWT o consultar tu base de datos de sesiones
	// 2. Verificar que no haya expirado
	// 3. Extraer el userID del token

	if token == "valid-test-token" {
		return 1, true
	}

	return 0, false
}
