package utils

import "github.com/caarlos0/env/v11"

type Config struct {
	DBUri string `env:"DB_URI" envDefault:"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
