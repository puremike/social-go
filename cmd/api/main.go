package main

import (
	"time"

	"github.com/puremike/social-go/internal/db"
	"github.com/puremike/social-go/internal/env"
	"github.com/puremike/social-go/internal/store"
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

	str := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  str,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.start(mux))
}
