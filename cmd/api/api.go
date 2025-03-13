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
	store  store.Storage
}

type config struct {
	port        string
	dbconfig    dbconfig
	environment string
}

type dbconfig struct {
	Addr                       string
	maxOpenConns, maxIdleConns int
	maxIdleTime                time.Duration
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes) // Automatically removes trailing slashes
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.health)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPost)
			r.Get("/", app.getAllPosts)
			r.Delete("/", app.deleteAllPosts)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", app.getPostById)
				r.Delete("/", app.deletePostByID)
				r.Patch("/", app.updatePost)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUser)

			r.Group(func(r chi.Router) {
				r.Get("/{id}/feed", app.getUserFeedsHandler)
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.userContextMiddleWare)
				r.Get("/", app.getUserByID)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})

		})
	})

	return r
}

func (app *application) start(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on port %s...\n", app.config.port)
	return srv.ListenAndServe()
}
