package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/puremike/social-go/docs"
	"github.com/puremike/social-go/internal/auth"
	"github.com/puremike/social-go/internal/mailer"
	"github.com/puremike/social-go/internal/ratelimiter"
	"github.com/puremike/social-go/internal/store"
	"github.com/puremike/social-go/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         *store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  *cache.Storage
	rateLimiter   *ratelimiter.FixedWindowRateLimiter
}

type config struct {
	port        string
	dbconfig    dbconfig
	environment string
	apiUrl      string
	mail        mailConfig
	frontEndURL string
	auth        authConfig
	redisConfig redisClientConfig
	rateLimiter rateLimiterConfig
}

type redisClientConfig struct {
	addr, pw string
	db       int
	enabled  bool
}
type authConfig struct {
	username    string
	password    string
	tokenSecret string
	iss         string
	auds        string
	tokenExp    time.Duration
}

type rateLimiterConfig struct {
	requestsPerTimeFrame int
	timeFrame            time.Duration
	enabled              bool
}
type mailConfig struct {
	invitationExp time.Duration
	fromEmail     string
	mailTrap      mailTrapConfig
	// sendgrid      sendGridConfig
}

//	type sendGridConfig struct {
//		apiKey string
//	}
type mailTrapConfig struct {
	apiKey string
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
	r.Use(app.rateLimiterMiddleware)

	// Set a timeout value on the request context, ctx, that will signal through ctx.Done() that the request has timed out and further processing should be canceled.

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.basicAuthMiddleware).Get("/health", app.health)

		docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docURL), //The url pointing to API definition
		))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPost)
			r.Get("/", app.getAllPosts)
			// r.Delete("/", app.deleteAllPosts)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postContextMiddleWare)
				r.Get("/", app.getPostById)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostByID))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePost))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Post("/", app.createUser)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserByID)
				r.Delete("/", app.deleteUserByIDHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/{id}/feeds", app.getUserFeedsHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.logger.Infow("shutdown signal", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Starting server on port:", "port", app.config.port, "env", app.config.environment)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return err
	}
	return nil
}
