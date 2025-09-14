package httphandler

import (
	"encoding/json"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	waitlistdom "github.com/bsach64/booked/internal/domain/waitlist"
)

func (c *CoreHandler) AddWaitlistHandler(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	var request waitlistdom.WaitlistReqeust

	user, ok := rCtx.Value(httputils.USER_CTX_KEY).(*userdom.User)
	if !ok {
		httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "could not decode request", err)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	err := c.usecases.WaitlistUC.AddToWaitlist(rCtx, user, &request)
	if err != nil {
		ae, ok := err.(*errordom.AppError)
		if !ok {
			httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
			return
		}

		if ae.CategoryCode == errordom.INVALID_UUID || ae.CategoryCode == errordom.ALREADY_IN_WAILIST || ae.CategoryCode == errordom.INVALID_WAITLIST_COUNT {
			httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
			return
		}

		httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
		return
	}
}
