package response

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	errResponse := ErrResponse{
		Status:  statusCode,
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(errResponse)
}
