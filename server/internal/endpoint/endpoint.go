package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Validator interface {
	Validate() error
}

type ErrResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteWithError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	errResponse := ErrResponse{
		Status:  statusCode,
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(errResponse)
}

func WriteWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func buildDecodeError(field string, want string, got string) error {
	return errors.New(fmt.Sprintf("Expected %s to be %s, got %s", field, want, got))
}

func HandleDecodeErr(err error, defaultErr error) error {
	if err, ok := err.(*json.UnmarshalTypeError); ok {
		decodeErr := buildDecodeError(err.Field, err.Type.String(), err.Value)
		return decodeErr
	}
	return defaultErr
}
