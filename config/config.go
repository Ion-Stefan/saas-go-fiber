package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Define the Config
type Config struct {
	PublicHost                string
	Port                      string
	DBUser                    string
	DBPassword                string
	DBName                    string
	DBHost                    string
	DBPort                    string
	WebsiteURL                string
	LocalURL                  string
	BuildURL                  string
	JWTSecret                 string
	ClientID                  string
	ClientSecret              string
	RedirectURL               string
	LemonSqueezyWebhookSecret string
}

// Envs is the global variable that holds the configuration
var Envs = initConfig()

func initConfig() Config {
	// Load the environment variables from the .env file
	if err := godotenv.Load("../cmd/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Return the configuration
	return Config{
		PublicHost:                getEnv("PUBLIC_HOST"),
		Port:                      getEnv("PORT"),
		DBUser:                    getEnv("DB_USER"),
		DBPassword:                getEnv("DB_PASSWORD"),
		DBName:                    getEnv("DB_NAME"),
		DBHost:                    getEnv("DB_HOST"),
		DBPort:                    getEnv("DB_PORT"),
		WebsiteURL:                getEnv("WEBSITE_URL"),
		LocalURL:                  getEnv("LOCAL_URL"),
		BuildURL:                  getEnv("BUILD_URL"),
		JWTSecret:                 getEnv("JWT_SECRET"),
		ClientID:                  getEnv("CLIENT_ID"),
		ClientSecret:              getEnv("CLIENT_SECRET"),
		RedirectURL:               getEnv("REDIRECT_URL"),
		LemonSqueezyWebhookSecret: getEnv("LEMONSQUEEZY_WEBHOOK_SECRET"),
	}
}

// getEnv is a helper function that returns the value of the environment variable
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
