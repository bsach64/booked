package eventrepo

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type impl struct {
	config  *utils.Config
	queries *db.Queries
	dbConn  *pgxpool.Pool
}

func (i *impl) GetEvents(ctx context.Context, limit int) ([]*eventdom.Event, int64, error) {
	var events []*eventdom.Event
	dbEvents, err := i.queries.GetEvents(ctx, int32(limit))
	if err != nil {
		return nil, 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read events", err)
	}

	if len(dbEvents) == 0 {
		return nil, 0, nil
	}

	lastUnixTime := utils.GetUTCUnixTime(dbEvents[len(dbEvents)-1].Time.Time)
	for _, dbEvent := range dbEvents {
		events = append(events, ToEventDomainFromEventsRow(dbEvent))
	}

	return events, lastUnixTime, nil
}

func (i *impl) GetNextEvents(ctx context.Context, unixTime int64, limit int) ([]*eventdom.Event, int64, error) {
	var events []*eventdom.Event
	getNextEventsParams := db.GetNextEventsParams{
		Time:  pgtype.Timestamp{Time: utils.GetUTCTimeFromUnix(unixTime), Valid: true},
		Limit: int32(limit),
	}
	dbEvents, err := i.queries.GetNextEvents(ctx, getNextEventsParams)
	if err != nil {
		return nil, 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read events", err)
	}

	if len(dbEvents) == 0 {
		return nil, 0, nil
	}

	lastUnixTime := utils.GetUTCUnixTime(dbEvents[len(dbEvents)-1].Time.Time)
	for _, dbEvent := range dbEvents {
		events = append(events, ToEventDomainFromNextEventsRow(dbEvent))
	}

	return events, lastUnixTime, nil
}

func (i *impl) CreateEvent(ctx context.Context, event *eventdom.Event) (uuid.UUID, error) {
	eventID := uuid.New()
	createEventParams := &db.CreateEventParams{
		ID:          pgtype.UUID{Bytes: eventID, Valid: true},
		Name:        event.Name,
		Address:     event.Address,
		Description: event.Description,
		Time:        pgtype.Timestamp{Time: utils.GetUTCTime(event.Time), Valid: true},
	}

	if event.Latitude == nil || event.Longitude == nil {
		createEventParams.Latitude.Valid = false
		createEventParams.Longitude.Valid = false
	} else {
		createEventParams.Latitude.Float64 = *event.Latitude
		createEventParams.Longitude.Float64 = *event.Longitude
		createEventParams.Latitude.Valid = true
		createEventParams.Longitude.Valid = true
	}

	err := i.queries.CreateEvent(ctx, *createEventParams)
	if err != nil {
		return [16]byte{}, errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not create event", err)
	}
	return eventID, nil
}

func (i *impl) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	_, err := i.queries.DeleteEvent(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return errordom.GetEventError(errordom.NO_EVENT_FOUND, "no event with given id", err)
	}

	if err != nil {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not delete event", err)
	}
	return nil
}

func (i *impl) UpdateEvent(ctx context.Context, Request *eventdom.UpdateEventRequest) error {
	tx, err := i.dbConn.Begin(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "could not start tx", err)
	}
	defer tx.Rollback(ctx)

	qtx := i.queries.WithTx(tx)
	eventUUID, err := uuid.Parse(Request.ID)
	if err != nil {
		return errordom.GetSystemError(errordom.INVALID_UUID, "invalid uuid for event", err)
	}

	event, err := qtx.GetEventByID(ctx, pgtype.UUID{Bytes: eventUUID, Valid: true})
	if err != nil {
		return err
	}

	Params := db.UpdateEventParams{
		ID: event.ID,
	}

	if Request.Name != nil {
		Params.Name = *Request.Name
	} else {
		Params.Name = event.Name
	}

	if Request.Description != nil {
		Params.Description = *Request.Description
	} else {
		Params.Description = event.Description
	}

	if Request.Address != nil {
		Params.Address = *Request.Address
	} else {
		Params.Address = event.Address
	}

	if Request.UnixTime != nil && *Request.UnixTime > 0 {
		Params.Time = pgtype.Timestamp{
			Time:  utils.GetUTCTimeFromUnix(*Request.UnixTime),
			Valid: true,
		}
	} else {
		Params.Time = event.Time
	}

	if Request.Latitude != nil && Request.Longitude != nil {
		Params.Latitude = pgtype.Float8{Float64: *Request.Latitude, Valid: true}
		Params.Longitude = pgtype.Float8{Float64: *Request.Longitude, Valid: true}
	} else {
		Params.Latitude = event.Latitude
		Params.Longitude = event.Longitude
	}

	if Request.SeatCount != nil {
		if *Request.SeatCount < int(event.TotalTickets) {
			return errordom.GetEventError(errordom.CANT_REDUCE_SEAT_COUNT, "cant reduce seat count", nil)
		}
		additionalSeats := *Request.SeatCount - int(event.TotalTickets)
		var tickets []db.CreateTicketsParams

		for range additionalSeats {
			tickets = append(tickets, db.CreateTicketsParams{
				ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
				EventID: event.ID,
			})
		}

		insertedCount, err := qtx.CreateTickets(ctx, tickets)
		if err != nil || insertedCount != int64(additionalSeats) {
			return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not create tickets", err)
		}
	}

	err = qtx.UpdateEvent(ctx, Params)
	if err != nil {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not update event", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "failed to commit tx", err)
	}
	return nil

}

func New(config *utils.Config, queries *db.Queries, pool *pgxpool.Pool) eventdom.Repository {
	return &impl{
		config:  config,
		queries: queries,
		dbConn:  pool,
	}
}
