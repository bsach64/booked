package ticketuc

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
)

type impl struct {
	config       *utils.Config
	repositories repo.Repositories
}

func (i *impl) ReserveTickets(ctx context.Context, reserveTickets *ticketdom.ReserveTicketRequest) (*ticketdom.ReserveTicketsResponse, error) {
	eventID, err := uuid.Parse(reserveTickets.EventID)
	if err != nil {
		return nil, errordom.GetSystemError(errordom.INVALID_UUID, "invalid uuid", err)
	}

	if reserveTickets.Count <= 0 {
		return nil, errordom.GetEventError(errordom.INVALID_SEAT_COUNT, "seat count is less than 0", err)
	}

	ticketIDs, err := i.repositories.Ticket.ReserveTickets(ctx, reserveTickets.UserID, eventID, reserveTickets.Count)
	if err != nil {
		return nil, err
	}

	// ticketIDs, err := i.repositories.Ticket.GetAvailiableTickets(ctx, eventID)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if len(ticketIDs) < reserveTickets.Count {
	// 	return nil, errordom.GetEventError(errordom.TOO_FEW_TICKETS, "not enough tickets", err)
	// }
	//
	// err = i.repositories.Ticket.ReserveTickets(ctx, reserveTickets.UserID, ticketIDs[:reserveTickets.Count])
	// if err != nil {
	// 	return nil, err
	// }
	//
	var ticketIDStrs []string
	for _, id := range ticketIDs[:reserveTickets.Count] {
		ticketIDStrs = append(ticketIDStrs, id.String())
	}

	return &ticketdom.ReserveTicketsResponse{
		TicketIDs: ticketIDStrs,
	}, nil
}

func (i *impl) BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []string) error {
	ids := uuid.UUIDs{}

	for _, idStr := range ticketIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errordom.GetSystemError(errordom.INVALID_UUID, "invalid uuid", err)
		}
		ids = append(ids, id)
	}
	return i.repositories.Ticket.BookTickets(ctx, userID, ids)
}

func New(config *utils.Config, repositories repo.Repositories) ticketdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
