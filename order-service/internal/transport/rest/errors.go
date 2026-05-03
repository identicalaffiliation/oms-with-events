package rest

import "net/http"

const (
	BAD_REQUEST = "BAD_REQUEST"
	INTERNAL    = "INTERNAL SERVER ERROR"
)

type HTTPError struct {
	Code    int     `json:"code"`
	Message *string `json:"message,omitempty"`
	Type    string  `json:"type"`
}

func (err *HTTPError) Error() string {
	return *err.Message
}

func NewBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")

	e := &HTTPError{
		Code:    http.StatusBadRequest,
		Message: &msg,
		Type:    BAD_REQUEST,
	}

	body, _ := encodeJSON(&e)
	w.Write(body)
}

func NewInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	e := &HTTPError{
		Code: http.StatusInternalServerError,
		Type: INTERNAL,
	}

	body, _ := encodeJSON(&e)
	w.Write(body)
}
