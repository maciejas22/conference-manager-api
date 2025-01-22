package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int    `mapstructure:"SERVER_PORT"`
	GoEnv       string `mapstructure:"GO_ENV"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

var AppConfig *Config

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

const (
	EnvSharedFile = ".env"
	EnvProdFile   = ".env.prod"
	EnvDevFile    = ".env.dev"
)

func Init() {
	err := godotenv.Load(EnvSharedFile)
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to parse SERVER_PORT")
	}

	AppConfig = &Config{
		Port:        port,
		GoEnv:       os.Getenv("GO_ENV"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
