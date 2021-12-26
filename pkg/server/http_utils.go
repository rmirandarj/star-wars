package server

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/go-playground/validator/v10"
)

var govalidator = validator.New()

func decodeAndValidate(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		writeJsonResponse(w, http.StatusBadRequest, errorMessage{Message: "failed to decode payload", ErrorCode: "WA:007"})
		return err
	}
	if err := govalidator.Struct(dest); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make([]map[string]string, 0, len(validationErrors))
			for _, v := range validationErrors {
				details = append(details, map[string]string{"name": v.Field(), "reason": v.Error()})
			}
			writeJsonResponse(w, http.StatusUnprocessableEntity, errorMessage{Message: "payload is invalid", ErrorCode: "WA:001", Details: details})
			return err
		}
		writeJsonResponse(w, http.StatusInternalServerError, errorMessage{Message: "failed to validate payload", ErrorCode: "WA:006"})
		return err
	}
	return nil
}

type errorMessage struct {
	ErrorCode string              `json:"error_code,omitempty"`
	Message   string              `json:"message,omitempty"`
	Details   []map[string]string `json:"details,omitempty"`
}

func writeJsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
