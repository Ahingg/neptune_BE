package websocketHand

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	webSocketService "neptune/backend/services/web_socket_service"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Change origin checkin prod.
		return true
	},
}

type WebSocketHandler struct {
	service webSocketService.WebSocketService
}

func NewWebSocketHandler(service webSocketService.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

func (h *WebSocketHandler) HandleSubmissionConnection(c *gin.Context) {
	submissionIDStr := c.Param("submissionId")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Ensure connection is cleaned up when the handler exits
	defer conn.Close()
	defer h.service.Unregister(submissionID, conn)

	// Register the new connection with our manager
	h.service.Register(submissionID, conn)

	// Start a read loop to keep the connection alive and detect when the client closes it.
	// This is crucial for the defer statements above to execute.
	for {
		if _, _, err := conn.NextReader(); err != nil {
			break // Exit loop if client closes connection
		}
	}
}
