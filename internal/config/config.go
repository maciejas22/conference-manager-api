package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                     int      `mapstructure:"SERVER_PORT"`
	GoEnv                    string   `mapstructure:"GO_ENV"`
	ServerPort               string   `mapstructure:"SERVER_PORT"`
	CorsAllowedOrigins       []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	JWTSecret                string   `mapstructure:"JWT_SECRET"`
	DatabaseURL              string   `mapstructure:"DATABASE_URL"`
	S3Region                 string   `mapstructure:"S3_REGION"`
	S3Endpoint               string   `mapstructure:"S3_ENDPOINT"`
	S3AccessKeyID            string   `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey        string   `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3UrlLifetimeInHours     int      `mapstructure:"S3_URL_LIFETIME_IN_HOURS"`
	S3BucketsConferenceFiles string   `mapstructure:"S3_BUCKETS_CONFERENCE_FILES"`
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
}
