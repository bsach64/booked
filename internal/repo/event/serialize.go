package eventrepo

import (
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
)

func ToEventDomainFromEventsRow(row db.GetEventsRow) *eventdom.Event {
	event := &eventdom.Event{
		Name:             row.Name,
		Address:          row.Address,
		Description:      row.Description,
		Time:             utils.GetUTCTime(row.Time.Time),
		ID:               row.ID.Bytes,
		SeatCount:        int(row.TotalTickets),
		AvailableTickets: int(row.AvailableTickets),
	}

	if row.Latitude.Valid && row.Longitude.Valid {
		event.Latitude = &row.Latitude.Float64
		event.Longitude = &row.Longitude.Float64
	}
	return event
}

func ToEventDomainFromNextEventsRow(row db.GetNextEventsRow) *eventdom.Event {
	event := &eventdom.Event{
		Name:             row.Name,
		Address:          row.Address,
		Description:      row.Description,
		Time:             utils.GetUTCTime(row.Time.Time),
		ID:               row.ID.Bytes,
		SeatCount:        int(row.TotalTickets),
		AvailableTickets: int(row.AvailableTickets),
	}

	if row.Latitude.Valid && row.Longitude.Valid {
		event.Latitude = &row.Latitude.Float64
		event.Longitude = &row.Longitude.Float64
	}
	return event
}
