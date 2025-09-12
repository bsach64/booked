package httphandler

import (
	"encoding/json"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (c *CoreHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	var newUser userdom.NewUser

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "", err).(*errordom.AppError)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	err := c.usecases.UserUC.CreateUser(rCtx, newUser)
	if err != nil {
		httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
