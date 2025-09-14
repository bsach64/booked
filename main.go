package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("could not get config err=%v\n", err)
	}

	pool, err := pgxpool.New(ctx, config.DBUri)
	if err != nil {
		log.Fatalf("could not establish connection to db err=%v\n", err)
	}

	valkeyURL, err := valkey.ParseURL(config.ValkeyURL)
	if err != nil {
		log.Fatalf("incorrect valkey url err=%v", err)
	}

	valkeyClient, err := valkey.NewClient(valkeyURL)
	if err != nil {
		log.Fatalf("could not establish connection to valkey err=%v\n", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	dbQueries := db.New(pool)
	repositories := repo.New(config, dbQueries, pool, valkeyClient)
	usecases := usecase.New(config, repositories)
	server := httpdelivery.New(config, usecases, repositories)

	go func() {
		slog.Info("starting http server", "url", config.ServerURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped unexpectedly", "err", err)
			cancel() // cancel context if server fails
		}
	}()

	// Run background notify job
	go BackgroundNotifyJob(ctx, usecases, config)

	<-stop
	slog.Info("shutting down...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("error while shutting down server", "err", err)
	}
}

func BackgroundNotifyJob(ctx context.Context, uc usecase.Usecase, config *utils.Config) {
	if err := uc.WaitlistUC.NotifyUsers(ctx); err != nil {
		slog.Error("error while running notify job", "err", err)
		return
	}

	timeInSeconds, err := strconv.Atoi(config.NotifyWaitlistInSeconds)
	if err != nil {
		slog.Error("could not parse notify time", "err", err)
		return
	}

	ticker := time.NewTicker(time.Duration(timeInSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("background notify job shutting down...")
			return
		case <-ticker.C:
			if err := uc.WaitlistUC.NotifyUsers(ctx); err != nil {
				slog.Error("error while running notify job", "err", err)
			}
		}
	}
}
