package services

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocketService struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	logger     *zap.Logger
}

func NewWebSocketService(logger *zap.Logger) *WebSocketService {
	return &WebSocketService{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		logger:     logger,
	}
}

func (s *WebSocketService) Start() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					s.logger.Error("Failed to send message", zap.Error(err))
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

// BroadcastTaskUpdate sends task updates to all connected clients
func (s *WebSocketService) BroadcastTaskUpdate(event string, data interface{}) {
	update := struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}{
		Event: event,
		Data:  data,
	}

	message, err := json.Marshal(update)
	if err != nil {
		s.logger.Error("Failed to marshal update", zap.Error(err))
		return
	}

	s.broadcast <- message
}

// Events that can be broadcast
const (
	TaskCreatedEvent  = "task.created"
	TaskUpdatedEvent  = "task.updated"
	TaskDeletedEvent  = "task.deleted"
	TaskStatusEvent   = "task.status"
	TaskProgressEvent = "task.progress"
) 