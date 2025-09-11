package utils

type Config struct {
	DBUri string `env:"DB_URI" default:"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}
