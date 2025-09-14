package eventuc

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	eventdom "github.com/bsach64/booked/internal/domain/event"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
)

type impl struct {
	config       *utils.Config
	repositories repo.Repositories
}

func (i *impl) GetEvents(ctx context.Context, limit int, lastFetchedUnixTime *int64) (*eventdom.GetEventsResponse, error) {
	var events []*eventdom.Event
	var err error
	var nextTimeToFetch int64
	if lastFetchedUnixTime == nil {
		events, nextTimeToFetch, err = i.repositories.Event.GetEvents(ctx, limit)
		if err != nil {
			return nil, err
		}
	} else {
		events, nextTimeToFetch, err = i.repositories.Event.GetNextEvents(ctx, *lastFetchedUnixTime, limit)
		if err != nil {
			return nil, err
		}
	}

	response := &eventdom.GetEventsResponse{}
	response.NextTimeToFetch = nextTimeToFetch

	for _, event := range events {
		requestEvent := &eventdom.EventResponse{
			Name:                event.Name,
			Description:         event.Description,
			Address:             event.Address,
			UnixTime:            utils.GetUTCUnixTime(event.Time),
			SeatCount:           event.SeatCount,
			ID:                  event.ID.String(),
			AvailableSeatsCount: event.AvailableTickets,
		}

		if event.Latitude != nil && event.Longitude != nil {
			requestEvent.Latitude = event.Latitude
			requestEvent.Longitude = event.Longitude
		}

		reservedTickets, err := i.repositories.Ticket.GetReservedTickets(ctx, event.ID)
		if err != nil {
			return nil, err
		}
		requestEvent.AvailableSeatsCount -= reservedTickets

		response.Events = append(response.Events, requestEvent)
	}

	return response, nil
}

func (i *impl) CreateEvent(ctx context.Context, eventRequest *eventdom.CreateEventRequest) error {
	event := &eventdom.Event{
		Name:        eventRequest.Name,
		Address:     eventRequest.Address,
		Description: eventRequest.Description,
		Time:        utils.GetUTCTimeFromUnix(eventRequest.UnixTime),
		SeatCount:   eventRequest.SeatCount,
	}

	if eventRequest.Name == "" || eventRequest.Description == "" || eventRequest.Address == "" || eventRequest.UnixTime == 0 {
		return errordom.GetEventError(errordom.INVALID_NEW_EVENT, "please provide valid event details", nil)
	}

	if eventRequest.SeatCount <= 0 {
		return errordom.GetEventError(errordom.INVALID_SEAT_COUNT, "seat count less than 0", nil)
	}

	if eventRequest.Latitude != nil && eventRequest.Longitude != nil {
		event.Latitude = eventRequest.Latitude
		event.Longitude = eventRequest.Longitude
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

func (i *impl) DeleteEvent(ctx context.Context, eventID string) error {
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return errordom.GetEventError(errordom.INVALID_EVENT_ID, "invalid uuid provided", err)
	}

	return i.repositories.Event.DeleteEvent(ctx, eventUUID)
}

func (i *impl) UpdateEvent(ctx context.Context, updateEventRequest *eventdom.UpdateEventRequest) error {
	return i.repositories.Event.UpdateEvent(ctx, updateEventRequest)
}

func New(config *utils.Config, repositories repo.Repositories) eventdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
