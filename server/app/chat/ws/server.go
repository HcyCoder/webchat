package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/team/webchat-server/common/token"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	hub          *Hub
	tokenManager *token.Manager
}

func NewServer(hub *Hub, tm *token.Manager) *Server {
	return &Server{hub: hub, tokenManager: tm}
}

func (s *Server) Listen(addr string) error {
	http.HandleFunc("/ws", s.handleWS)
	log.Printf("chat WebSocket on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	userIDStr, err := s.tokenManager.Validate(r.Context(), tokenStr)
	if err != nil || userIDStr == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := NewClient(s.hub, conn, userID)
	s.hub.Register(client)
	go client.WritePump()
	go client.ReadPump()
}
