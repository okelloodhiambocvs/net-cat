package chat

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Hub struct {
	Clients  map[net.Conn]*Client
	Messages []string
	Mutex    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:  make(map[net.Conn]*Client),
		Messages: make([]string, 0),
	}
}

func (h *Hub) AddClient(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Clients[c.Conn] = c

	// Send chat history
	for _, msg := range h.Messages {
		c.Conn.Write([]byte(msg + "\n"))
	}
}

func (h *Hub) RemoveClient(conn net.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, conn)
}

func (h *Hub) Broadcast(sender net.Conn, message string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	client := h.Clients[sender]
	if client == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("[%s][%s]: %s", timestamp, client.Name, message)

	h.Messages = append(h.Messages, formatted)

	for conn, c := range h.Clients {
		if conn == sender {
			continue
		}
		c.Conn.Write([]byte(formatted + "\n"))
	}
}

func (h *Hub) SystemMessage(msg string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("[%s][System]: %s", timestamp, msg)

	h.Messages = append(h.Messages, formatted)

	for _, c := range h.Clients {
		c.Conn.Write([]byte(formatted + "\n"))
	}
}