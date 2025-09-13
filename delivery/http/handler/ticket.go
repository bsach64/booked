package httphandler

import (
	"encoding/json"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	ticketdom "github.com/bsach64/booked/internal/domain/ticket"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (c *CoreHandler) ReserveTickets(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	user, ok := rCtx.Value(httputils.USER_CTX_KEY).(*userdom.User)
	if !ok {
		httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
		return
	}

	var reserveTicketRequest ticketdom.ReserveTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&reserveTicketRequest); err != nil {
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

	httputils.SendJson(w, http.StatusOK, nil, resp)
}

func (c *CoreHandler) BookTickets(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	bookTicketRequest := struct {
		TicketIDs []string `json:"ticket_ids"`
	}{}

	user, ok := rCtx.Value(httputils.USER_CTX_KEY).(*userdom.User)
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

func (c *CoreHandler) CancelTickets(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	user, ok := rCtx.Value(httputils.USER_CTX_KEY).(*userdom.User)
	if !ok {
		httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
		return
	}

	var cancelTicketRequest ticketdom.CancelTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&cancelTicketRequest); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "", err).(*errordom.AppError)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	err := c.usecases.TicketUC.CancelTickets(rCtx, user, &cancelTicketRequest)
	if err != nil {
		ae, ok := err.(*errordom.AppError)
		if !ok {
			httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
			return
		}

		if ae.CategoryCode == errordom.ErrorCategoryCode(errordom.TOO_FEW_TICKETS) ||
			ae.CategoryCode == errordom.ErrorCategoryCode(errordom.INVALID_UUID) {
			httputils.SendAppError(w, http.StatusBadRequest, nil, err)
			return
		}
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}
}
