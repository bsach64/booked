package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	eventdom "github.com/bsach64/booked/internal/domain/event"
)

func (c *CoreHandler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	var createEventRequest eventdom.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&createEventRequest); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "", err).(*errordom.AppError)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	err := c.usecases.EventUC.CreateEvent(rCtx, &createEventRequest)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *CoreHandler) GetPaginatedEvents(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	limit := 10
	var nextTimeToFetch *int64
	limitStr := r.URL.Query().Get("limit")
	nextTimeToFetchStr := r.URL.Query().Get("timestamp")

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	if t, err := strconv.ParseInt(nextTimeToFetchStr, 10, 64); err == nil && t > 0 {
		nextTimeToFetch = &t
	}

	eventsResponse, err := c.usecases.EventUC.GetEvents(rCtx, limit, nextTimeToFetch)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	httputils.SendJson(w, http.StatusOK, nil, eventsResponse)
}

func (c *CoreHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	id := r.URL.Query().Get("id")

	statusCode := http.StatusInternalServerError
	err := c.usecases.EventUC.DeleteEvent(rCtx, id)
	if err != nil {
		ae, ok := err.(*errordom.AppError)
		if !ok {
			httputils.SendAppError(w, statusCode, nil, err)
			return
		}

		if ae.CategoryCode == errordom.NO_EVENT_FOUND {
			statusCode = http.StatusBadRequest
		}

		httputils.SendAppError(w, statusCode, nil, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
