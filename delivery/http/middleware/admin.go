package httpmiddleware

import (
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (m *CoreMiddleware) Admin(next http.Handler) http.Handler {
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

		if user.Role != userdom.ADMIN {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
