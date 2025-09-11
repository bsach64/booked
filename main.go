package main

import (
	"database/sql"
	"fmt"

	httpdelivery "github.com/bsach64/booked/delivery/http"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/internal/usecase"
	_ "github.com/lib/pq"

	"github.com/bsach64/booked/utils"
)

func main() {
	// minimum to test everything
	config, err := utils.GetConfig()
	if err != nil {
		fmt.Print(err)
		return
	}

	dbCon, err := sql.Open("postgres", config.DBUri)
	if err != nil {
		fmt.Print(err)
		return
	}

	dbQueries := db.New(dbCon)
	repositories := repo.New(config, dbQueries)
	usecases := usecase.New(config, repositories)
	server := httpdelivery.New(config, usecases, repositories)
	server.StartServer()
}
