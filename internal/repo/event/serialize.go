package eventrepo

import (
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
)

func ToEventDomain(dbEvent db.Event) *eventdom.Event {
	event := eventdom.Event{
		Name:        dbEvent.Name,
		Address:     dbEvent.Address,
		Description: dbEvent.Description,
		SeatCount:   int(dbEvent.SeatCount),
		Time:        utils.GetUTCTime(dbEvent.Time.Time),
		ID:          dbEvent.ID.Bytes,
	}

	if dbEvent.Latitude.Valid && dbEvent.Longitude.Valid {
		event.Latitude = &dbEvent.Latitude.Float64
		event.Longitude = &dbEvent.Latitude.Float64
	}

	return &event
}
