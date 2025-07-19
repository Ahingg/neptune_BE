package webSocketService

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	SendUpdateToClient(submissionID uuid.UUID, payload interface{})
	Register(submissionID uuid.UUID, conn *websocket.Conn)
	Unregister(submissionID uuid.UUID, connToRemove *websocket.Conn)
}
