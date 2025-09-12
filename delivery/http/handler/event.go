package httphandler

import (
	"encoding/json"
	"net/http"

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
