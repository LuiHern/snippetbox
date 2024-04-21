package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"time"
)

// send a generic 500 server error response to the user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{CurrentYear: time.Now().Year(), Flash: app.sessionManager.PopString(r.Context(), "flash")}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// write the template to the buffer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// write buffer to ResponseWriter
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		// check for this specific error that can be a server side issue
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}

	return nil
}
