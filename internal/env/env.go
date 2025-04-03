package env

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Port            string `env:"PORT" envDefault:"9080"`
	TargetUrl       string `env:"TARGET_URL" envDefault:"http://127.0.0.1:3000"`
	AccessTokenName string `env:"ACCESS_TOKEN_NAME" envDefault:"access"`
	DisableLogs     int    `env:"DISABLE_LOGS" envDefault:"0"`
}

var cfg = Config{}
var loaded = false

func load() {
	loaded = true
	envFilePath := flag.String("env", ".env", "Config file to load")
	flag.Parse()
	err := godotenv.Load(*envFilePath)
	if err != nil {
		log.Println("Config file not found. System environment variables are used. For set config file use option env, example [docker-black-hole --env=.env.prod]")

	} else {
		log.Printf("Loading env variables from %s", *envFilePath)

	}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("Error parsing environment variables: %v\n", err)
	}
}

func GetEnv() Config {
	if !loaded {
		load()
	}
	return cfg
}
