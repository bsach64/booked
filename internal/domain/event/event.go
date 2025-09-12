package eventdom

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Name        string
	Time        time.Time
	Address     string
	Description string
	Latitude    float64
	Longitude   float64
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

type Usecase interface {
	CreateEvent(ctx context.Context, eventRequest *CreateEventRequest) error
}

type Repository interface {
	CreateEvent(ctx context.Context, event *Event) (uuid.UUID, error)
}
