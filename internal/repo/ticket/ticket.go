package ticketrepo

import (
	"context"
	"errors"
	"time"

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

func (i *impl) GetAvailiableTickets(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	dbTicketIDs, err := i.dbConn.GetAvailableTickets(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err != nil {
		return nil, err
	}
	var ticketIDs uuid.UUIDs
	for _, dbID := range dbTicketIDs {
		id := uuid.UUID(dbID.Bytes)
		_, err = i.valkeyClient.Do(ctx, i.valkeyClient.B().Get().Key(id.String()).Build()).ToString()
		if err == nil {
			continue
		}
		if !errors.Is(err, valkey.Nil) {
			return nil, err
		}

		ticketIDs = append(ticketIDs, id)
	}
	return ticketIDs, nil
}

func (i *impl) ReserveTickets(ctx context.Context, userID uuid.UUID, ticketIDs []uuid.UUID) error {
	for _, ticket := range ticketIDs {
		err := i.valkeyClient.Do(ctx, i.valkeyClient.B().Set().Key(ticket.String()).Value(userID.String()).Ex(10*time.Minute).Build()).Error()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *impl) BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []uuid.UUID) error {
	dbTicketIDs := []pgtype.UUID{}
	dbUserID := pgtype.UUID{Bytes: userID, Valid: true}
	for _, ticket := range ticketIDs {
		userIDForTicket, err := i.valkeyClient.Do(ctx, i.valkeyClient.B().Get().Key(ticket.String()).Build()).ToString()
		if err != nil {
			return err
		}

		if userIDForTicket != userID.String() {
			// create err
			return err
		}

		dbTicketIDs = append(dbTicketIDs, pgtype.UUID{Bytes: ticket, Valid: true})
	}

	err := i.dbConn.BookTickets(ctx, db.BookTicketsParams{UserID: dbUserID, Column2: dbTicketIDs})
	if err != nil {
		return err
	}
	return nil
}

func New(config *utils.Config, db *db.Queries, valkeyClient valkey.Client) ticketdom.Repository {
	return &impl{
		config:       config,
		dbConn:       db,
		valkeyClient: valkeyClient,
	}
}
