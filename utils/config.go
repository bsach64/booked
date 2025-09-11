package utils

import "github.com/caarlos0/env/v11"

type Config struct {
	DBUri     string `env:"DB_URI" envDefault:"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	JwtSecret string `env:"JWT_SECRET" envDefault:"verysecret"`
	ServerURL string `env:"SERVER_URL" envDefault:"localhost:8080"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
