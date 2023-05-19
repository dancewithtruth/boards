package response

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteWithError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	errResponse := ErrResponse{
		Status:  statusCode,
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(errResponse)
}

func WriteWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
