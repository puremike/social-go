package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, DB_URI string
}

func GetPort() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal()
	}

	port := os.Getenv(":PORT")
	if port == "" {
		port = ":5100"
	}

	db_uri := os.Getenv("DB_URI")
	if db_uri == "" {
        log.Fatal("DB_URI is not set")
    }

	return Config{
		Port: port,
        DB_URI: db_uri,
	}
}