package settings

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Nexivent/nexivent-backend/internal"
)

type Application struct {
	Config Config
	Logger *slog.Logger
}

// El método logError() es un helper genérico para registrar un mensaje de error junto
// con el método de la request actual y la URL como atributos en la entrada del log.
func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.Logger.Error(err.Error(), "method", method, "uri", uri)
}

// El método errorResponse() es un helper genérico para enviar mensajes de error
// formateados en JSON al cliente con un código de estado dado. Nota que estamos usando el tipo
// any para el parámetro message, en lugar de solo un tipo string, ya que esto nos
// da más flexibilidad sobre los valores que podemos incluir en la respuesta.
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := internal.Envelope{"error": message}
	// Escribe la respuesta usando el helper writeJSON(). Si esto devuelve un
	// error, entonces lo registra y recurre a enviar al cliente una respuesta vacía con un
	// código de estado 500 Internal Server Error.
	err := internal.WriteJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// El método serverErrorResponse() se usará cuando nuestra Application encuentre un
// problema inesperado en tiempo de ejecución. Registra el mensaje de error detallado, luego usa el
// helper errorResponse() para enviar un código de estado 500 Internal Server Error y una respuesta
// JSON (que contiene un mensaje de error genérico) al cliente.
func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// El método notFoundResponse() se usará para enviar un código de estado 404 Not Found y
// una respuesta JSON al cliente.
func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// El método methodNotAllowedResponse() se usará para enviar un código de estado 405 Method Not
// Allowed y una respuesta JSON al cliente.
func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// El método BadRequestResponse() se usará para enviar un código de estado 400 Bad Request
// y una respuesta JSON al cliente con el mensaje de error proporcionado.
func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
