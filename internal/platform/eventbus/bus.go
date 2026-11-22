// Package eventbus provides a non-blocking, in-process topic publish/subscribe
// bus safe for concurrent publishers and subscribers.
package eventbus

import (
	"sync"
	"time"
)

// Event represents a published message delivered to topic subscribers.
type Event struct {
	Topic     string
	Payload   any
	Timestamp time.Time
}

// Bus routes events by topic to active subscribers without blocking publishers.
type Bus struct {
	mu     sync.RWMutex
	subs   map[string]map[uint64]chan Event
	nextID uint64
}

// New creates a new in-process Event Bus.
func New() *Bus {
	return &Bus{
		subs: make(map[string]map[uint64]chan Event),
	}
}

// Publish emits an event to all subscribers registered for topic. Slow
// subscribers whose channel buffer is full drop the excess event instead of
// blocking the publisher.
func (b *Bus) Publish(topic string, payload any) {
	evt := Event{
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	}

	// Hold RLock across the non-blocking send so unsubscribe (which acquires
	// Lock) cannot close a channel concurrently with a send.
	b.mu.RLock()
	defer b.mu.RUnlock()

	topicSubs := b.subs[topic]
	for _, ch := range topicSubs {
		select {
		case ch <- evt:
		default:
		}
	}
}

// Subscribe registers a buffered channel of size buf for events matching topic
// and returns the receive-only channel along with an idempotent unsubscribe
// function that removes the subscription and closes the channel.
func (b *Bus) Subscribe(topic string, buf int) (<-chan Event, func()) {
	if buf < 0 {
		buf = 0
	}

	ch := make(chan Event, buf)

	b.mu.Lock()
	b.nextID++
	id := b.nextID
	if b.subs[topic] == nil {
		b.subs[topic] = make(map[uint64]chan Event)
	}
	b.subs[topic][id] = ch
	b.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			b.mu.Lock()
			defer b.mu.Unlock()

			topicSubs, ok := b.subs[topic]
			if !ok {
				return
			}
			if _, exists := topicSubs[id]; !exists {
				return
			}
			delete(topicSubs, id)
			if len(topicSubs) == 0 {
				delete(b.subs, topic)
			}
			close(ch)
		})
	}

	return ch, unsubscribe
}
