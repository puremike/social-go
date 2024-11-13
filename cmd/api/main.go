package main

import (
	"log"

	"github.com/puremike/social-go/internal/env"
	"github.com/puremike/social-go/internal/store"
)

func main () {
	cfg := config {
		port: env.GetPort(),
	}

	str := store.NewStorage(nil)

	app := &application{
        config: cfg, 
		store: str,
    }

	mux := app.mount()
	log.Fatal(app.start(mux))
}