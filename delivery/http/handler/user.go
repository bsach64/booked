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
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "could not decode request", err)
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

func (c *CoreHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	var loginRequest userdom.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		ae := errordom.GetSystemError(errordom.JSON_DECODE_ERROR, "could not decode request", err)
		httputils.SendAppError(w, http.StatusBadRequest, nil, ae)
		return
	}

	token, err := c.usecases.UserUC.LoginUser(rCtx, &loginRequest)
	if err != nil {
		ae, ok := err.(*errordom.AppError)
		if !ok {
			httputils.SendAppError(w, http.StatusInternalServerError, nil, err)
			return
		}

		if ae.CategoryCode == errordom.USER_NOT_FOUND || ae.CategoryCode == errordom.INVALID_PASSWORD {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		httputils.SendAppError(w, http.StatusInternalServerError, nil, ae)
		return
	}

	httputils.SendJson(w, http.StatusOK, nil, token)
}
