package ws

import (
	"encoding/json"

	"github.com/Wave-95/boards/server/internal/api/post"
	"github.com/Wave-95/boards/server/pkg/validator"
)

func HandleMessage(c *Client, message []byte) {
	// Identify message type
	var messageRequest MessageRequest
	err := json.Unmarshal(message, &messageRequest)
	if err != nil {
		sendMessageBadRequest(c)
		return
	}

	// Route message type and handle accordingly
	switch messageRequest.Type {
	case TypeCreatePost:
		// Unmarshal payload
		var payload post.CreatePostInput
		err := json.Unmarshal(messageRequest.Payload, &payload)
		if err != nil {
			sendMessageBadRequest(c)
			return
		}
		// Validate create post payload
		err = c.api.validator.Struct(payload)
		if err != nil {
			sendMessageValidationErr(c, payload, err)
			return
		}

		// TODO: Create post
		// post, err := c.api.postService.CreatePost(payload)
		// if err != nil {
		// 	sendMessageInternalServerErr(c, ErrorMessageInternalCreatePost)
		// 	return
		// }

	default:
		errorResponse := buildErrorResponse(TypeBadRequest, ErrorMessageBadType)
		c.send <- errorResponse
		return
	}
}

func sendMessageBadRequest(c *Client) {
	errorResponse := buildErrorResponse(TypeBadRequest, ErrorMessageBadRequest)
	c.send <- errorResponse
}

func sendMessageValidationErr(c *Client, s interface{}, err error) {
	errMsg := "Invalid request"
	validationErrMsg := validator.GetValidationErrMsg(s, err)
	if validationErrMsg != "" {
		errMsg = validationErrMsg
	}
	c.send <- buildErrorResponse(TypeInvalidRequest, errMsg)
}

func sendMessageInternalServerErr(c *Client, errMsg string) {
	errorResponse := buildErrorResponse(TypeInternalServerError, errMsg)
	c.send <- errorResponse
}

func buildErrorResponse(messageType string, messageError string) []byte {
	messageResponse := MessageResponse{
		Type:         messageType,
		ErrorMessage: messageError,
	}
	messageResponseBytes, _ := json.Marshal(messageResponse)
	return messageResponseBytes
}
