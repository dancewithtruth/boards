package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	InvalidRequest      = "Invalid request"
	JsonDecodingFailure = "Unable to decode json request"
)

type Validator interface {
	Validate() error
}

type ErrResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	errResponse := ErrResponse{
		Status:  statusCode,
		Message: errMsg,
	}
	json.NewEncoder(w).Encode(errResponse)
}

func WriteWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func buildDecodeErrorMsg(field string, want string, got string) string {
	return fmt.Sprintf("Expected %s to be %s, got %s", field, want, got)
}

func HandleDecodeErr(w http.ResponseWriter, err error) {
	errMsg := JsonDecodingFailure
	if err, ok := err.(*json.UnmarshalTypeError); ok {
		errMsg = buildDecodeErrorMsg(err.Field, err.Type.String(), err.Value)
	}
	WriteWithError(w, http.StatusBadRequest, errMsg)
}

func HandleValidationErr(w http.ResponseWriter, err error) {
	errMsg := InvalidRequest
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErr := fieldErrors[0]
		switch fieldErr.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is a required field", fieldErr.Field())
		default:
			errMsg = fmt.Sprintf("something wrong on %s; %s", fieldErr.Field(), fieldErr.Tag())
		}
	}
	WriteWithError(w, http.StatusBadRequest, errMsg)
}
