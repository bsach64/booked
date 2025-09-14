package waitlistuc

import (
	"context"
	"log/slog"

	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	waitlistdom "github.com/bsach64/booked/internal/domain/waitlist"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
)

type impl struct {
	config       *utils.Config
	repositories repo.Repositories
}

func (i *impl) AddToWaitlist(ctx context.Context, user *userdom.User, request *waitlistdom.WaitlistReqeust) error {
	eventID, err := uuid.Parse(request.EventID)
	if err != nil {
		return errordom.GetSystemError(errordom.INVALID_UUID, "invalid uuid", err)
	}

	return i.repositories.Waitlist.AddToWaitlist(ctx, user, eventID, request.Count)
}

func (i *impl) NotifyUsers(ctx context.Context) error {
	notifications, err := i.repositories.Waitlist.GetWaitlistNotifications(ctx)
	if err != nil {
		return err
	}

	notifiedIDs := uuid.UUIDs{}
	for _, notification := range notifications {
		availableTickets, err := i.repositories.Ticket.GetAvailableTickets(ctx, notification.EventID)
		if err != nil {
			return err
		}

		if availableTickets < notification.SeatCount {
			continue
		}

		// notify user
		slog.Info(
			"Notifying User",
			"Name",
			notification.UserName,
			"Email",
			notification.UserEmail,
			"EventName",
			notification.EventName,
			"EventTime",
			notification.EventTime,
			"Ticket Count",
			notification.SeatCount,
		)

		notifiedIDs = append(notifiedIDs, notification.WaitlistID)
	}

	err = i.repositories.Waitlist.UpdateNotificationStatus(ctx, notifiedIDs, waitlistdom.NOTIFIED)
	if err != nil {
		return err
	}
	return nil
}

func New(config *utils.Config, repositories repo.Repositories) waitlistdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
