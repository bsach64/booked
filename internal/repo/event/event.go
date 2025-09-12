package eventrepo

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type impl struct {
	config *utils.Config
	dbConn *db.Queries
}

func (i *impl) CreateEvent(ctx context.Context, event *eventdom.Event) (uuid.UUID, error) {
	eventID := uuid.New()
	createEventParams := &db.CreateEventParams{
		ID:          pgtype.UUID{Bytes: eventID, Valid: true},
		Name:        event.Name,
		Address:     event.Address,
		Description: event.Description,
		Time:        pgtype.Timestamp{Time: event.Time, Valid: true},
	}

	if event.Latitude == 0 || event.Longitude == 0 {
		createEventParams.Latitude.Valid = false
		createEventParams.Longitude.Valid = false
	} else {
		createEventParams.Latitude.Float64 = event.Latitude
		createEventParams.Longitude.Float64 = event.Longitude
		createEventParams.Latitude.Valid = true
		createEventParams.Longitude.Valid = true
	}

	err := i.dbConn.CreateEvent(ctx, *createEventParams)
	if err != nil {
		return [16]byte{}, errordom.GetDBError(errordom.DB_WRITE_ERROR, "", err)
	}
	return eventID, nil
}

func New(config *utils.Config, db *db.Queries) eventdom.Repository {
	return &impl{
		config: config,
		dbConn: db,
	}
}
