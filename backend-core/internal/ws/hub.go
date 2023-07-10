package ws

import (
	"github.com/Wave-95/boards/backend-core/internal/models"
)

// Hub represents a single board that clients are connected to. It encapsulates logic to
// register & unregister clients, broadcast messages to clients, and destroy itself to free
// up resources.
type Hub struct {
	// Board ID
	boardID string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Destroy sends a destroy request to delete this Hub from BoardHubs map
	destroy chan<- string
}

func newHub(boardID string, destroy chan<- string) *Hub {
	return &Hub{
		boardID:    boardID,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		destroy:    destroy,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Destroy the hub if all clients are disconnected to free resources
				if len(h.clients) == 0 {
					h.destroy <- h.boardID
					return
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) listConnectedUsers() []models.User {
	users := []models.User{}
	for client := range h.clients {
		users = append(users, *client.user)
	}

	return users
}
