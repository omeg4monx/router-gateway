package eventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
)

func TestPublishNeverBlocksSlowSubscriber(t *testing.T) {
	bus := eventbus.New()
	ch, unsubscribe := bus.Subscribe("usage.recorded", 1)

	start := time.Now()
	for i := range 100 {
		bus.Publish("usage.recorded", i)
	}
	if elapsed := time.Since(start); elapsed >= 100*time.Millisecond {
		t.Fatalf("Publish took %v, expected < 100ms", elapsed)
	}

	// Only 1 buffered event should be present (excess dropped without blocking).
	select {
	case evt, ok := <-ch:
		if !ok {
			t.Fatal("channel closed unexpectedly before unsubscribe")
		}
		if evt.Topic != "usage.recorded" {
			t.Errorf("evt.Topic = %q, want %q", evt.Topic, "usage.recorded")
		}
		if evt.Timestamp.IsZero() {
			t.Error("evt.Timestamp is zero")
		}
	default:
		t.Fatal("expected at least 1 buffered event on channel")
	}

	// Calling unsubscribe twice must be idempotent and close the channel.
	unsubscribe()
	unsubscribe()

	if _, ok := <-ch; ok {
		t.Fatal("expected channel to be closed after unsubscribe")
	}

	// Subsequent Publish must succeed without panicking.
	bus.Publish("usage.recorded", "after-unsub")
}

func TestPublishDeliversToMultipleSubscribersOnMatchingTopic(t *testing.T) {
	bus := eventbus.New()
	ch1, unsub1 := bus.Subscribe("config.changed", 4)
	defer unsub1()
	ch2, unsub2 := bus.Subscribe("config.changed", 4)
	defer unsub2()
	chOther, unsubOther := bus.Subscribe("other", 4)
	defer unsubOther()

	payload := map[string]string{"key": "log_level"}
	bus.Publish("config.changed", payload)

	for idx, ch := range []<-chan eventbus.Event{ch1, ch2} {
		select {
		case evt := <-ch:
			if evt.Topic != "config.changed" {
				t.Errorf("sub%d topic = %q, want config.changed", idx+1, evt.Topic)
			}
			got, ok := evt.Payload.(map[string]string)
			if !ok || got["key"] != "log_level" {
				t.Errorf("sub%d payload = %#v, want %#v", idx+1, evt.Payload, payload)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("sub%d did not receive event", idx+1)
		}
	}

	select {
	case evt := <-chOther:
		t.Fatalf("other subscriber received unexpected event: %+v", evt)
	default:
	}
}

func TestConcurrentPublishSubscribeUnsubscribe(t *testing.T) {
	bus := eventbus.New()
	var wg sync.WaitGroup

	// Concurrent publishers
	for p := range 6 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range 200 {
				bus.Publish("gateway.request", id*1000+i)
			}
		}(p)
	}

	// Concurrent subscribers joining and leaving while publishing happens
	for s := range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				ch, unsub := bus.Subscribe("gateway.request", 2)
				select {
				case <-ch:
				default:
				}
				unsub()
				unsub()
			}
		}()
		_ = s
	}

	wg.Wait()
}
