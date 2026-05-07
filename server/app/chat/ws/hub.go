package ws

import (
	"encoding/json"
	"sync"
)

type WsMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Seq  int64           `json:"seq"`
}

type Hub struct {
	mu        sync.RWMutex
	clients   map[int64]*Client
	OnMessage func(userID int64, msg *WsMessage)
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int64]*Client)}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c.UserID] = c
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c.UserID)
	h.mu.Unlock()
}

func (h *Hub) SendToUser(userID int64, msg *WsMessage) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		c.Send(msg)
	}
}

func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
