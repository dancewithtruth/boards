// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Wave-95/boards/server/internal/api/post"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	api *API

	userId string

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		handleMessage(c, message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleMessage(c *Client, message []byte) {
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
		fmt.Println(payload)
		if err != nil {
			sendMessageValidationErr(c, err)
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

func sendMessageValidationErr(c *Client, err error) {
	errMsg := "Invalid request"
	validationErrMsg := validator.GetValidationErrMsg(err)
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
