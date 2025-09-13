package httpmiddleware

import (
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
)

func (m *CoreMiddleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		user, ok := rCtx.Value(httputils.USER_CTX_KEY).(*userdom.User)
		if !ok {
			ae := errordom.GetUserError(errordom.USER_NOT_FOUND, "user not in ctx", nil)
			httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
			return
		}

		if user.Role != userdom.ADMIN {
			ae := errordom.GetUserError(errordom.INVALID_USER_ROLE, "user is not admin", nil)
			httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
			return
		}

		next.ServeHTTP(w, r)
	})
}
