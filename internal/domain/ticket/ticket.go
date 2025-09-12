package ticketdom

import (
	"context"

	eventdom "github.com/bsach64/booked/internal/domain/event"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/google/uuid"
)

type Status string

const (
	BOOKED     Status = "booked"
	AVAILIABLE Status = "availiable"
)

type Ticket struct {
	Status Status
	ID     uuid.UUID
	User   *userdom.User
	Event  *eventdom.Event
}

type Usecase interface {
	ReserveTicket(ctx context.Context, eventID uuid.UUID, userEmail string) error
	BookTicket(ctx context.Context, ticketID uuid.UUID) error
}

type Repository interface {
	CreateTickets(ctx context.Context, eventID uuid.UUID, count int) error
	ReserveTicket(ctx context.Context, ticketID uuid.UUID) error
	BookTicket(ctx context.Context, ticketID uuid.UUID) error
}
