package repo

import (
	eventdom "github.com/bsach64/booked/internal/domain/event"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	userdom "github.com/bsach64/booked/internal/domain/user"
	eventrepo "github.com/bsach64/booked/internal/repo/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	ticketrepo "github.com/bsach64/booked/internal/repo/ticket"
	userrepo "github.com/bsach64/booked/internal/repo/user"
	"github.com/bsach64/booked/utils"
)

type Repositories struct {
	Config *utils.Config
	db     *db.Queries
	User   userdom.Repository
	Event  eventdom.Repository
	Ticket ticketdom.Repository
}

func New(config *utils.Config, dbConn *db.Queries) Repositories {
	return Repositories{
		Config: config,
		db:     dbConn,
		User:   userrepo.New(config, dbConn),
		Event:  eventrepo.New(config, dbConn),
		Ticket: ticketrepo.New(config, dbConn),
	}
}
