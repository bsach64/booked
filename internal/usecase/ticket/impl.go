package ticketuc

import (
	"context"

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
	return nil, nil
}

func (i *impl) BookTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
	return nil
}

func New(config *utils.Config, repositories repo.Repositories) ticketdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
