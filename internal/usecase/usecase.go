package usecase

import (
	eventdom "github.com/bsach64/booked/internal/domain/event"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo"
	eventuc "github.com/bsach64/booked/internal/usecase/event"
	ticketuc "github.com/bsach64/booked/internal/usecase/ticket"
	useruc "github.com/bsach64/booked/internal/usecase/user"
	"github.com/bsach64/booked/utils"
)

type Usecase struct {
	config   *utils.Config
	UserUC   userdom.Usecase
	EventUC  eventdom.Usecase
	TicketUC ticketdom.Usecase
}

func New(config *utils.Config, repositories repo.Repositories) Usecase {
	return Usecase{
		config:   config,
		UserUC:   useruc.New(config, repositories),
		EventUC:  eventuc.New(config, repositories),
		TicketUC: ticketuc.New(config, repositories),
	}
}
