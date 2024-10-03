package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port                      int      `mapstructure:"SERVER_PORT"`
	GoEnv                     string   `mapstructure:"GO_ENV"`
	ServerPort                string   `mapstructure:"SERVER_PORT"`
	CorsAllowedOrigins        []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	DatabaseURL               string   `mapstructure:"DATABASE_URL"`
	AWSRegion                 string   `mapstructure:"AWS_REGION"`
	AWSEndpoint               string   `mapstructure:"AWS_ENDPOINT_URL_S3"`
	AWSAccessKeyId            string   `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey        string   `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSBucketsConferenceFiles string   `mapstructure:"AWS_BUCKETS_CONFERENCE_FILES"`
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

func setAWSConfig() {
	os.Setenv("AWS_REGION", AppConfig.AWSRegion)
	os.Setenv("AWS_ENDPOINT_URL_S3", AppConfig.AWSEndpoint)
	os.Setenv("AWS_ACCESS_KEY_ID", AppConfig.AWSAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", AppConfig.AWSSecretAccessKey)
}

func LoadConfig() {
	viper.AutomaticEnv()
	viper.SetConfigFile(EnvSharedFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cant find the file %s: %v", EnvSharedFile, err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("config cant be loaded: %v", err)
	}

	env := viper.GetString("GO_ENV")
	if env != EnvProd {
		env = EnvDev
	}

	envFile := EnvDevFile
	if env == EnvProd {
		envFile = EnvProdFile
	}

	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("cant find the file %s: %v", envFile, err)
	}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("config cant be loaded: %v", err)
	}

	log.Printf("app is running in %s mode\n", env)
	log.Printf("cors allowed origins: %v\n", AppConfig.CorsAllowedOrigins)
	setAWSConfig()
}
