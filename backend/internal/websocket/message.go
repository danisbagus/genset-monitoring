package websocket

import (
	"encoding/json"
	"time"
)

// Envelope defines the standard structure for websocket messages.
type Envelope struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

// NewMessage creates a new JSON-encoded message envelope.
func NewMessage(event string, data any) ([]byte, error) {
	env := Envelope{
		Event:     event,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	return json.Marshal(env)
}
