package httputils

import (
	"encoding/json"
	"net/http"

	errordom "github.com/bsach64/booked/internal/domain/error"
)

const USER_CTX_KEY = "USER"
const EMAIL_HEADER = "X-EMAIL"

type ErrorResponse struct {
	Category     string `json:"category"`
	CategoryCode string `json:"category_code"`
	Msg          string `json:"msg"`
	Error        error  `json:"error,omitempty"`
}

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
	errRes := &ErrorResponse{
		Category:     string(errordom.CATEGORY_SYSTEM),
		CategoryCode: string(errordom.UNKNOWN_ERROR),
		Msg:          "an unknowned error occurred, please try again later",
		Error:        err,
	}

	if err != nil {
		ae, ok := err.(*errordom.AppError)
		if ok {
			errRes.Category = string(ae.Category)
			errRes.CategoryCode = string(ae.CategoryCode)
			errRes.Msg = ae.Msg
			if ae.ErrorToWrap != nil {
				errRes.Error = ae.ErrorToWrap
			}
		}
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errRes)
}
