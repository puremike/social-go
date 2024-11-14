package main

import (
	"log"
	"time"

	"github.com/puremike/social-go/internal/db"
	"github.com/puremike/social-go/internal/env"
	"github.com/puremike/social-go/internal/store"
)

func main () {

	envData := env.GetPort()

	cfg := config {
		port: envData.Port,
		dbconfig: dbconfig {
			Addr:     envData.DB_URI,
            maxOpenConns: 10,
            maxIdleConns: 5,
            maxIdleTime:  15 * time.Minute, // 15 minutes
		},
	}
	db, err := db.NewDB(cfg.dbconfig.Addr, cfg.dbconfig.maxOpenConns, cfg.dbconfig.maxIdleConns, cfg.dbconfig.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	str := store.NewStorage(db)

	app := &application{
        config: cfg, 
		store: str,
    }

	mux := app.mount()
	log.Fatal(app.start(mux))
}