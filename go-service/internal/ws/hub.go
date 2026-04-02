package ws

import (
	"encoding/json"
	"sync"
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

type Hub struct {
	mu          sync.Mutex
	connections map[*Conn]struct{}
}

type Conn struct {
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[*Conn]struct{}),
	}
}

func (h *Hub) Register(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[c] = struct{}{}
}

func (h *Hub) Unregister(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, c)
	close(c.send)
}

func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.connections {
		select {
		case c.send <- data:
		default:
		}
	}
}
