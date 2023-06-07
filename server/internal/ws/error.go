package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func buildErrorResponse(msg Request, errMsg string) ResponseBase {
	return ResponseBase{
		Event:        msg.Event,
		Success:      false,
		ErrorMessage: errMsg,
	}
}

func sendErrorMessage(c *Client, msgRes interface{}) {
	msgResBytes, err := json.Marshal(msgRes)
	if err != nil {
		log.Printf("Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.send <- msgResBytes
}
