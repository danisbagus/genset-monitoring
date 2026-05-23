package pubsub

import (
	"sync"
)

// Topic defines the type for pub/sub event topics.
type Topic string

const (
	// TopicDeviceCRUD triggered when a device is created, updated, or deleted
	TopicDeviceCRUD Topic = "device.crud"
	// TopicAlertCreated triggered when a new alert is created
	TopicAlertCreated Topic = "alert.created"
	// TopicHeartbeatReceived triggered when a device heartbeat is recorded
	TopicHeartbeatReceived Topic = "heartbeat.received"
	// TopicTelemetryEngineCreated triggered when a telemetry engine is created
	TopicTelemetryEngineCreated Topic = "telemetry.engine.created"
	// TopicTelemetryElectricalCreated triggered when a telemetry electrical is created
	TopicTelemetryElectricalCreated Topic = "telemetry.electrical.created"
)

// Broker defines the interface for event publishing and subscribing.
type Broker interface {
	Publish(topic Topic, payload interface{})
	Subscribe(topic Topic) <-chan interface{}
	Close()
}

// inMemoryBroker implements an in-memory topic-based Pub/Sub broker.
type inMemoryBroker struct {
	mu          sync.RWMutex
	subscribers map[Topic][]chan interface{}
	closed      bool
}

// NewBroker creates a new instance of Broker.
func NewBroker() Broker {
	return &inMemoryBroker{
		subscribers: make(map[Topic][]chan interface{}),
	}
}

// Publish sends a payload to all subscribers of a given topic.
// It is non-blocking to prevent slow subscribers from stalling the main flow.
func (b *inMemoryBroker) Publish(topic Topic, payload interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return
	}

	for _, ch := range b.subscribers[topic] {
		select {
		case ch <- payload:
		default:
			// Buffer full, drop message to prevent blocking
		}
	}
}

// Subscribe returns a read-only channel for a given topic.
func (b *inMemoryBroker) Subscribe(topic Topic) <-chan interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		// Return closed channel if broker is already closed
		ch := make(chan interface{})
		close(ch)
		return ch
	}

	// We use a generous buffer to prevent message loss on burst events
	ch := make(chan interface{}, 200)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}

// Close cleanly shuts down the broker and closes all subscriber channels.
func (b *inMemoryBroker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}
	b.closed = true

	for _, channels := range b.subscribers {
		for _, ch := range channels {
			close(ch)
		}
	}
	b.subscribers = nil
}
