package main

import (
	"log"

	"github.com/puremike/social-go/internal/db"
	"github.com/puremike/social-go/internal/env"
	"github.com/puremike/social-go/internal/store"
)

func main() {

	port := env.GetPort()

	conn, err := db.NewDB(port.DB_URI, 3, 3, 15)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close();
	store := store.NewStorage(conn)

	db.Seed(store)
}