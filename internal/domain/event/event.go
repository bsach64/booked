package eventdom

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID
	Name             string
	Time             time.Time
	Address          string
	Description      string
	SeatCount        int
	AvailableTickets int
	Latitude         *float64
	Longitude        *float64
}

type CreateEventRequest struct {
	Name        string   `json:"name"`
	UnixTime    int64    `json:"unix_time"`
	Address     string   `json:"address"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	SeatCount   int      `json:"seat_count"`
}

type EventResponse struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	UnixTime            int64    `json:"unix_time"`
	Address             string   `json:"address"`
	Description         string   `json:"description"`
	SeatCount           int      `json:"seat_count"`
	AvailableSeatsCount int      `json:"available_seats_count"`
	Latitude            *float64 `json:"latitude,omitempty"`
	Longitude           *float64 `json:"longitude,omitempty"`
}

type GetEventsResponse struct {
	NextTimeToFetch int64            `json:"next_time_to_fetch"`
	Events          []*EventResponse `json:"events"`
}

type Usecase interface {
	CreateEvent(ctx context.Context, eventRequest *CreateEventRequest) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvents(ctx context.Context, limit int, lastFetchedUnixTime *int64) (*GetEventsResponse, error)
}

type Repository interface {
	CreateEvent(ctx context.Context, event *Event) (uuid.UUID, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEvents(ctx context.Context, limit int) ([]*Event, int64, error)
	GetNextEvents(ctx context.Context, unixTime int64, limit int) ([]*Event, int64, error)
}
