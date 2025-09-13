package ticketuc

import (
	"context"
	"fmt"
	"log/slog"

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
		// crazy
		return nil, fmt.Errorf("uuid")
	}
	slog.Info("got uuid")

	ticketIDs, err := i.repositories.Ticket.GetAvailiableTickets(ctx, eventID)
	if err != nil {
		return nil, err
	}
	slog.Info("got availiable tickets", "ticket_ids", ticketIDs)

	if len(ticketIDs) < reserveTickets.Count {
		// return an error here
		return nil, fmt.Errorf("to few")

	}

	err = i.repositories.Ticket.ReserveTickets(ctx, reserveTickets.UserID, ticketIDs[:reserveTickets.Count])
	if err != nil {
		// do error handling
		return nil, err
	}
	slog.Info("reserved tickets")

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
			return err
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
