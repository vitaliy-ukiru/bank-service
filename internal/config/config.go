package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Database string `env:"DB_DATABASE"`
	Password string `env:"DB_PASSWORD"`
}

type WebServerConfig struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

type Env string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

type Config struct {
	Database DatabaseConfig
	Server   WebServerConfig
	Env      Env `env:"APP_ENV" env-default:"dev"`
}

var cfg Config
var cfgLoaded bool

func LoadConfig(envFile string) error {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil && os.IsNotExist(err) {
			return err
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return err
	}
	cfgLoaded = true
	return nil
}

func Get() Config {
	if !cfgLoaded {
		panic("config not loaded")
	}
	return cfg
}
