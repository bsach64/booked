package eventuc

import (
	"context"
	"time"

	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/utils"
)

type impl struct {
	config       *utils.Config
	repositories repo.Repositories
}

func (i *impl) CreateEvent(ctx context.Context, eventRequest *eventdom.CreateEventRequest) error {
	event := &eventdom.Event{
		Name:        eventRequest.Name,
		Address:     eventRequest.Address,
		Description: eventRequest.Description,
		Time:        time.Unix(eventRequest.UnixTime, 0),
	}

	if eventRequest.Latitude != nil && eventRequest.Longitude != nil {
		event.Latitude = *eventRequest.Latitude
		event.Longitude = *eventRequest.Longitude
	}

	eventID, err := i.repositories.Event.CreateEvent(ctx, event)
	if err != nil {
		return err
	}

	err = i.repositories.Ticket.CreateTickets(ctx, eventID, eventRequest.SeatCount)
	if err != nil {
		return err
	}
	return nil
}

func New(config *utils.Config, repositories repo.Repositories) eventdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
