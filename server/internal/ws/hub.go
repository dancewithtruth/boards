package ws

import "github.com/Wave-95/boards/server/internal/models"

type Hub struct {
	// Board ID
	boardId string

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

func newHub(boardId string, destroy chan<- string) *Hub {
	return &Hub{
		boardId:    boardId,
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
				if len(h.clients) == 0 {
					h.destroy <- h.boardId
					break
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
		users = append(users, client.user)
	}
	return users
}
