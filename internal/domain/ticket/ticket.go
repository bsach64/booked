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

type Usecase interface {
	ReserveTickets(ctx context.Context, reserveTickets *ReserveTicketRequest) (*ReserveTicketsResponse, error)
	BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []string) error
}

type Repository interface {
	CreateTickets(ctx context.Context, eventID uuid.UUID, count int) error
	ReserveTickets(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, count int) ([]uuid.UUID, error)
	BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []uuid.UUID) error
}
