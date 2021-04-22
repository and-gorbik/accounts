package app

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"

	"accounts/app/controller"
)

func Serve() error {
	router := Router(&controller.Controller{})
	return http.ListenAndServe("0.0.0.0:8888", router)
}

func Router(c Controller) http.Handler {
	router := chi.NewRouter()
	router.Route("/accounts", func(r chi.Router) {
		r.Get("/filter/", c.FilterAccounts)
		r.Get("/group/", c.GroupAccounts)
		r.Get("/{id}/recommend/", c.GetRecommends)
		r.Get("/{id}/suggest/", c.GetSuggestions)
		r.Post("/new/", c.CreateAccount)
		r.Post("/{id}/", c.UpdateAccount)
		r.Post("/likes/", c.AddLikes)
	})

	return router
}

type Controller interface {
	FilterAccounts(w http.ResponseWriter, r *http.Request)
	GroupAccounts(w http.ResponseWriter, r *http.Request)
	GetRecommends(w http.ResponseWriter, r *http.Request)
	GetSuggestions(w http.ResponseWriter, r *http.Request)
	CreateAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	AddLikes(w http.ResponseWriter, r *http.Request)
}
