package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	port string
}

func (app *application) mount() http.Handler {
	
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
 	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)
	r.Route("/v1", func (r chi.Router) {
		r.Get("/health", app.health)
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