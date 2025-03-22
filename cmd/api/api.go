package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/puremike/social-go/docs"
	"github.com/puremike/social-go/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	port        string
	dbconfig    dbconfig
	environment string
	apiUrl string
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

		docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)
		r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(docURL), //The url pointing to API definition
	))

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
				r.Get("/{id}/feeds", app.getUserFeedsHandler)
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.userContextMiddleWare)
				r.Get("/", app.getUserByID)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})
		})

		r.Route("/authentication", func (r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})


	})

	return r
}

func (app *application) start(mux http.Handler) error {

	// swagger docs

	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = app.config.apiUrl

	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Starting server on port", "port", app.config.port, "env", app.config.environment)
	return srv.ListenAndServe()
}
