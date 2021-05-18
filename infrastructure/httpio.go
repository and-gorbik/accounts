package infrastructure

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ReadRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func ReadURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func WriteSuccessResponse(w http.ResponseWriter, body []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		w.Write(body)
	}
}
