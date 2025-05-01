package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, DB_URI, SWAGGER_API_URL, SENDGRID_API_KEY, MAILTRAP_API_KEY, FROM_EMAIL, FRONTEND_URL, AUTH_HEADER_USERNAME, AUTH_HEADER_PASSWORD, AUTH_TOKEN_SECRET, REDIS_ADDR, REDIS_PW string
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

	frontendURl := os.Getenv(":FRONTEND_URL")
	if frontendURl == "" {
		frontendURl = "http://localhost:4022"
	}

	db_uri := os.Getenv("DB_URI")
	if db_uri == "" {
		log.Fatal("DB_URI is not set")
	}

	api_url := os.Getenv("SWAGGER_API_URL")
	if api_url == "" {
		api_url = "localhost:5100"
	}

	// sendgrid_api_key := os.Getenv("SENDGRID_API_KEY")
	// if sendgrid_api_key == "" {
	// 	log.Fatal("SENDGRID_API_KEY NOT SET")
	// }

	mailTrap_api_key := os.Getenv("MAILTRAP_API_KEY")
	if mailTrap_api_key == "" {
		log.Fatal("MAILTRAP_API_KEY NOT SET")
	}

	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		log.Fatal("FROM_EMAIl NOT SET")
	}

	authHeader_user := os.Getenv("AUTH_HEADER_USERNAME")
	if authHeader_user == "" {
		log.Fatal("AUTH_HEADER_USERNAME NOT SET")
	}

	authHeader_pass := os.Getenv("AUTH_HEADER_PASSWORD")
	if authHeader_pass == "" {
		log.Fatal("AUTH_HEADER_PASSWORD NOT SET")
	}

	authTokenSecret := os.Getenv("AUTH_TOKEN_SECRET")
	if authTokenSecret == "" {
		log.Fatal("AUTH_TOKEN_SECRET NOT SET")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPw := os.Getenv("REDIS_PW")
	if redisPw == "" {
		redisPw = ""
	}

	return Config{
		Port:            port,
		DB_URI:          db_uri,
		SWAGGER_API_URL: api_url,
		// SENDGRID_API_KEY: sendgrid_api_key,
		MAILTRAP_API_KEY:     mailTrap_api_key,
		FROM_EMAIL:           fromEmail,
		FRONTEND_URL:         frontendURl,
		AUTH_HEADER_USERNAME: authHeader_user,
		AUTH_HEADER_PASSWORD: authHeader_pass,
		AUTH_TOKEN_SECRET:    authTokenSecret,
		REDIS_ADDR:           redisAddr,
		REDIS_PW:             redisPw,
	}
}
