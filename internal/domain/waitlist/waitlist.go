package waitlistdom

import (
	"context"
	"time"

	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/google/uuid"
)

type NotificationStatus string

const (
	TO_NOTIFY NotificationStatus = "to_notify"
	NOTIFIED  NotificationStatus = "notified"
)

type WaitlistEntry struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	EventID uuid.UUID
	Count   int
}

type WaitlistNotification struct {
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	EventName  string    `json:"event_name"`
	EventTime  time.Time `json:"event_time"`
	SeatCount  int       `json:"seat_count"`
	EventID    uuid.UUID
	WaitlistID uuid.UUID
}

type WaitlistReqeust struct {
	EventID string `json:"event_id"`
	Count   int    `json:"count"`
}

type Usecase interface {
	AddToWaitlist(ctx context.Context, user *userdom.User, request *WaitlistReqeust) error
	NotifyUsers(ctx context.Context) error
}

type Repository interface {
	AddToWaitlist(ctx context.Context, user *userdom.User, eventID uuid.UUID, count int) error
	UpdateNotificationStatus(ctx context.Context, IDs uuid.UUIDs, status NotificationStatus) error
	GetWaitlistNotifications(ctx context.Context) ([]*WaitlistNotification, error)
}
