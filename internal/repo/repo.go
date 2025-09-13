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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
)

type Repositories struct {
	Config       *utils.Config
	queries      *db.Queries
	valkeyClient valkey.Client
	User         userdom.Repository
	Event        eventdom.Repository
	Ticket       ticketdom.Repository
}

func New(config *utils.Config, queries *db.Queries, dbConn *pgxpool.Pool, valkeyClient valkey.Client) Repositories {
	return Repositories{
		Config:       config,
		queries:      queries,
		valkeyClient: valkeyClient,
		User:         userrepo.New(config, queries),
		Event:        eventrepo.New(config, queries),
		Ticket:       ticketrepo.New(config, queries, dbConn, valkeyClient),
	}
}
