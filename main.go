package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bsach64/booked/internal/repo/sql/db"
	userrepo "github.com/bsach64/booked/internal/repo/user"
	useruc "github.com/bsach64/booked/internal/usecase/user"
	_ "github.com/lib/pq"

	"github.com/bsach64/booked/utils"
)

func main() {
	// minimum to test everything
	ctx := context.Background()
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

	dbConn := db.New(dbCon)
	userRepo := userrepo.New(config, dbConn)
	userUsecase := useruc.New(config, userRepo)

	// err = userUsecase.CreateUser(ctx, userdom.User{
	// 	Name:           "Bhavik Sachdev",
	// 	Email:          "b.sachdev1904@gmail.com",
	// 	HashedPassword: "1234",
	// 	Role:           userdom.USER,
	// })
	user, err := userUsecase.GetUserByEmail(ctx, "b.sachdev1904@gmail.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
}
