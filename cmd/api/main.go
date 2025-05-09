package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/social-go/internal/auth"
	"github.com/puremike/social-go/internal/db"
	"github.com/puremike/social-go/internal/env"
	"github.com/puremike/social-go/internal/mailer"
	"github.com/puremike/social-go/internal/ratelimiter"
	"github.com/puremike/social-go/internal/store"
	"github.com/puremike/social-go/internal/store/cache"
	"go.uber.org/zap"
)

//	@title	Social_Go API

//	@description	This is an API for my Social_Go
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

var envData env.Config

const version = "1.1.0"

func main() {

	envData := env.GetPort()

	cfg := config{
		port: envData.Port,
		dbconfig: dbconfig{
			Addr:         envData.DB_URI,
			maxOpenConns: 10,
			maxIdleConns: 5, maxIdleTime: 15 * time.Minute, // 15 minutes /
		},
		environment: "development",
		apiUrl:      envData.SWAGGER_API_URL,
		mail: mailConfig{
			invitationExp: time.Hour * 24 * 3,
			fromEmail:     envData.FROM_EMAIL,
			// sendgrid: sendGridConfig{
			// 	apiKey: envData.SENDGRID_API_KEY,
			// },
			mailTrap: mailTrapConfig{
				apiKey: envData.MAILTRAP_API_KEY,
			},
		},
		frontEndURL: envData.FRONTEND_URL,
		auth: authConfig{
			username:    envData.AUTH_HEADER_USERNAME,
			password:    envData.AUTH_HEADER_PASSWORD,
			tokenSecret: envData.AUTH_TOKEN_SECRET,
			tokenExp:    time.Hour * 24 * 3, // 3 days
			iss:         "SocialGo",
			auds:        "SocialGo",
		},
		redisConfig: redisClientConfig{
			addr:    envData.REDIS_ADDR,
			pw:      envData.REDIS_PW,
			db:      0,
			enabled: false,
		},
		rateLimiter: rateLimiterConfig{
			requestsPerTimeFrame: 20,
			timeFrame:            5 * time.Second,
			enabled:              true,
		},
	}

	// Logger - using SugaredLogger
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewDB(cfg.dbconfig.Addr, cfg.dbconfig.maxOpenConns, cfg.dbconfig.maxIdleConns, cfg.dbconfig.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("Database connected successfully")

	// cache

	var rdb *redis.Client

	if cfg.redisConfig.enabled {
		rdb = cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.pw, cfg.redisConfig.db)
		defer rdb.Close()

		logger.Info("Redis client connected successfully")
	}

	//Rate limiter
	var rateLimit *ratelimiter.FixedWindowRateLimiter
	if cfg.rateLimiter.enabled {
		rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
			cfg.rateLimiter.requestsPerTimeFrame, cfg.rateLimiter.timeFrame,
		)

		rateLimit = rateLimiter

		logger.Info("Rate limiter is enabled")
	}

	// mailer := mailer.NewSendGridMailer(cfg.mail.fromEmail, cfg.mail.sendgrid.apiKey)
	mailer, err := mailer.NewMailTrapMailer(cfg.mail.fromEmail, cfg.mail.mailTrap.apiKey)
	if err != nil {
		logger.Errorw("Error: %v", err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.tokenSecret, cfg.auth.iss, cfg.auth.auds)

	str := store.NewStorage(db)
	cacheStorage := cache.NewRdbStorage(rdb)

	app := &application{
		config:        cfg,
		store:         str,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cacheStorage:  cacheStorage,
		rateLimiter:   rateLimit,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.start(mux))
}
