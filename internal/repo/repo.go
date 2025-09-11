package repo

import (
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo/sql/db"
	userrepo "github.com/bsach64/booked/internal/repo/user"
	"github.com/bsach64/booked/utils"
)

type Repositories struct {
	Config *utils.Config
	db     *db.Queries
	User   userdom.Repository
}

func New(config *utils.Config, dbConn *db.Queries) Repositories {
	return Repositories{
		Config: config,
		db:     dbConn,
		User:   userrepo.New(config, dbConn),
	}
}
