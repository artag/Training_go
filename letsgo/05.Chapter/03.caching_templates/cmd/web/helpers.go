package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Helper. Writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {

	// debug.Stack() get a stack trace.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // 2 - call depth. To ignore this method in stack trace

	// http.StatusText() generate a human-friendly text representation
	// of a given HTTP status code.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Helper. Sends a specific status code and description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Helper. Sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Render templates from cache
func (app *application) render(
	w http.ResponseWriter, status int, page string, data *templateData) {

	// Retrieve the template set from the cache.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the templat %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Write out the provided HTTP status.
	w.WriteHeader(status)

	// Execute the template set and write the response body.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
