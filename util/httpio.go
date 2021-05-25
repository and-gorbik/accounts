package util

import (
	"io/ioutil"
	"log"
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
	w.WriteHeader(status)
	log.Println(err)
}

func WriteSuccessResponse(w http.ResponseWriter, body []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		w.Write(body)
	}
}
