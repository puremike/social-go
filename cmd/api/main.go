package main

import (
	"log"

	"github.com/puremike/social-go/pkg/env"
	"github.com/puremike/social-go/pkg/store"
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