package httphandler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (c *CoreHandler) ReserveTickets(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	user, ok := rCtx.Value("USER").(*userdom.User)
	if !ok {
		slog.Info("did not get user")
		httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
		return
	}

	var reserveTicketRequest ticketdom.ReserveTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&reserveTicketRequest); err != nil {
		slog.Info("could not get body")
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "", err).(*errordom.AppError)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	reserveTicketRequest.UserID = user.ID
	resp, err := c.usecases.TicketUC.ReserveTickets(rCtx, &reserveTicketRequest)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	slog.Info("got resp", "resp", resp)

	httputils.SendJson(w, http.StatusOK, nil, resp)
}

func (c *CoreHandler) BookTickets(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	bookTicketRequest := struct {
		TicketIDs []string `json:"ticket_ids"`
	}{}

	user, ok := rCtx.Value("USER").(*userdom.User)
	if !ok {
		httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&bookTicketRequest); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "", err).(*errordom.AppError)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	err := c.usecases.TicketUC.BookTickets(rCtx, user.ID, bookTicketRequest.TicketIDs)
	if err != nil {
		httputils.SendAppError(w, http.StatusBadRequest, nil, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
