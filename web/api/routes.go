package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	chirender "github.com/go-chi/render"
	"github.com/mholt/binding"
	"github.com/pkg/errors"

	"git.iiens.net/edouardparis/town/app"
	"git.iiens.net/edouardparis/town/failures"
	"git.iiens.net/edouardparis/town/logging"
	"git.iiens.net/edouardparis/town/web/middlewares"
)

func NewRouter(a *app.App) http.Handler {
	r := chi.NewRouter()
	r.Route("/articles", articlesRoutes(a))
	return r
}

func render(w http.ResponseWriter, r *http.Request, resource interface{}, httpStatus int) error {
	chirender.Status(r, httpStatus)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(resource)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	return err
}

type view func(*app.App, middlewares.HandleError) http.HandlerFunc

func newHandle(a *app.App) func(view) http.HandlerFunc {
	return func(fn view) http.HandlerFunc {
		return fn(a, handleError(a.Logger))
	}
}

func handleError(logger logging.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}

		var status int
		switch cerr := errors.Cause(err).(type) {
		case failures.Error:
			status = cerr.Code
			err = cerr
		case binding.Errors:
			status = http.StatusBadRequest
			err = cerr
		default:
			logger.Error(cerr.Error())
			status = http.StatusInternalServerError
		}

		chirender.Status(r, status)
		chirender.JSON(w, r, err)
	}
}
