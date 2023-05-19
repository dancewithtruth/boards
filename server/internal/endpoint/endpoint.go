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

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, request Validator, defaultErr error) error {
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("decode err", err)

		if err, ok := err.(*json.UnmarshalTypeError); ok {
			decodeErr := buildDecodeError(err.Field, err.Type.String(), err.Value)
			return decodeErr
		}
		return defaultErr
	}
	defer r.Body.Close()
	if err := request.Validate(); err != nil {
		return err
	}
	return nil
}
