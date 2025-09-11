package httphandler

import (
	"net/http"

	httputils "github.com/bsach64/booked/delivery/http/utils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	httputils.SendJson(w, http.StatusOK, nil, map[string]string{
		"health":  "ok",
		"message": "Hello from booked!",
	})
}
