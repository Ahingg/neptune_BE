package webSocketService

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type webSocketService struct {
	clients map[uuid.UUID][]*websocket.Conn
	sync.RWMutex
}

// Register adds a new websocket connection to the manager.
func (s *webSocketService) Register(submissionID uuid.UUID, conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	s.clients[submissionID] = append(s.clients[submissionID], conn)
	log.Printf("Registered new websocket client for submission %s", submissionID)
}

// Unregister removes a websocket connection.
func (s *webSocketService) Unregister(submissionID uuid.UUID, connToRemove *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	connections := s.clients[submissionID]
	for i, conn := range connections {
		if conn == connToRemove {
			// Remove the connection from the slice
			s.clients[submissionID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	// If no connections are left for this submission, clean up the map entry
	if len(s.clients[submissionID]) == 0 {
		delete(s.clients, submissionID)
	}
	log.Printf("Unregistered websocket client for submission %s", submissionID)
}

func (s *webSocketService) SendUpdateToClient(submissionID uuid.UUID, payload interface{}) {
	s.RLock() // Use a Read Lock as we are only reading the map
	defer s.RUnlock()

	connections, found := s.clients[submissionID]
	if !found {
		// No client is currently listening for this submission, which is fine.
		return
	}

	log.Printf("Broadcasting update to %d clients for submission %s", len(connections), submissionID)
	for _, conn := range connections {
		// WriteJSON is safe for concurrent use on a single connection.
		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("Websocket write error: %v", err)
			// Optional: In a real app, you might try to unregister this failed connection.
		}
	}
}

func NewWebSocketService() WebSocketService {
	return &webSocketService{
		clients: make(map[uuid.UUID][]*websocket.Conn),
	}
}
