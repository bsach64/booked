package ticketrepo

import (
	"context"
	"errors"
	"fmt"

	"strconv"

	errordom "github.com/bsach64/booked/internal/domain/error"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
)

type impl struct {
	config       *utils.Config
	queries      *db.Queries
	dbConn       *pgxpool.Pool
	valkeyClient valkey.Client
}

func (i *impl) CreateTickets(ctx context.Context, eventID uuid.UUID, count int) error {
	var tickets []db.CreateTicketsParams

	for range count {
		tickets = append(tickets, db.CreateTicketsParams{
			ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
			EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		})
	}

	insertedCount, err := i.queries.CreateTickets(ctx, tickets)
	if err != nil || insertedCount != int64(count) {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not create tickets", err)
	}
	return nil
}

func (i *impl) ReserveTickets(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, count int) ([]uuid.UUID, error) {
	tx, err := i.dbConn.Begin(ctx)
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_TX_ERROR, "could not start tx", err)
	}
	defer tx.Rollback(ctx)

	qtx := i.queries.WithTx(tx)
	dbTicketIDs, err := qtx.GetAvailableTickets(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not get available tickets", err)
	}

	if len(dbTicketIDs) < count {
		return nil, errordom.GetEventError(errordom.TOO_FEW_TICKETS, "too few tickets availiable", nil)
	}

	var dbTicketIDStrs []string
	for _, id := range dbTicketIDs {
		dbTicketIDStrs = append(dbTicketIDStrs, id.String())
	}

	/*
		We know that we have N available tickets, suppose we R reserved tickets.

		We can only reserve "count" tickets if:
			N - R >= count

		If this is the case, then we should reserve them.
		All of this should happen simultaneously.

		The best way to do this I found was through lua script.
	*/

	luaScript := `
		local user_id = ARGV[1]
		local ttl = tonumber(ARGV[2])
		local count = tonumber(ARGV[3])

		local available_ids = {}
		for i = 1, #KEYS do
			if redis.call("EXISTS", KEYS[i]) == 0 then
				table.insert(available_ids, KEYS[i])
			end
		end

		if #available_ids < count then
			return {}
		end

		local reserved_tickets = {}
		for i = 1, count do
			redis.call("SET", available_ids[i], user_id, "EX", ttl, "NX")
			table.insert(reserved_tickets,  available_ids[i])
		end

		return reserved_tickets
	`

	script := valkey.NewLuaScript(luaScript)

	reservedTickets, err := script.Exec(ctx, i.valkeyClient, dbTicketIDStrs, []string{userID.String(), "600", strconv.Itoa(count)}).ToArray()
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_TX_ERROR, "failed to reserve tickets", err)
	}

	if len(reservedTickets) != count {
		return nil, errordom.GetEventError(errordom.TOO_FEW_TICKETS, "majority tickets reserved", err)
	}

	reservedTicketsIDs := uuid.UUIDs{}

	for _, ticket := range reservedTickets {
		idStr, err := ticket.ToString()
		if err != nil {
			return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "failed to get redis key as string", err)
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errordom.GetSystemError(errordom.INVALID_UUID, "could not get uuid", err)
		}

		reservedTicketsIDs = append(reservedTicketsIDs, id)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_TX_ERROR, "failed to commit tx", err)
	}

	return reservedTicketsIDs, nil
}

func (i *impl) BookTickets(ctx context.Context, userID uuid.UUID, ticketIDs []uuid.UUID) error {
	dbTicketIDs := []pgtype.UUID{}
	dbUserID := pgtype.UUID{Bytes: userID, Valid: true}

	tx, err := i.dbConn.Begin(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "could not start tx", err)
	}
	defer tx.Rollback(ctx)

	qtx := i.queries.WithTx(tx)
	for _, ticket := range ticketIDs {
		// it should have been reserved for this person
		// we dont really need to worry about concurrency here
		userIDForTicket, err := i.valkeyClient.Do(ctx, i.valkeyClient.B().Get().Key(ticket.String()).Build()).ToString()
		if err != nil {
			if errors.Is(err, valkey.Nil) {
				return errordom.GetTicketError(errordom.TICKET_NOT_RESERVED, "ticket not reserved or reservation expired", nil)
			}
			return errordom.GetDBError(errordom.DB_READ_ERROR, "could not read from cache", err)
		}

		if userIDForTicket != userID.String() {
			return errordom.GetTicketError(errordom.NOT_YOUR_TICKET, "ticket reserved by someone else", nil)
		}

		dbTicketIDs = append(dbTicketIDs, pgtype.UUID{Bytes: ticket, Valid: true})
	}

	err = qtx.BookTickets(ctx, db.BookTicketsParams{UserID: dbUserID, Column2: dbTicketIDs})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "failed to commit tx", err)
	}
	return nil
}

func (i *impl) GetPastBookings(ctx context.Context, userID uuid.UUID) ([]*ticketdom.PastBookingsResponse, error) {
	pastBookings := []*ticketdom.PastBookingsResponse{}
	dbResp, err := i.queries.GetBookingHistory(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not get past booking", err)
	}

	for _, row := range dbResp {
		pastBooking := &ticketdom.PastBookingsResponse{
			EventID:          row.ID.String(),
			EventName:        row.Name,
			EventAddress:     row.Address,
			EventDescription: row.Address,
			EventUnixTime:    utils.GetUTCUnixTime(row.Time.Time),
			NumberOfTickets:  int(row.Count),
			BookingUnixTime:  utils.GetUTCUnixTime(row.Column9.Time),
		}

		if row.Latitude.Valid && row.Longitude.Valid {
			pastBooking.EventLatitude = &row.Latitude.Float64
			pastBooking.EventLongitude = &row.Longitude.Float64
		}
		pastBookings = append(pastBookings, pastBooking)
	}
	return pastBookings, nil
}

func (i *impl) CancelTickets(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, count int) error {
	tx, err := i.dbConn.Begin(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "could not start tx", err)
	}
	defer tx.Rollback(ctx)

	qtx := i.queries.WithTx(tx)
	getBookedTickets := db.GetBookedTicketsParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	}

	dbTicketIDs, err := qtx.GetBookedTickets(ctx, getBookedTickets)
	if len(dbTicketIDs) < count {
		msg := fmt.Sprintf("count=%v is greater than number of tickets booked", count)
		return errordom.GetTicketError(errordom.TOO_FEW_TICKETS, msg, nil)
	}

	err = qtx.CancelTickets(ctx, dbTicketIDs[:count])
	if err != nil {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not canncel tickets", nil)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errordom.GetDBError(errordom.DB_TX_ERROR, "failed to commit tx", err)
	}
	return nil
}

func (i *impl) GetReservedTickets(ctx context.Context, eventID uuid.UUID) (int, error) {
	dbTicketIDs, err := i.queries.GetAvailableTickets(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err != nil {
		return 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not get available tickets", err)
	}

	idStrs := []string{}
	for _, id := range dbTicketIDs {
		idStrs = append(idStrs, id.String())
	}

	count, err := i.valkeyClient.Do(ctx, i.valkeyClient.B().Exists().Key(idStrs...).Build()).AsInt64()
	if err != nil {
		return 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not get reserved ticket count", err)
	}
	return int(count), nil
}

func (i *impl) GetAnalytics(ctx context.Context) ([]*ticketdom.Analytics, error) {
	response := []*ticketdom.Analytics{}
	stats, err := i.queries.GetAnalytics(ctx)
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read tickets table", err)
	}

	for _, stat := range stats {
		data := &ticketdom.Analytics{
			EventID:             stat.EventID.String(),
			TotalSeats:          int(stat.TotalSeats),
			SoldTickets:         int(stat.BookedTickets),
			CapacityUtilisation: float64(stat.BookedTickets) / float64(stat.TotalSeats),
		}
		response = append(response, data)
	}
	return response, nil
}

func (i *impl) GetDailyBookings(ctx context.Context) ([]*ticketdom.DailyAnalytics, error) {
	response := []*ticketdom.DailyAnalytics{}
	stats, err := i.queries.GetDailyBookings(ctx)
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read tickets table", err)
	}

	for _, stat := range stats {
		data := &ticketdom.DailyAnalytics{
			EventID:          stat.EventID.String(),
			TodaySoldTickets: int(stat.TodayBookedTickets),
		}
		response = append(response, data)
	}
	return response, nil
}

func (i *impl) GetCancellationRates(ctx context.Context) ([]*ticketdom.CancellationRates, error) {
	response := []*ticketdom.CancellationRates{}
	stats, err := i.queries.GetCancellationRates(ctx)
	if err != nil {
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read tickets table", err)
	}

	for _, stat := range stats {
		data := &ticketdom.CancellationRates{
			EventID:          stat.EventID.String(),
			CancellationRate: stat.CancellationRate,
		}
		response = append(response, data)
	}
	return response, nil
}

func (i *impl) GetAvailableTickets(ctx context.Context, eventID uuid.UUID) (int, error) {
	ids, err := i.queries.GetAvailableTickets(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err != nil {
		return 0, errordom.GetDBError(errordom.DB_READ_ERROR, "could not read tickets table", err)
	}

	return len(ids), nil
}

func New(config *utils.Config, queries *db.Queries, dbConn *pgxpool.Pool, valkeyClient valkey.Client) ticketdom.Repository {
	return &impl{
		config:       config,
		queries:      queries,
		dbConn:       dbConn,
		valkeyClient: valkeyClient,
	}
}
