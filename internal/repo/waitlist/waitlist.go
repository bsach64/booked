package waitlistrepo

import (
	"context"
	"errors"

	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	waitlistdom "github.com/bsach64/booked/internal/domain/waitlist"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
)

type impl struct {
	config       *utils.Config
	queries      *db.Queries
	dbConn       *pgxpool.Pool
	valkeyClient valkey.Client
}

func (i *impl) AddToWaitlist(ctx context.Context, user *userdom.User, eventID uuid.UUID, count int) error {
	tx, err := i.dbConn.Begin(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "could not start tx", err)
	}
	defer tx.Rollback(ctx)

	qtx := i.queries.WithTx(tx)

	_, err = qtx.GetWaitlistEntry(ctx, db.GetWaitlistEntryParams{
		UserID:  pgtype.UUID{Bytes: user.ID, Valid: true},
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return errordom.GetDBError(errordom.DB_READ_ERROR, "could not read from db", err)
	}

	if err == nil {
		return errordom.GetDBError(errordom.ALREADY_IN_WAILIST, "already in waitlist", err)
	}

	if count <= 0 {
		return errordom.GetWaitlistError(errordom.INVALID_WAITLIST_COUNT, "can't waitlist for 0 tickets", nil)
	}

	err = qtx.AddToWaitlist(ctx, db.AddToWaitlistParams{
		ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:  pgtype.UUID{Bytes: user.ID, Valid: true},
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		Count:   int32(count),
	})

	err = tx.Commit(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "failed to commit tx", err)
	}
	return nil
}

func (i *impl) UpdateNotificationStatus(ctx context.Context, ID uuid.UUIDs, status waitlistdom.NotificationStatus) error {
	params := db.UpdateWaitlistStatusParams{
		Status: db.NotificationStatus(status),
	}

	for _, id := range ID {
		params.Column2 = append(params.Column2, pgtype.UUID{Bytes: id, Valid: true})
	}

	err := i.queries.UpdateWaitlistStatus(ctx, params)
	if err != nil {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not write from db", err)
	}
	return nil
}

func (i *impl) GetWaitlistNotifications(ctx context.Context) ([]*waitlistdom.WaitlistNotification, error) {
	dbEntries, err := i.queries.GetWaitlistNotificationDetails(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read from db", err)
	}

	entries := []*waitlistdom.WaitlistNotification{}

	for _, dbE := range dbEntries {
		entry := &waitlistdom.WaitlistNotification{
			UserName:   dbE.UserName,
			EventName:  dbE.EventName,
			UserEmail:  dbE.UserEmail,
			EventTime:  dbE.EventTime.Time,
			SeatCount:  int(dbE.Count),
			EventID:    dbE.EventID.Bytes,
			WaitlistID: dbE.WaitlistID.Bytes,
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func New(config *utils.Config, queries *db.Queries, dbConn *pgxpool.Pool, valkeyClient valkey.Client) waitlistdom.Repository {
	return &impl{
		config:       config,
		queries:      queries,
		dbConn:       dbConn,
		valkeyClient: valkeyClient,
	}
}
