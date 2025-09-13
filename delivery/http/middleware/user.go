package httpmiddleware

import (
	"context"
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
)

func (m *CoreMiddleware) SetUserInCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		email := r.Header.Get(httputils.EMAIL_HEADER)
		if email == "" {
			ae := errordom.GetUserError(errordom.EMPTY_EMAIL, "email header is empty", nil)
			httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
			return
		}

		user, err := m.usecases.UserUC.GetUserByEmail(rCtx, email)
		if err != nil {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, err)
			return
		}
		rCtx = context.WithValue(rCtx, httputils.USER_CTX_KEY, user)
		r = r.WithContext(rCtx)
		next.ServeHTTP(w, r)
	})
}
