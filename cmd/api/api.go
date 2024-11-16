package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/puremike/social-go/internal/store"
)

type application struct {
	config config
	store store.Storage
}

type config struct {
	port string
	dbconfig dbconfig
	environment string
}

type dbconfig struct {
	Addr string
	maxOpenConns, maxIdleConns int
	maxIdleTime time.Duration
}

func (app *application) mount() http.Handler {
	
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
 	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)
	r.Route("/v1", func (r chi.Router) {
		r.Get("/health", app.health)
		r.Route("/posts", func (r chi.Router) {
			r.Post("/", app.CreatePost)
		})
	})

	return r
}

func (app *application) start(mux http.Handler) error {

	srv := &http.Server{
		Addr: app.config.port,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Starting server on port %s...\n", app.config.port)
	return srv.ListenAndServe()
}