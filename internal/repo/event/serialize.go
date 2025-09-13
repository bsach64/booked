package eventrepo

import (
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToEventDomainFromEventsRow(row db.GetEventsRow) *eventdom.Event {
	return &eventdom.Event{
		Name:             row.Name,
		Address:          row.Address,
		Description:      row.Description,
		Time:             utils.GetUTCTime(row.Time.Time),
		ID:               row.ID.Bytes,
		SeatCount:        int(row.TotalTickets),
		AvailableTickets: int(row.AvailableTickets),
		Latitude:         getPtrIfValid(row.Latitude),
		Longitude:        getPtrIfValid(row.Longitude),
	}
}

func ToEventDomainFromNextEventsRow(row db.GetNextEventsRow) *eventdom.Event {
	return &eventdom.Event{
		Name:             row.Name,
		Address:          row.Address,
		Description:      row.Description,
		Time:             utils.GetUTCTime(row.Time.Time),
		ID:               row.ID.Bytes,
		SeatCount:        int(row.TotalTickets),
		AvailableTickets: int(row.AvailableTickets),
		Latitude:         getPtrIfValid(row.Latitude),
		Longitude:        getPtrIfValid(row.Longitude),
	}
}

func getPtrIfValid(f pgtype.Float8) *float64 {
	if f.Valid {
		return &f.Float64
	}
	return nil
}
