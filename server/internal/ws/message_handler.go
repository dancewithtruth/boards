package ws

import (
	"encoding/json"
	"log"

	"github.com/Wave-95/boards/server/internal/post"
	"github.com/Wave-95/boards/server/pkg/validator"
)

func HandleMessage(c *Client, message []byte) {
	// Identify message type
	var messageRequest MessageRequest
	err := json.Unmarshal(message, &messageRequest)
	if err != nil {
		sendMessageBadType(c)
		return
	}

	// Route message type and handle accordingly
	switch messageRequest.Type {
	case "":
		errorResponse := buildErrorResponse("", ErrorMessageMissingType)
		c.send <- errorResponse
		return
	case TypeCreatePost:
		handleCreatePost(c, messageRequest)
	default:
		errorResponse := buildErrorResponse(messageRequest.Type, ErrorMessageUnsupportedType)
		c.send <- errorResponse
		return
	}
}

func handleCreatePost(c *Client, messageRequest MessageRequest) {
	// Unmarshal input
	var input post.CreatePostInput
	err := json.Unmarshal(messageRequest.Data, &input)
	if err != nil {
		errorResponse := buildErrorResponse(messageRequest.Type, ErrorMessageBadInput)
		c.send <- errorResponse
		return
	}
	// Validate create post payload
	err = c.api.validator.Struct(input)
	if err != nil {
		sendMessageValidationErr(c, messageRequest.Type, input, err)
		return
	}

	// TODO: Create post
	// post, err := c.api.postService.CreatePost(payload)
	// if err != nil {
	// 	sendMessageInternalServerErr(c, ErrorMessageInternalCreatePost)
	// 	return
	// }
}

func sendMessageBadType(c *Client) {
	errorResponse := buildErrorResponse("", ErrorMessageBadType)
	c.send <- errorResponse
}

func sendMessageValidationErr(c *Client, messageType string, s interface{}, err error) {
	errMsg := "Invalid request"
	validationErrMsg := validator.GetValidationErrMsg(s, err)
	if validationErrMsg != "" {
		errMsg = validationErrMsg
	}
	c.send <- buildErrorResponse(messageType, errMsg)
}

func buildErrorResponse(messageType string, messageError string) []byte {
	messageResponse := MessageResponse{
		Type:         messageType,
		ErrorMessage: messageError,
	}
	messageResponseBytes, err := json.Marshal(messageResponse)
	if err != nil {
		log.Printf("Failed to marshal error response into json:%v", err)
	}
	return messageResponseBytes
}
