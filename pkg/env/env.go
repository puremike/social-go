package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


func GetPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5100"
	}

	return ":" + port
}