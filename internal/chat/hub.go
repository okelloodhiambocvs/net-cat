package chat

import (
	"net"
	"sync"
)

type Hub struct {
	Clients map[net.Conn]*Client
	Mutex   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[net.Conn]*Client),
	}
}

func (h *Hub) AddClient(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Clients[c.Conn] = c
}

func (h *Hub) RemoveClient(conn net.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, conn)
}