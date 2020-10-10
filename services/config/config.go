package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APPPort string `envconfig:"APP_PORT"`

	DBHost               string `envconfig:"DB_HOST"`
	DBPort               string `envconfig:"DB_PORT"`
	DBName               string `envconfig:"DB_NAME"`
	DBUser               string `envconfig:"DB_USER"`
	DBPass               string `envconfig:"DB_PASS"`
	DBConnection         string `envconfig:"DB_CONNECTION"`
	DBDialTimeout        int    `envconfig:"DB_DIAL_TIMEOUT"`
	DBReadTimeout        int    `envconfig:"DB_READ_TIMEOUT"`
	DBWriteTimeout       int    `envconfig:"DB_WRITE_TIMEOUT"`
	DBMaxOpenConnections int    `envconfig:"DB_MAX_OPEN_CONNECTIONS"`
	DBMaxIdleConnections int    `envconfig:"DB_MAX_IDLE_CONNECTIONS"`
}

var (
	conf Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		godotenv.Load()
		err := envconfig.Process("", &conf)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	return &conf
}
