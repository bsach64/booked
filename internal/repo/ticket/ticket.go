package ticketrepo

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/valkey-io/valkey-go"
)

type impl struct {
	config       *utils.Config
	dbConn       *db.Queries
	valkeyClient valkey.Client
}

func (i *impl) CreateTickets(ctx context.Context, eventID uuid.UUID, count int) error {
	var tickets []db.CreateTicketsParams

	for range count {
		tickets = append(tickets, db.CreateTicketsParams{
			ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
			EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		})
	}

	insertedCount, err := i.dbConn.CreateTickets(ctx, tickets)
	if err != nil || insertedCount != int64(count) {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "", err)
	}
	return nil
}

func (i *impl) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
	return nil
}

func (i *impl) BookTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
	return nil
}

func New(config *utils.Config, db *db.Queries, valkeyClient valkey.Client) ticketdom.Repository {
	return &impl{
		config:       config,
		dbConn:       db,
		valkeyClient: valkeyClient,
	}
}
