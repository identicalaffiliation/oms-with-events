package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func decodeJSON(r *http.Request, i any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(i); err != nil {
		return err
	}

	return nil
}

func encodeJSON(i any) ([]byte, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func encodeResponse(w http.ResponseWriter, i any, code int) {
	switch code {
	case http.StatusCreated:
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		body, _ := encodeJSON(i)
		w.Write(body)
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		body, _ := encodeJSON(i)
		w.Write(body)
	default:
		return
	}
}

func ValidateRequest(v *validator.Validate, i any) error {
	return v.Struct(i)
}
