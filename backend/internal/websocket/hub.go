package websocket

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins in development; restrict in production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client represents a connected WebSocket client.
type Client struct {
	ID   string
	conn *websocket.Conn
	send chan []byte
}

// Hub manages all active WebSocket clients and message broadcasting.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
	log     *zap.Logger
}

// NewHub creates a new WebSocket Hub.
func NewHub(log *zap.Logger) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		log:     log,
	}
}

// ServeWS handles WebSocket upgrade and client lifecycle.
func (h *Hub) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("websocket upgrade failed", zap.Error(err))
		return
	}

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = c.ClientIP()
	}

	client := &Client{
		ID:   clientID,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register(client)
	h.log.Info("websocket client connected", zap.String("client_id", clientID))

	go h.writePump(client)
	h.readPump(client)
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.send <- msg:
		default:
			h.log.Warn("websocket send buffer full, dropping message",
				zap.String("client_id", client.ID))
		}
	}
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	h.clients[c.ID] = c
	h.mu.Unlock()
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c.ID]; ok {
		delete(h.clients, c.ID)
		close(c.send)
	}
	h.mu.Unlock()
}

func (h *Hub) readPump(c *Client) {
	defer func() {
		h.unregister(c)
		c.conn.Close()
		h.log.Info("websocket client disconnected", zap.String("client_id", c.ID))
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Warn("websocket unexpected close", zap.Error(err))
			}
			break
		}
	}
}

func (h *Hub) writePump(c *Client) {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			h.log.Warn("websocket write error",
				zap.String("client_id", c.ID),
				zap.Error(err))
			return
		}
	}
}
