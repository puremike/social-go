package main

import (
	"log"
)

func main () {
	cfg := config {
		port: ":5100",
	}

	app := &application{
        config: cfg,
    }
	mux := app.mount()
	log.Fatal(app.start(mux))
}