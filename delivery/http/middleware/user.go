package httpmiddleware

import (
	"context"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
)

func (m *CoreMiddleware) SetUserInCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		email := r.Header.Get("X-EMAIL")
		if email == "" {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
			return
		}

		user, err := m.usecases.UserUC.GetUserByEmail(rCtx, email)
		if err != nil {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, err)
			return
		}

		if user == nil {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, err)
			return
		}

		rCtx = context.WithValue(rCtx, "USER", user)
		r = r.WithContext(rCtx)
		next.ServeHTTP(w, r)
	})
}
