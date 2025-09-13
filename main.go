package main

import (
	"context"
	"log"

	httpdelivery "github.com/bsach64/booked/delivery/http"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/valkey-io/valkey-go"

	"github.com/bsach64/booked/utils"
)

func main() {
	ctx := context.Background()
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("could not get config err=%v\n", err)
		return
	}

	pool, err := pgxpool.New(ctx, config.DBUri)
	if err != nil {
		log.Fatalf("could not establish connection to db err=%v\n", err)
		return
	}

	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{config.ValkeyURL},
	})
	if err != nil {
		log.Fatalf("could not establish connection to valkey err=%v\n", err)
	}

	dbQueries := db.New(pool)
	repositories := repo.New(config, dbQueries, pool, valkeyClient)
	usecases := usecase.New(config, repositories)
	server := httpdelivery.New(config, usecases, repositories)
	server.StartServer()
}
