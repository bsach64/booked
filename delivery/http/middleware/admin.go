package httpmiddleware

import (
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (m *CoreMiddleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		user, ok := rCtx.Value("USER").(*userdom.User)
		if !ok {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
			return
		}

		if user.Role != userdom.ADMIN {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
