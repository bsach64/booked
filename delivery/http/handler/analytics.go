package httphandler

import (
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
)

func (c *CoreHandler) TotalBookingsHandler(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	resp, err := c.repositories.Ticket.GetAnalytics(rCtx)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	httputils.SendJson(w, http.StatusOK, nil, resp)
}

func (c *CoreHandler) DailyBookingsHandler(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	resp, err := c.repositories.Ticket.GetDailyBookings(rCtx)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	httputils.SendJson(w, http.StatusOK, nil, resp)
}

func (c *CoreHandler) CancellationRatesHandler(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

	resp, err := c.repositories.Ticket.GetCancellationRates(rCtx)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	httputils.SendJson(w, http.StatusOK, nil, resp)
}
