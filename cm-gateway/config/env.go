package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type BucketNames struct {
	ConferenceFiles string
}

type Config struct {
	Port                  int
	GoEnv                 string
	StripeSecretKey       string
	InfoServiceAddr       string
	AuthServiceAddr       string
	ConferenceServiceAddr string
	AWSRegion             string
	AWSEndpoint           string
	AWSAccessKeyId        string
	AWSSecretAccessKey    string
	Buckets               BucketNames
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
		log.Fatal("failed to parse SERVER_PORT")
	}

	goEnv := os.Getenv("GO_ENV")
	if goEnv == EnvProd {
		goEnv = EnvProd
	} else {
		goEnv = EnvDev
	}
	log.Println("app running in", goEnv, "mode")

	AppConfig = &Config{
		Port:                  port,
		GoEnv:                 goEnv,
		StripeSecretKey:       os.Getenv("STRIPE_SECRET_KEY"),
		InfoServiceAddr:       os.Getenv("INFO_SERVICE_ADDR"),
		AuthServiceAddr:       os.Getenv("AUTH_SERVICE_ADDR"),
		ConferenceServiceAddr: os.Getenv("CONFERENCE_SERVICE_ADDR"),
		AWSRegion:             os.Getenv("AWS_REGION"),
		AWSEndpoint:           os.Getenv("AWS_ENDPOINT_URL_S3"),
		AWSAccessKeyId:        os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Buckets: BucketNames{
			ConferenceFiles: os.Getenv("AWS_BUCKETS_CONFERENCES_FILES"),
		},
	}
}
