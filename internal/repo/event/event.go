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
)

type impl struct {
	config *utils.Config
	dbConn *db.Queries
}

func (i *impl) GetEvents(ctx context.Context, limit int) ([]*eventdom.Event, int64, error) {
	var events []*eventdom.Event
	dbEvents, err := i.dbConn.GetEvents(ctx, int32(limit))
	if err != nil {
		return nil, 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read events", err)
	}

	if len(dbEvents) == 0 {
		return nil, 0, nil
	}

	lastUnixTime := utils.GetUTCUnixTime(dbEvents[len(dbEvents)-1].Time.Time)
	for _, dbEvent := range dbEvents {
		events = append(events, ToEventDomain(dbEvent))
	}

	return events, lastUnixTime, nil
}

func (i *impl) GetNextEvents(ctx context.Context, unixTime int64, limit int) ([]*eventdom.Event, int64, error) {
	var events []*eventdom.Event
	getNextEventsParams := db.GetNextEventsParams{
		Time:  pgtype.Timestamp{Time: utils.GetUTCTimeFromUnix(unixTime), Valid: true},
		Limit: int32(limit),
	}
	dbEvents, err := i.dbConn.GetNextEvents(ctx, getNextEventsParams)
	if err != nil {
		return nil, 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read events", err)
	}

	if len(dbEvents) == 0 {
		return nil, 0, nil
	}

	lastUnixTime := utils.GetUTCUnixTime(dbEvents[len(dbEvents)-1].Time.Time)
	for _, dbEvent := range dbEvents {
		events = append(events, ToEventDomain(dbEvent))
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
		SeatCount:   int64(event.SeatCount),
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

	err := i.dbConn.CreateEvent(ctx, *createEventParams)
	if err != nil {
		return [16]byte{}, errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not create event", err)
	}
	return eventID, nil
}

func (i *impl) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	_, err := i.dbConn.DeleteEvent(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return errordom.GetEventError(errordom.NO_EVENT_FOUND, "no event with given id", err)
	}
	return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not delete event", err)
}

func New(config *utils.Config, db *db.Queries) eventdom.Repository {
	return &impl{
		config: config,
		dbConn: db,
	}
}
