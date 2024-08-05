package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                     int    `mapstructure:"SERVER_PORT"`
	AppEnv                   string `mapstructure:"APP_ENV"`
	ServerPort               string `mapstructure:"SERVER_PORT"`
	JWTSecret                string `mapstructure:"JWT_SECRET"`
	DatabaseURL              string `mapstructure:"DATABASE_URL"`
	S3Region                 string `mapstructure:"S3_REGION"`
	S3Endpoint               string `mapstructure:"S3_ENDPOINT"`
	S3AccessKeyID            string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey        string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3UrlLifetimeInHours     int    `mapstructure:"S3_URL_LIFETIME_IN_HOURS"`
	S3BucketsConferenceFiles string `mapstructure:"S3_BUCKETS_CONFERENCE_FILES"`
}

var AppConfig *Config

func LoadConfig() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	AppConfig = &Config{}
	err = viper.Unmarshal(AppConfig)
	if err != nil {
		log.Fatal("Config can't be loaded: ", err)
	}

	if AppConfig.AppEnv == "development" {
		log.Println("The App is running in development env")
	}
}
