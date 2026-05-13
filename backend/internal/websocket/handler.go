package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, we should check origin. For now, allow all.
		return true
	},
}

// ServeWS handles websocket requests from the peer.
func (h *Hub) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("Websocket upgrade failed", zap.Error(err))
		return
	}

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = c.ClientIP()
	}

	client := &Client{
		ID:   clientID,
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	h.log.Info("Websocket client connected",
		zap.String("client_id", client.ID),
		zap.String("remote_addr", conn.RemoteAddr().String()),
	)

	// Start pumps in new goroutines
	go client.writePump()
	go client.readPump()
}
