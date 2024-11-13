package main

import (
	"log"

	"github.com/puremike/social-go/pkg/env"
)

func main () {
	cfg := config {
		port: env.GetPort(),
	}

	app := &application{
        config: cfg,
    }
	mux := app.mount()
	log.Fatal(app.start(mux))
}