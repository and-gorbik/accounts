package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"

	"accounts/app/controller"
	"accounts/app/repository"
	"accounts/app/service"
)

func Serve() error {
	connStr := flag.String("conn", "", "connection string")
	flag.Parse()
	if connStr == nil || *connStr == "" {
		return fmt.Errorf("connection string is empty")
	}

	conn, err := pgxpool.Connect(context.Background(), *connStr)
	if err != nil {
		return err
	}

	router := Router(
		controller.New(
			service.New(
				repository.New(conn),
			),
		),
	)

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
