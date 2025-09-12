package httputils

import (
	"encoding/json"
	"net/http"

	errordom "github.com/bsach64/booked/internal/domain/error"
)

func SendJson(w http.ResponseWriter, status int, headers map[string]string, data any) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Just to make things explicit that we are sending an error
func SendAppError(w http.ResponseWriter, status int, headers map[string]string, err error) {
	ae, ok := err.(*errordom.AppError)
	if !ok {
		ae = errordom.GetSystemError(errordom.UNKNOWN_ERROR, "", err).(*errordom.AppError)
	}
	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if ae != nil {
		json.NewEncoder(w).Encode(ae)
	}
}
