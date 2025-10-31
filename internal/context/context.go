package context

import (
	"context"

	"github.com/Nexivent/nexivent-backend/internal/settings"
)

// contextKey es un tipo personalizado para las keys del contexto
// Esto evita colisiones con otras keys en el contexto
type contextKey string

const (
	applicationContextKey contextKey = "application"
	userIDContextKey      contextKey = "userID"
	tokenContextKey       contextKey = "token"
	requestIDContextKey   contextKey = "requestID"
	organizerIDContextKey contextKey = "organizerID"
)

// WithApplication añade la instancia de Application al contexto
func WithApplication(ctx context.Context, app *settings.Application) context.Context {
	return context.WithValue(ctx, applicationContextKey, app)
}

// GetApplication recupera la instancia de Application del contexto
// Retorna nil si no existe en el contexto
func GetApplication(ctx context.Context) *settings.Application {
	app, ok := ctx.Value(applicationContextKey).(*settings.Application)
	if !ok {
		return nil
	}
	return app
}

// WithUserID añade el ID del usuario al contexto
func WithUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// GetUserID recupera el ID del usuario del contexto
// Retorna 0 y false si no existe en el contexto
func GetUserID(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uint64)
	return userID, ok
}

// WithToken añade el token de autenticación al contexto
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

// GetToken recupera el token de autenticación del contexto
// Retorna una cadena vacía y false si no existe en el contexto
func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenContextKey).(string)
	return token, ok
}

// WithRequestID añade un ID único de request al contexto (útil para logging/tracing)
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

// GetRequestID recupera el ID único del request del contexto
// Retorna una cadena vacía y false si no existe en el contexto
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	return requestID, ok
}

// WithOrganizerID añade el ID del organizador al contexto
func WithOrganizerID(ctx context.Context, organizerID uint64) context.Context {
	return context.WithValue(ctx, organizerIDContextKey, organizerID)
}

// GetOrganizerID recupera el ID del organizador del contexto
// Retorna 0 y false si no existe en el contexto
func GetOrganizerID(ctx context.Context) (uint64, bool) {
	organizerID, ok := ctx.Value(organizerIDContextKey).(uint64)
	return organizerID, ok
}
