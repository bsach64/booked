package httpmiddleware

import (
	"net/http"
	"strings"

	httputils "github.com/bsach64/booked/delivery/http/utils"
	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

func (m *CoreMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenBearerString := r.Header.Get("Authorization")
		if tokenBearerString == "" {
			ae := errordom.GetUserError(errordom.INVALID_TOKEN, "jwt token not found", nil)
			httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
			return
		}

		tokenString := strings.TrimPrefix(tokenBearerString, "Bearer ")

		claims := &userdom.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return []byte(m.config.JwtSecret), nil
		})

		if err != nil {
			httputils.SendAppError(w, http.StatusUnauthorized, nil, err)
			return
		}

		if !token.Valid {
			ae := errordom.GetUserError(errordom.INVALID_TOKEN, "invalid jwt token", nil)
			httputils.SendAppError(w, http.StatusUnauthorized, nil, ae)
			return
		}

		r.Header.Set(httputils.EMAIL_HEADER, claims.Email)
		next.ServeHTTP(w, r)
	})
}
