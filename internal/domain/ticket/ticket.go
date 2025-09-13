package ticketdom

import (
	"context"

	eventdom "github.com/bsach64/booked/internal/domain/event"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/google/uuid"
)

type Status string

const (
	BOOKED    Status = "booked"
	AVAILABLE Status = "available"
)

type Ticket struct {
	Status Status
	ID     uuid.UUID
	User   *userdom.User
	Event  *eventdom.Event
}

type ReserveTicketRequest struct {
	UserID  uuid.UUID `json:"-"`
	EventID string    `json:"event_id"`
	Count   int       `json:"count"`
}

type ReserveTicketsResponse struct {
	TicketIDs []string `json:"ticket_ids"`
}

type CancelTicketRequest struct {
	EventID string `json:"event_id"`
	Count   int    `json:"count"`
}

type PastBookingsResponse struct {
	EventID          string   `json:"event_id"`
	EventName        string   `json:"event_name"`
	EventUnixTime    int64    `json:"event_unix_time"`
	EventAddress     string   `json:"event_address"`
	EventDescription string   `json:"event_description"`
	EventLatitude    *float64 `json:"event_latitude,omitempty"`
	EventLongitude   *float64 `json:"event_longitude,omitempty"`
	NumberOfTickets  int      `json:"number_of_tickets"`
	BookingUnixTime  int64    `json:"booking_unix_time"`
}

type Usecase interface {
	ReserveTickets(ctx context.Context, reserveTickets *ReserveTicketRequest) (*ReserveTicketsResponse, error)
	BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []string) error
	GetPastBookings(ctx context.Context, user *userdom.User) ([]*PastBookingsResponse, error)
	CancelTickets(ctx context.Context, user *userdom.User, cancelTicketRequest *CancelTicketRequest) error
}

type Repository interface {
	CreateTickets(ctx context.Context, eventID uuid.UUID, count int) error
	ReserveTickets(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, count int) ([]uuid.UUID, error)
	BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []uuid.UUID) error
	GetPastBookings(ctx context.Context, userID uuid.UUID) ([]*PastBookingsResponse, error)
	CancelTickets(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, count int) error
	GetReservedTickets(ctx context.Context, eventID uuid.UUID) (int, error)
}
