package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bsach64/booked/internal/repo/sql/db"
	_ "github.com/lib/pq"

	"github.com/bsach64/booked/utils"
)

func main() {
	ctx := context.Background()
	config, err := utils.GetConfig()
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("dbURI ", config.DBUri)
	dbCon, err := sql.Open("postgres", config.DBUri)
	if err != nil {
		fmt.Print(err)
		return
	}

	queries := db.New(dbCon)

	fmt.Println(queries.GetAllUsers(ctx))
}
